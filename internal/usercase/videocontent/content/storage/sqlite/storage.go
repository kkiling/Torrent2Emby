package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kkiling/goplatform/log"
	"github.com/kkiling/goplatform/storagebase/sqlitebase"

	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/common"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/content"
)

type Storage struct {
	base *sqlitebase.Storage
}

func NewStorage(config sqlitebase.Config, logger log.Logger) (*Storage, error) {
	s, err := sqlitebase.NewStorage(config, logger)
	if err != nil {
		return nil, err
	}
	return &Storage{
		base: s,
	}, nil
}

func NewTestStorage(base *sqlitebase.Storage) *Storage {
	return &Storage{
		base: base,
	}
}

func (s *Storage) CreateVideoContent(ctx context.Context, videoContent *content.VideoContent) error {
	statesJSON, err := json.Marshal(videoContent.State)
	if err != nil {
		return fmt.Errorf("error marshal states: %w", err)
	}

	var movieID *uint64
	var tvshowID *uint64
	var seasonNumber *uint8
	if videoContent.ContentID.MovieID != nil {
		movieID = videoContent.ContentID.MovieID
	}
	if videoContent.ContentID.TVShow != nil {
		tvshowID = &videoContent.ContentID.TVShow.ID
		seasonNumber = &videoContent.ContentID.TVShow.SeasonNumber
	}

	query := `
        INSERT OR REPLACE INTO video_content (
            id,
            created_at,
            movie_id,
            tvshow_id,
            season_number,
            delivery_status,
            states_json
        ) VALUES (?, ?, ?, ?, ?, ?, ?)
    `

	_, err = s.base.Next(ctx).ExecContext(ctx, query,
		videoContent.ID.String(),
		videoContent.CreatedAt.Format(time.RFC3339),
		movieID,
		tvshowID,
		seasonNumber,
		videoContent.DeliveryStatus,
		string(statesJSON),
	)
	if err != nil {
		return s.base.HandleError(err)
	}

	return nil
}

func (s *Storage) getVideoContents(rows *sql.Rows) ([]content.VideoContent, error) {
	var results []content.VideoContent
	for rows.Next() {
		var (
			vc                content.VideoContent
			movieID           sql.NullInt64
			tvshowID          sql.NullInt64
			seasonNumber      sql.NullInt16
			deliveryStatusStr string
			statesJSON        sql.NullString
		)

		err := rows.Scan(
			&vc.ID,
			&vc.CreatedAt,
			&movieID,
			&tvshowID,
			&seasonNumber,
			&deliveryStatusStr,
			&statesJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan video_content: %w", err)
		}

		// rebuild ContentID
		var cid common.ContentID
		if movieID.Valid {
			mid := uint64(movieID.Int64)
			cid.MovieID = &mid
		}
		if tvshowID.Valid {
			tid := uint64(tvshowID.Int64)
			sn := uint8(0)
			if seasonNumber.Valid {
				sn = uint8(seasonNumber.Int16)
			}
			cid.TVShow = &common.TVShowID{ID: tid, SeasonNumber: sn}
		}
		vc.ContentID = cid

		// DeliveryStatus
		vc.DeliveryStatus = content.DeliveryStatus(deliveryStatusStr)

		// unmarshal states
		if statesJSON.Valid && statesJSON.String != "" {
			if err := json.Unmarshal([]byte(statesJSON.String), &vc.State); err != nil {
				return nil, fmt.Errorf("failed to unmarshal states_json: %w", err)
			}
		}

		results = append(results, vc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return results, nil
}

func (s *Storage) GetVideoContents(ctx context.Context, contentID common.ContentID) ([]content.VideoContent, error) {

	next := s.base.Next(ctx)
	var rows *sql.Rows
	var err error

	if contentID.MovieID != nil {
		rows, err = next.QueryContext(ctx, `
            SELECT 
                id,
                created_at,
                movie_id,
                tvshow_id,
                season_number,
                delivery_status,
                states_json
            FROM video_content
            WHERE movie_id = ?
        `, *contentID.MovieID)
	} else if contentID.TVShow != nil {
		rows, err = next.QueryContext(ctx, `
            SELECT 
                id,
                created_at,
                movie_id,
                tvshow_id,
                season_number,
                delivery_status,
                states_json
            FROM video_content
            WHERE tvshow_id = ? AND season_number = ?
        `, contentID.TVShow.ID, contentID.TVShow.SeasonNumber)
	} else {
		return nil, fmt.Errorf("invalid contentID: both MovieID and TVShow are nil")
	}

	if err != nil {
		return nil, s.base.HandleError(err)
	}
	defer rows.Close()

	return s.getVideoContents(rows)
}

func (s *Storage) UpdateVideoContent(ctx context.Context, id uuid.UUID, videoContent *content.UpdateVideoContent) error {
	query := `
        UPDATE video_content
        SET delivery_status = ?
        WHERE id = ?
    `
	_, err := s.base.Next(ctx).ExecContext(ctx, query,
		videoContent.DeliveryStatus,
		id,
	)
	if err != nil {
		return s.base.HandleError(err)
	}

	return nil
}

func (s *Storage) GetVideoContentsByStatus(ctx context.Context, status content.DeliveryStatus, limit int) ([]content.VideoContent, error) {
	query := `
		SELECT
			id,
			created_at,
			movie_id,
			tvshow_id,
			season_number,
			delivery_status,
			states_json
		FROM video_content
		WHERE delivery_status = ?
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := s.base.Next(ctx).QueryContext(ctx, query, status, limit)
	if err != nil {
		return nil, s.base.HandleError(err)
	}
	defer rows.Close()

	return s.getVideoContents(rows)
}

func (s *Storage) RunTransaction(ctx context.Context, txFunc func(ctxTx context.Context) error) error {
	return s.base.RunTransaction(ctx, txFunc)
}

package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/samber/lo"

	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
)

func (s *Storage) saveOrUpdateImage(ctx context.Context, img *tvshowlibrary.Image) error {
	if img == nil {
		return nil
	}
	_, err := s.base.Next(ctx).ExecContext(ctx, `
        INSERT OR REPLACE INTO images (id, w342, original) VALUES (?, ?, ?)
    `, img.ID, img.W342, img.Original)

	if err != nil {
		return s.base.HandleError(err)
	}

	return nil
}

func (s *Storage) saveOrUpdateSeason(ctx context.Context, tvID uint64, season tvshowlibrary.Season) error {
	err := s.saveOrUpdateImage(ctx, season.Poster)
	if err != nil {
		return fmt.Errorf("error saving image: %w", err)
	}

	seasonQuery := `
            INSERT OR REPLACE INTO seasons (
                id, tv_show_id, air_date, episode_count,
                name, overview, poster_id,
                season_number, vote_average
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        `

	posterID := lo.EmptyableToPtr(lo.FromPtr(season.Poster).ID)
	_, err = s.base.Next(ctx).ExecContext(ctx, seasonQuery,
		season.ID,
		tvID,
		season.AirDate,
		season.EpisodeCount,
		season.Name,
		season.Overview,
		posterID,
		season.SeasonNumber,
		season.VoteAverage,
	)

	if err != nil {
		return s.base.HandleError(err)
	}

	return nil
}

func (s *Storage) getImage(ctx context.Context, imageID string) (*tvshowlibrary.Image, error) {
	var img tvshowlibrary.Image
	err := s.base.Next(ctx).QueryRowContext(ctx, `
        SELECT id, w342, original FROM images WHERE id = ?
    `, imageID).Scan(&img.ID, &img.W342, &img.Original)

	if err != nil {
		return nil, s.base.HandleError(err)
	}

	return &img, nil
}

func (s *Storage) getSeasons(ctx context.Context, tvID uint64) ([]tvshowlibrary.Season, error) {
	rows, err := s.base.Next(ctx).QueryContext(ctx, `
        SELECT 
            id, air_date, episode_count, name,
            overview, poster_id, season_number, vote_average
        FROM seasons
        WHERE tv_show_id = ?
        ORDER BY season_number
    `, tvID)

	if err != nil {
		return nil, fmt.Errorf("failed to query seasons: %w", err)
	}
	defer rows.Close()

	var seasons []tvshowlibrary.Season
	for rows.Next() {
		var season tvshowlibrary.Season
		var posterID sql.NullString

		err := rows.Scan(
			&season.ID,
			&season.AirDate,
			&season.EpisodeCount,
			&season.Name,
			&season.Overview,
			&posterID,
			&season.SeasonNumber,
			&season.VoteAverage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan season: %w", err)
		}

		if posterID.Valid {
			season.Poster, err = s.getImage(ctx, posterID.String)
			if err != nil {
				return nil, fmt.Errorf("failed to get season poster: %w", err)
			}
		}
		seasons = append(seasons, season)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return seasons, nil
}

func (s *Storage) SaveOrUpdateTVShow(ctx context.Context, tvShow *tvshowlibrary.TVShow) error {
	return s.RunTransaction(ctx, func(tCtx context.Context) error {
		// Сохраняем изображения (постер и backdrop)
		err := s.saveOrUpdateImage(tCtx, tvShow.Poster)
		if err != nil {
			return fmt.Errorf("error saving image: %w", err)
		}

		err = s.saveOrUpdateImage(tCtx, tvShow.Backdrop)
		if err != nil {
			return fmt.Errorf("error saving image: %w", err)
		}

		// Преобразуем массивы в JSON
		genresJSON, err := json.Marshal(tvShow.Genres)
		if err != nil {
			return fmt.Errorf("error marshal genres: %w", err)
		}

		originCountryJSON, err := json.Marshal(tvShow.OriginCountry)
		if err != nil {
			return fmt.Errorf("error marshal origin: %w", err)
		}

		// Сохраняем основную информацию о сериале
		query := `
        INSERT OR REPLACE INTO tv_shows (
            id, name, original_name, overview, poster_id,
            first_air_date, vote_average, vote_count, popularity,
            backdrop_id, genres, last_air_date, next_episode_to_air,
            number_of_episodes, number_of_seasons, origin_country,
            status, tagline, type
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

		_, err = s.base.Next(tCtx).ExecContext(ctx, query,
			tvShow.ID,
			tvShow.Name,
			tvShow.OriginalName,
			tvShow.Overview,
			lo.EmptyableToPtr(lo.FromPtr(tvShow.Poster).ID),
			tvShow.FirstAirDate,
			tvShow.VoteAverage,
			tvShow.VoteCount,
			tvShow.Popularity,
			lo.EmptyableToPtr(lo.FromPtr(tvShow.Backdrop).ID),
			genresJSON,
			tvShow.LastAirDate,
			tvShow.NextEpisodeToAir,
			tvShow.NumberOfEpisodes,
			tvShow.NumberOfSeasons,
			originCountryJSON,
			tvShow.Status,
			tvShow.Tagline,
			tvShow.Type,
		)

		if err != nil {
			return s.base.HandleError(err)
		}

		// Сохраняем сезоны
		for _, season := range tvShow.Seasons {
			err = s.saveOrUpdateSeason(tCtx, tvShow.ID, season)
			if err != nil {
				return fmt.Errorf("error saving season: %w", err)
			}
		}

		return nil
	})
}

func (s *Storage) GetTVShow(ctx context.Context, tvID uint64) (*tvshowlibrary.TVShow, error) {
	var tvShow tvshowlibrary.TVShow
	var genresJSON, originCountryJSON string
	var posterID, backdropID sql.NullString

	// Загружаем основную информацию о сериале
	err := s.base.Next(ctx).QueryRowContext(ctx, `
        SELECT 
            id, name, original_name, overview, poster_id,
            first_air_date, vote_average, vote_count, popularity,
            backdrop_id, genres, last_air_date, next_episode_to_air,
            number_of_episodes, number_of_seasons, origin_country,
            status, tagline, type
        FROM tv_shows
        WHERE id = ?
    `, tvID).Scan(
		&tvShow.ID,
		&tvShow.Name,
		&tvShow.OriginalName,
		&tvShow.Overview,
		&posterID,
		&tvShow.FirstAirDate,
		&tvShow.VoteAverage,
		&tvShow.VoteCount,
		&tvShow.Popularity,
		&backdropID,
		&genresJSON,
		&tvShow.LastAirDate,
		&tvShow.NextEpisodeToAir,
		&tvShow.NumberOfEpisodes,
		&tvShow.NumberOfSeasons,
		&originCountryJSON,
		&tvShow.Status,
		&tvShow.Tagline,
		&tvShow.Type,
	)

	if err != nil {
		return nil, s.base.HandleError(err)
	}

	// Загружаем изображения
	if posterID.Valid {
		tvShow.Poster, err = s.getImage(ctx, posterID.String)
		if err != nil {
			return nil, fmt.Errorf("failed to get poster image: %w", err)
		}
	}

	if backdropID.Valid {
		tvShow.Backdrop, err = s.getImage(ctx, backdropID.String)
		if err != nil {
			return nil, fmt.Errorf("failed to get backdrop image: %w", err)
		}
	}

	// Десериализуем массивы
	if err := json.Unmarshal([]byte(genresJSON), &tvShow.Genres); err != nil {
		return nil, fmt.Errorf("failed to unmarshal genres: %w", err)
	}

	if err := json.Unmarshal([]byte(originCountryJSON), &tvShow.OriginCountry); err != nil {
		return nil, fmt.Errorf("failed to unmarshal origin country: %w", err)
	}

	// Загружаем сезоны
	seasons, err := s.getSeasons(ctx, tvID)
	if err != nil {
		return nil, fmt.Errorf("failed to get seasons: %w", err)
	}
	tvShow.Seasons = seasons
	return &tvShow, nil
}

func (s *Storage) GetTVShows(ctx context.Context) ([]tvshowlibrary.TVShowShort, error) {
	rows, err := s.base.Next(ctx).QueryContext(ctx, `
        SELECT 
            ts.id, ts.name, ts.original_name, ts.overview,
            ts.first_air_date, ts.vote_average, ts.vote_count, ts.popularity,
            poster.w342, poster.original
        FROM tv_shows ts
        LEFT JOIN images poster ON ts.poster_id = poster.id
        ORDER BY ts.popularity DESC
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to query tv shows: %w", err)
	}
	defer rows.Close()

	var shows []tvshowlibrary.TVShowShort
	for rows.Next() {
		var show tvshowlibrary.TVShowShort
		var firstAirDateStr string
		var posterW342, posterOriginal sql.NullString

		err := rows.Scan(
			&show.ID,
			&show.Name,
			&show.OriginalName,
			&show.Overview,
			&firstAirDateStr,
			&show.VoteAverage,
			&show.VoteCount,
			&show.Popularity,
			&posterW342,
			&posterOriginal,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tv show: %w", err)
		}

		// Парсим дату
		if firstAirDateStr != "" {
			show.FirstAirDate, _ = time.Parse("2006-01-02", firstAirDateStr)
		}

		// Заполняем изображение постера если есть
		if posterW342.Valid && posterOriginal.Valid {
			show.Poster = &tvshowlibrary.Image{
				W342:     posterW342.String,
				Original: posterOriginal.String,
			}
		}

		shows = append(shows, show)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return shows, nil
}

func (s *Storage) GetSeasonEpisodes(ctx context.Context, tvID uint64, seasonNumber int) ([]tvshowlibrary.Episode, error) {
	// First get the season ID for the given TV show and season number
	var seasonID uint64
	err := s.base.Next(ctx).QueryRowContext(ctx, `
        SELECT id FROM seasons 
        WHERE tv_show_id = ? AND season_number = ?
    `, tvID, seasonNumber).Scan(&seasonID)

	if err != nil {
		return nil, fmt.Errorf("failed to get season ID: %w", s.base.HandleError(err))
	}

	// Now get all episodes for this season
	rows, err := s.base.Next(ctx).QueryContext(ctx, `
        SELECT 
            e.id, e.air_date, e.episode_number, e.episode_type,
            e.name, e.overview, e.runtime, e.still_id,
            e.vote_average, e.vote_count
        FROM episodes e
        WHERE e.season_id = ?
        ORDER BY e.episode_number
    `, seasonID)

	if err != nil {
		return nil, fmt.Errorf("failed to query episodes: %w", s.base.HandleError(err))
	}
	defer rows.Close()

	var episodes []tvshowlibrary.Episode
	for rows.Next() {
		var episode tvshowlibrary.Episode
		var stillID sql.NullString

		err = rows.Scan(
			&episode.ID,
			&episode.AirDate,
			&episode.EpisodeNumber,
			&episode.EpisodeType,
			&episode.Name,
			&episode.Overview,
			&episode.Runtime,
			&stillID,
			&episode.VoteAverage,
			&episode.VoteCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan episode: %w", err)
		}

		// Get still image if exists
		if stillID.Valid {
			episode.Still, err = s.getImage(ctx, stillID.String)
			if err != nil {
				return nil, fmt.Errorf("failed to get still image: %w", err)
			}
		}

		episodes = append(episodes, episode)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return episodes, nil
}

func (s *Storage) SaveOrUpdateSeasonEpisode(ctx context.Context, tvID uint64, seasonNumber int, episodes []tvshowlibrary.Episode) error {
	return s.RunTransaction(ctx, func(tCtx context.Context) error {
		// First get the season ID for the given TV show and season number
		var seasonID uint64
		err := s.base.Next(tCtx).QueryRowContext(tCtx, `
            SELECT id FROM seasons 
            WHERE tv_show_id = ? AND season_number = ?
        `, tvID, seasonNumber).Scan(&seasonID)

		if err != nil {
			return fmt.Errorf("failed to get season ID: %w", s.base.HandleError(err))
		}

		// Prepare the query for inserting/updating episodes
		episodeQuery := `
            INSERT OR REPLACE INTO episodes (
                id, season_id, air_date, episode_number,
                episode_type, name, overview, runtime,
                still_id, vote_average, vote_count
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        `

		// Process each episode
		for _, episode := range episodes {
			// Save still image if exists
			if episode.Still != nil {
				err = s.saveOrUpdateImage(tCtx, episode.Still)
				if err != nil {
					return fmt.Errorf("error saving still image: %w", err)
				}
			}

			// Insert/update episode
			stillID := lo.EmptyableToPtr(lo.FromPtr(episode.Still).ID)
			_, err = s.base.Next(tCtx).ExecContext(tCtx, episodeQuery,
				episode.ID,
				seasonID,
				episode.AirDate,
				episode.EpisodeNumber,
				episode.EpisodeType,
				episode.Name,
				episode.Overview,
				episode.Runtime,
				stillID,
				episode.VoteAverage,
				episode.VoteCount,
			)

			if err != nil {
				return s.base.HandleError(err)
			}
		}

		return nil
	})
}

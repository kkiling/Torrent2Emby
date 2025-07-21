package sqlite

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary"
	"github.com/kkiling/torrent2emby/internal/usercase/tvshowlibrary/storage"
)

func TestStorage_SaveAndGetTVShow(t *testing.T) {
	t.Parallel()

	s := setupTestDB(t)
	ctx := context.Background()

	initTvShow := func() *tvshowlibrary.TVShow {
		// Тестовые данные
		testImage := &tvshowlibrary.Image{
			ID:       "1234566", // uuid.NewString(),
			W342:     "path/w342",
			Original: "path/original",
		}

		testSeason := tvshowlibrary.Season{
			ID:           1,
			AirDate:      time.Now().AddDate(-13, 0, 0), // "2023-01-01",
			EpisodeCount: 10,
			Name:         "Season 1",
			Overview:     "Test season overview",
			Poster:       testImage,
			SeasonNumber: 1,
			VoteAverage:  8.5,
		}

		testTVShow := &tvshowlibrary.TVShow{
			TVShowShort: tvshowlibrary.TVShowShort{
				ID:           1,
				Name:         "Test Show",
				OriginalName: "Test Show Original",
				Overview:     "Test overview",
				Poster:       testImage,
				FirstAirDate: time.Now(),
				VoteAverage:  7.8,
				VoteCount:    100,
				Popularity:   50.5,
			},
			Backdrop:         testImage,
			Genres:           []string{"Drama", "Comedy"},
			LastAirDate:      time.Now(),
			NextEpisodeToAir: time.Now().AddDate(0, 1, 0),
			NumberOfEpisodes: 10,
			NumberOfSeasons:  1,
			OriginCountry:    []string{"US", "UK"},
			Status:           "Running",
			Tagline:          "Test tagline",
			Type:             "Scripted",
			Seasons:          []tvshowlibrary.Season{testSeason},
		}
		return testTVShow
	}

	t.Run("SaveTVShow success", func(t *testing.T) {
		testTVShow := initTvShow()

		err := s.SaveOrUpdateTVShow(ctx, testTVShow)
		require.NoError(t, err)

		// Проверяем, что данные сохранились
		var count int
		err = s.next(ctx).QueryRowContext(ctx, "SELECT COUNT(*) FROM tv_shows WHERE id = ?", testTVShow.ID).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("SaveTVShow with nil images", func(t *testing.T) {
		testTVShow := initTvShow()
		testTVShow.Poster = nil
		testTVShow.Backdrop = nil
		testTVShow.Seasons[0].Poster = nil

		err := s.SaveOrUpdateTVShow(ctx, testTVShow)
		require.NoError(t, err)
	})

	t.Run("SaveTVShow empty arrays", func(t *testing.T) {
		testTVShow := initTvShow()
		testTVShow.Genres = []string{}
		testTVShow.OriginCountry = []string{}

		err := s.SaveOrUpdateTVShow(ctx, testTVShow)
		require.NoError(t, err)
	})

	t.Run("GetTVShow success", func(t *testing.T) {
		testTVShow := initTvShow()
		// Сначала сохраняем тестовые данные
		err := s.SaveOrUpdateTVShow(ctx, testTVShow)
		require.NoError(t, err)

		// Получаем данные
		result, err := s.GetTVShow(ctx, testTVShow.ID)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Сравниваем основные поля
		require.Equal(t, testTVShow.ID, result.ID)
		require.Equal(t, testTVShow.Name, result.Name)
		require.Equal(t, testTVShow.OriginalName, result.OriginalName)
		require.Equal(t, testTVShow.Overview, result.Overview)
		require.Equal(t, testTVShow.Genres, result.Genres)
		require.Equal(t, testTVShow.OriginCountry, result.OriginCountry)

		// Проверяем изображения
		require.NotNil(t, result.Poster)
		require.Equal(t, testTVShow.Poster.W342, result.Poster.W342)
		require.Equal(t, testTVShow.Poster.Original, result.Poster.Original)

		// Проверяем сезоны
		require.Len(t, result.Seasons, 1)
		require.Equal(t, testTVShow.Seasons[0].Name, result.Seasons[0].Name)
		require.NotNil(t, result.Seasons[0].Poster)
	})

	t.Run("GetTVShow not found", func(t *testing.T) {
		_, err := s.GetTVShow(ctx, 9999)
		require.ErrorIs(t, err, storage.ErrNotFound)
	})

	t.Run("GetTVShow with nil images", func(t *testing.T) {
		testTVShow := initTvShow()
		testTVShow.Poster = nil
		testTVShow.Backdrop = nil
		testTVShow.Seasons[0].Poster = nil

		err := s.SaveOrUpdateTVShow(ctx, testTVShow)
		require.NoError(t, err)

		result, err := s.GetTVShow(ctx, testTVShow.ID)
		require.NoError(t, err)
		require.Nil(t, result.Poster)
		require.Nil(t, result.Backdrop)
		require.Nil(t, result.Seasons[0].Poster)
	})

	t.Run("Save and Get roundtrip", func(t *testing.T) {
		testTVShow := initTvShow()
		err := s.SaveOrUpdateTVShow(ctx, testTVShow)
		require.NoError(t, err)

		result, err := s.GetTVShow(ctx, testTVShow.ID)
		require.NoError(t, err)

		// Сравниваем все поля через JSON для глубокого сравнения
		expectedJSON, err := json.Marshal(testTVShow)
		require.NoError(t, err)

		actualJSON, err := json.Marshal(result)
		require.NoError(t, err)

		require.JSONEq(t, string(expectedJSON), string(actualJSON))
	})
}

func TestGetTVShows(t *testing.T) {
	t.Parallel()

	s := setupTestDB(t)
	ctx := context.Background()

	// Тестовые данные
	now := time.Now()
	testImage := &tvshowlibrary.Image{
		ID:       uuid.NewString(),
		W342:     "path/w342",
		Original: "path/original",
	}

	testShows := []*tvshowlibrary.TVShow{
		{
			TVShowShort: tvshowlibrary.TVShowShort{
				ID:           1235667,
				Name:         "Show 1",
				OriginalName: "Show 1 Original",
				Overview:     "Overview 1",
				Poster:       testImage,
				FirstAirDate: now.AddDate(0, -1, 0),
				VoteAverage:  8.5,
				VoteCount:    100,
				Popularity:   90.0,
			},
		},
		{
			TVShowShort: tvshowlibrary.TVShowShort{
				ID:           24441123,
				Name:         "Show 2",
				OriginalName: "Show 2 Original",
				Overview:     "Overview 2",
				Poster:       nil, // Тест с отсутствующим постером
				FirstAirDate: now.AddDate(0, -2, 0),
				VoteAverage:  7.2,
				VoteCount:    50,
				Popularity:   80.0,
			},
		},
	}

	// Подготовка: сохраняем тестовые данные
	for _, show := range testShows {
		err := s.SaveOrUpdateTVShow(ctx, show)
		require.NoError(t, err)
	}

	t.Run("successful get all shows", func(t *testing.T) {
		result, err := s.GetTVShows(ctx)

		result = lo.Filter(result, func(item tvshowlibrary.TVShowShort, _ int) bool {
			return item.ID == testShows[0].ID || item.ID == testShows[1].ID
		})

		require.NoError(t, err)
		require.Len(t, result, len(testShows))

		require.Equal(t, testShows[0].ID, result[0].ID)
		require.Equal(t, testShows[0].Popularity, result[0].Popularity)
		require.Equal(t, testShows[0].OriginalName, result[0].OriginalName)
		require.NotNil(t, result[0].Poster)
		require.Equal(t, testShows[0].Poster.W342, result[0].Poster.W342)

		require.Equal(t, testShows[1].ID, result[1].ID)
		require.Equal(t, testShows[1].Popularity, result[1].Popularity)
		require.Equal(t, testShows[1].OriginalName, result[1].OriginalName)
		require.Nil(t, result[1].Poster)
	})
}

func TestStorage_SeasonEpisodes(t *testing.T) {
	t.Parallel()

	s := setupTestDB(t)
	ctx := context.Background()

	// Helper function to initialize test data
	initTestData := func() (*tvshowlibrary.TVShow, []tvshowlibrary.Episode) {
		testImage := &tvshowlibrary.Image{
			ID:       "episode_still_123",
			W342:     "path/w342/still",
			Original: "path/original/still",
		}

		testSeason := tvshowlibrary.Season{
			ID:           1,
			AirDate:      time.Now().AddDate(-1, 0, 0),
			EpisodeCount: 2,
			Name:         "Test Season",
			Overview:     "Season overview",
			Poster:       testImage,
			SeasonNumber: 1,
			VoteAverage:  8.0,
		}

		testTVShow := &tvshowlibrary.TVShow{
			TVShowShort: tvshowlibrary.TVShowShort{
				ID:           1,
				Name:         "Test Show",
				OriginalName: "Test Show Original",
				Overview:     "Test overview",
				Poster:       testImage,
				FirstAirDate: time.Now(),
				VoteAverage:  7.8,
				VoteCount:    100,
				Popularity:   50.5,
			},
			Seasons: []tvshowlibrary.Season{testSeason},
		}

		testEpisodes := []tvshowlibrary.Episode{
			{
				ID:            101,
				AirDate:       time.Now().AddDate(0, 0, -7),
				EpisodeNumber: 1,
				EpisodeType:   "standard",
				Name:          "Episode 1",
				Overview:      "First episode",
				Runtime:       45,
				Still:         testImage,
				VoteAverage:   8.5,
				VoteCount:     50,
			},
			{
				ID:            102,
				AirDate:       time.Now().AddDate(0, 0, -6),
				EpisodeNumber: 2,
				EpisodeType:   "standard",
				Name:          "Episode 2",
				Overview:      "Second episode",
				Runtime:       42,
				Still:         nil, // Test nil still
				VoteAverage:   8.7,
				VoteCount:     60,
			},
		}

		return testTVShow, testEpisodes
	}

	t.Run("SaveOrUpdateSeasonEpisode success", func(t *testing.T) {
		testTVShow, testEpisodes := initTestData()

		// First save the TV show with season
		err := s.SaveOrUpdateTVShow(ctx, testTVShow)
		require.NoError(t, err)

		// Save episodes
		err = s.SaveOrUpdateSeasonEpisode(ctx, testTVShow.ID, 1, testEpisodes)
		require.NoError(t, err)

		// Verify episodes were saved
		var count int
		err = s.next(ctx).QueryRowContext(ctx,
			"SELECT COUNT(*) FROM episodes WHERE season_id = ?",
			testTVShow.Seasons[0].ID).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, len(testEpisodes), count)
	})

	t.Run("SaveOrUpdateSeasonEpisode with empty episodes", func(t *testing.T) {
		testTVShow, _ := initTestData()
		err := s.SaveOrUpdateTVShow(ctx, testTVShow)
		require.NoError(t, err)

		err = s.SaveOrUpdateSeasonEpisode(ctx, testTVShow.ID, 1, []tvshowlibrary.Episode{})
		require.NoError(t, err)
	})

	t.Run("GetSeasonEpisodes success", func(t *testing.T) {
		testTVShow, testEpisodes := initTestData()
		err := s.SaveOrUpdateTVShow(ctx, testTVShow)
		require.NoError(t, err)

		err = s.SaveOrUpdateSeasonEpisode(ctx, testTVShow.ID, 1, testEpisodes)
		require.NoError(t, err)

		// Retrieve episodes
		episodes, err := s.GetSeasonEpisodes(ctx, testTVShow.ID, 1)
		require.NoError(t, err)
		require.Len(t, episodes, len(testEpisodes))

		// Verify first episode
		require.Equal(t, testEpisodes[0].ID, episodes[0].ID)
		require.Equal(t, testEpisodes[0].Name, episodes[0].Name)
		require.Equal(t, testEpisodes[0].Overview, episodes[0].Overview)
		require.NotNil(t, episodes[0].Still)
		require.Equal(t, testEpisodes[0].Still.W342, episodes[0].Still.W342)

		// Verify second episode has nil still
		require.Nil(t, episodes[1].Still)
	})

	t.Run("GetSeasonEpisodes not found", func(t *testing.T) {
		testTVShow, _ := initTestData()
		err := s.SaveOrUpdateTVShow(ctx, testTVShow)
		require.NoError(t, err)

		// Try to get episodes for non-existent season
		_, err = s.GetSeasonEpisodes(ctx, testTVShow.ID, 999)
		require.ErrorIs(t, err, storage.ErrNotFound)

		// Try to get episodes for non-existent show
		_, err = s.GetSeasonEpisodes(ctx, 999, 1)
		require.ErrorIs(t, err, storage.ErrNotFound)
	})

	t.Run("Save and Get roundtrip", func(t *testing.T) {
		testTVShow, testEpisodes := initTestData()
		err := s.SaveOrUpdateTVShow(ctx, testTVShow)
		require.NoError(t, err)

		err = s.SaveOrUpdateSeasonEpisode(ctx, testTVShow.ID, 1, testEpisodes)
		require.NoError(t, err)

		retrievedEpisodes, err := s.GetSeasonEpisodes(ctx, testTVShow.ID, 1)
		require.NoError(t, err)

		// Compare all fields through JSON for deep comparison
		expectedJSON, err := json.Marshal(testEpisodes)
		require.NoError(t, err)

		// Need to handle nil Still manually before marshaling
		for i := range retrievedEpisodes {
			if retrievedEpisodes[i].Still != nil && retrievedEpisodes[i].Still.ID == "" {
				retrievedEpisodes[i].Still = nil
			}
		}

		actualJSON, err := json.Marshal(retrievedEpisodes)
		require.NoError(t, err)

		require.JSONEq(t, string(expectedJSON), string(actualJSON))
	})

	t.Run("Update existing episodes", func(t *testing.T) {
		testTVShow, testEpisodes := initTestData()
		err := s.SaveOrUpdateTVShow(ctx, testTVShow)
		require.NoError(t, err)

		// First save
		err = s.SaveOrUpdateSeasonEpisode(ctx, testTVShow.ID, 1, testEpisodes)
		require.NoError(t, err)

		// Update episodes
		updatedEpisodes := make([]tvshowlibrary.Episode, len(testEpisodes))
		copy(updatedEpisodes, testEpisodes)
		updatedEpisodes[0].Name = "Updated Episode Name"
		updatedEpisodes[0].Overview = "Updated overview"
		updatedEpisodes[0].VoteAverage = 9.0

		err = s.SaveOrUpdateSeasonEpisode(ctx, testTVShow.ID, 1, updatedEpisodes)
		require.NoError(t, err)

		// Retrieve and verify updates
		retrievedEpisodes, err := s.GetSeasonEpisodes(ctx, testTVShow.ID, 1)
		require.NoError(t, err)

		require.Equal(t, "Updated Episode Name", retrievedEpisodes[0].Name)
		require.Equal(t, "Updated overview", retrievedEpisodes[0].Overview)
		require.Equal(t, 9.0, retrievedEpisodes[0].VoteAverage)
	})
}

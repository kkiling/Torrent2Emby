package themoviedb

import (
	"time"
)

// Language represents supported languages
type Language string

const (
	LanguageRU Language = "ru-ru"
	LanguageEN Language = "en-en"
)

// Image contains URLs for different image sizes
type Image struct {
	W92      string
	W154     string
	W185     string
	W342     string
	W500     string
	W780     string
	Original string
}

// MovieShort contains basic movie information
type MovieShort struct {
	ID            uint64
	Title         string
	OriginalTitle string
	Overview      string
	PosterPath    *Image
	ReleaseDate   time.Time
	VoteAverage   float64
	VoteCount     int
	Popularity    float64
}

// Movie contains detailed movie information
type Movie struct {
	MovieShort
	BackdropPath     *Image
	Budget           int64
	Genres           []string
	ImdbID           string
	OriginCountry    []string
	OriginalLanguage string
	Revenue          int64
	Runtime          int
	Status           string
	Tagline          string
}

// TVShowShort contains basic TV show information
type TVShowShort struct {
	ID           uint64
	Name         string
	OriginalName string
	Overview     string
	PosterPath   *Image
	FirstAirDate time.Time
	VoteAverage  float64
	VoteCount    int
	Popularity   float64
}

// TVShow contains detailed TV show information
type TVShow struct {
	TVShowShort
	BackdropPath     *Image
	Genres           []string
	LastAirDate      time.Time
	NextEpisodeToAir time.Time
	NumberOfEpisodes int
	NumberOfSeasons  int
	OriginCountry    []string
	Seasons          []Season
	Status           string
	Tagline          string
	Type             string
}

// Season contains TV show season information
type Season struct {
	ID           int
	AirDate      string
	EpisodeCount int
	Name         string
	Overview     string
	PosterPath   *Image
	SeasonNumber int
	VoteAverage  float64
}

// Episode contains TV show episode information
type Episode struct {
	ID int
	// Дата выхода
	AirDate time.Time
	// Номер эпизода в сезоне
	EpisodeNumber int
	// Какой то тип сезона (standart)
	EpisodeType string
	// Наименование эпизода
	Name string
	// Описание эпизода
	Overview string
	// Продолжительность эпизода (секунды)
	Runtime int
	// Превью эпизода
	StillPath *Image
	// Средний рейтинг эпизода
	VoteAverage float64
	// Количество оценок
	VoteCount int
}

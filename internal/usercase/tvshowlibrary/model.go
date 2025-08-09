package tvshowlibrary

import "time"

type Image struct {
	ID       string
	W342     string
	Original string
}

// TVShowShort базовая информация о сериале
type TVShowShort struct {
	ID           uint64
	Name         string
	OriginalName string
	Overview     string
	Poster       *Image
	FirstAirDate time.Time
	VoteAverage  float64
	VoteCount    uint32
	Popularity   float64
}

// TVShow расширенная информация о сериале
type TVShow struct {
	TVShowShort
	Backdrop         *Image
	Genres           []string
	LastAirDate      time.Time
	NextEpisodeToAir time.Time
	NumberOfEpisodes uint32
	NumberOfSeasons  uint32
	OriginCountry    []string
	Status           string
	Tagline          string
	Type             string
	Seasons          []Season
}

// Season базовая информация о сезоне сериала
type Season struct {
	ID           uint64
	AirDate      time.Time
	EpisodeCount uint32
	Name         string
	Overview     string
	Poster       *Image
	SeasonNumber uint8
	VoteAverage  float64
}

// Episode информация о эпизоде сезона
type Episode struct {
	ID uint64
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
	Still *Image
	// Средний рейтинг эпизода
	VoteAverage float64
	// Количество оценок
	VoteCount int
}

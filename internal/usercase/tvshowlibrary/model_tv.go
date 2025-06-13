package tvshowlibrary

import "time"

type TVShowMetadata struct {
	// Время добавления в библиотеку
	// Время последнего обновления
	// Наличие раздачи
	// Флаг указывающий на то что раздача когда то была добавлена и скачена
}

// TVShowShort базовая информация о сериале
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

// TVShow расширенная информация о сериале
type TVShow struct {
	TVShowShort
	BackdropPath     *Image
	Genres           []string
	LastAirDate      time.Time
	NextEpisodeToAir time.Time
	NumberOfEpisodes int
	NumberOfSeasons  int
	OriginCountry    []string
	Status           string
	Tagline          string
	Type             string
	Seasons          []Season
}

// Season базовая информация о сезоне сериала
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

// Episode информация о эпизоде сезона
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

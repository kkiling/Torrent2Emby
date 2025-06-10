package search

import "time"

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
}

// TVShowSearch расширенная информация о сериале для поиска
type TVShowSearch struct {
	TVShow
	Seasons []SeasonShort
}

// TvShowLibrary информация о сериале для библиотеки
type TvShowLibrary struct {
	TVShow
	Seasons []Season
	// Время добавления
	// Имеются ли раздачи
}

// SeasonShort базовая информация о сезоне сериала
type SeasonShort struct {
	ID           int
	AirDate      string
	EpisodeCount int
	Name         string
	Overview     string
	PosterPath   *Image
	SeasonNumber int
	VoteAverage  float64
}

// Season расширенная информация о сезоне сериала
type Season struct {
	SeasonShort
	Episodes []Episode
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

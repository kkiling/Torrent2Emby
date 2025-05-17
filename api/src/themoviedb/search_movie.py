from dataclasses import dataclass
from enum import Enum
from typing import List, Optional
import requests
import math


class Language(Enum):
    RU = "ru"
    EN = "en"

    def to_api_format(self) -> str:
        return f"{self.value}-{self.value.upper()}"


@dataclass
class SearchQuery:
    query: str
    language: Language
    page: int
    per_page: int

    def __post_init__(self):
        if len(self.query) < 3:
            raise ValueError("Query must be at least 3 characters long")
        if self.page < 1:
            raise ValueError("Page must be greater than or equal to 1")
        if not (1 <= self.per_page <= 20):
            raise ValueError("PerPage must be between 1 and 20")


@dataclass
class Movie:
    adult: bool
    backdrop_path: Optional[str]
    genre_ids: List[int]
    id: int
    original_language: str
    original_title: str
    overview: str
    popularity: float
    poster_path: Optional[str]
    release_date: str
    title: str
    video: bool
    vote_average: float
    vote_count: int

@dataclass
class TVShow:
    adult: bool
    backdrop_path: Optional[str]
    genre_ids: List[int]
    id: int
    origin_country: List[str]
    original_language: str
    original_name: str
    overview: str
    popularity: float
    poster_path: Optional[str]
    first_air_date: str
    name: str
    vote_average: float
    vote_count: int

@dataclass
class MovieSearchResponse:
    page: int
    total_pages: int
    total_results: int
    results: List[Movie]

@dataclass
class TVShowSearchResponse:
    page: int
    total_pages: int
    total_results: int
    results: List[TVShow]



class TheMovieDBAPI:
    def __init__(self, api_key: str, base_url: str = "https://api.themoviedb.org/3", max_pages: int = 5):
        self.api_key = api_key
        self.base_url = base_url
        self.max_pages = max_pages

    def __search_movie_unsorted(self, params: SearchQuery) -> MovieSearchResponse:
        """
        Search for movies using The Movie Database API.

        Args:
            params: MovieSearchQuery object containing search parameters

        Returns:
            MovieSearchResponse object containing search results

        Raises:
            ValueError: If the API request fails
            requests.RequestException: If there's a network error
        """
        url = f"{self.base_url}/search/movie"

        api_params = {
            "api_key": self.api_key,
            "query": params.query,
            "language": params.language.to_api_format(),
            "page": params.page,
        }

        try:
            response = requests.get(url, params=api_params)
            response.raise_for_status()
            data = response.json()

            # Преобразуем результаты в типизированные объекты
            movies = [
                Movie(
                    adult=item["adult"],
                    backdrop_path=item.get("backdrop_path"),
                    genre_ids=item.get("genre_ids", []),
                    id=item["id"],
                    original_language=item["original_language"],
                    original_title=item["original_title"],
                    overview=item["overview"],
                    popularity=item.get("popularity", 0),
                    poster_path=item.get("poster_path"),
                    release_date=item.get("release_date", ""),
                    title=item["title"],
                    video=item.get("video", False),
                    vote_average=item.get("vote_average", 0.0),
                    vote_count=item.get("vote_count", 0.0),
                )
                for item in data["results"]
            ]

            return MovieSearchResponse(
                page=data["page"],
                total_pages=data["total_pages"],
                total_results=data["total_results"],
                results=movies
            )

        except requests.RequestException as e:
            raise ValueError(f"Failed to search movies: {str(e)}")

    def search_movie(self, params: SearchQuery) -> MovieSearchResponse:
        """
        Search for movies with sorting by popularity.
        Fetches multiple pages and returns sorted results based on popularity.

        Args:
            params: MovieSearchQuery object containing search parameters

        Returns:
            MovieSearchResponse object with sorted results

        Raises:
            ValueError: If the API request fails
        """
        # Получаем первую страницу для определения общего количества страниц
        initial_response = self.__search_movie_unsorted(SearchQuery(
            query=params.query,
            language=params.language,
            page=1,
            per_page=params.per_page # Не участвует тут
        ))

        # Определяем, сколько страниц нужно получить
        pages_to_fetch = min(
            initial_response.total_pages,
            self.max_pages
        )

        # Собираем все фильмы
        all_movies: List[Movie] = initial_response.results

        # Получаем остальные страницы
        for page in range(2, pages_to_fetch + 1):
            page_response = self.__search_movie_unsorted(SearchQuery(
                query=params.query,
                language=params.language,
                page=page,
                per_page=params.per_page # Не участвует тут
            ))
            all_movies.extend(page_response.results)

        # Сортируем все результаты по popularity по убыванию
        sorted_movies = sorted(all_movies, key=lambda x: x.popularity, reverse=True)

        # Вычисляем индексы для пагинации
        start_idx = (params.page - 1) * params.per_page
        end_idx = start_idx + params.per_page
        paginated_movies = sorted_movies[start_idx:end_idx]

        # Вычисляем общее количество страниц для отсортированных результатов
        total_results = len(sorted_movies)
        total_pages = math.ceil(total_results / params.per_page)

        return MovieSearchResponse(
            page=params.page,
            total_pages=total_pages,
            total_results=total_results,
            results=paginated_movies
        )

    def __search_tv_unsorted(self, params: SearchQuery) -> TVShowSearchResponse:
        """
        Search for TV shows using The Movie Database API.

        Args:
            params: TVShowSearchQuery object containing search parameters

        Returns:
            TVShowSearchResponse object containing search results

        Raises:
            ValueError: If the API request fails
            requests.RequestException: If there's a network error
        """
        url = f"{self.base_url}/search/tv"

        api_params = {
            "api_key": self.api_key,
            "query": params.query,
            "language": params.language.to_api_format(),
            "page": params.page,
        }

        try:
            response = requests.get(url, params=api_params)
            response.raise_for_status()
            data = response.json()

            # Преобразуем результаты в типизированные объекты
            tv_shows = [
                TVShow(
                    adult=item.get("adult", False),
                    backdrop_path=item.get("backdrop_path"),
                    genre_ids=item.get("genre_ids", []),
                    id=item["id"],
                    origin_country=item.get("origin_country", []),
                    original_language=item["original_language"],
                    original_name=item["original_name"],
                    overview=item["overview"],
                    popularity=item.get("popularity", 0),
                    poster_path=item.get("poster_path"),
                    first_air_date=item.get("first_air_date", ""),
                    name=item["name"],
                    vote_average=item.get("vote_average", 0.0),
                    vote_count=item.get("vote_count", 0),
                )
                for item in data["results"]
            ]

            return TVShowSearchResponse(
                page=data["page"],
                total_pages=data["total_pages"],
                total_results=data["total_results"],
                results=tv_shows
            )

        except requests.RequestException as e:
            raise ValueError(f"Failed to search TV shows: {str(e)}")

    def search_tv(self, params: SearchQuery) -> TVShowSearchResponse:
        """
        Search for TV shows with sorting by popularity.
        Fetches multiple pages and returns sorted results based on popularity.

        Args:
            params: TVShowSearchQuery object containing search parameters

        Returns:
            TVShowSearchResponse object with sorted results

        Raises:
            ValueError: If the API request fails
        """
        # Получаем первую страницу для определения общего количества страниц
        initial_response = self.__search_tv_unsorted(SearchQuery(
            query=params.query,
            language=params.language,
            page=1,
            per_page=params.per_page
        ))

        # Определяем, сколько страниц нужно получить
        pages_to_fetch = min(
            initial_response.total_pages,
            self.max_pages
        )

        # Собираем все сериалы
        all_shows: List[TVShow] = initial_response.results

        # Получаем остальные страницы
        for page in range(2, pages_to_fetch + 1):
            page_response = self.__search_tv_unsorted(SearchQuery(
                query=params.query,
                language=params.language,
                page=page,
                per_page=params.per_page
            ))
            all_shows.extend(page_response.results)

        # Сортируем все результаты по popularity по убыванию
        sorted_shows = sorted(all_shows, key=lambda x: x.popularity, reverse=True)

        # Вычисляем индексы для пагинации
        start_idx = (params.page - 1) * params.per_page
        end_idx = start_idx + params.per_page
        paginated_shows = sorted_shows[start_idx:end_idx]

        # Вычисляем общее количество страниц для отсортированных результатов
        total_results = len(sorted_shows)
        total_pages = math.ceil(total_results / params.per_page)

        return TVShowSearchResponse(
            page=params.page,
            total_pages=total_pages,
            total_results=total_results,
            results=paginated_shows
        )

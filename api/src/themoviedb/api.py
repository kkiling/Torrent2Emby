from typing import List, Optional
import requests
import math
from src.themoviedb.model import (SearchQuery, MovieShort, TVShowShort,
                                  MovieSearchResponse, TVShowSearchResponse,
                                  Movie, Language, Image)


class TheMovieDBAPI:
    def __init__(self, api_key: str,
                 base_api_url: str = "https://api.themoviedb.org/3",
                 base_img_url: str = "https://image.tmdb.org/t/p",
                 max_pages: int = 5):
        self.api_key = api_key
        self.base_api_url = base_api_url
        self.base_img_url = base_img_url
        self.max_pages = max_pages

    def __get_image(self, url: str) -> Optional[Image]:
        if url == "":
            return None
        return Image(
            w92=f"{self.base_img_url}/w92{url}",
            w154=f"{self.base_img_url}/w154{url}",
            w185=f"{self.base_img_url}/w185{url}",
            w342=f"{self.base_img_url}/w342{url}",
            w500=f"{self.base_img_url}/w500{url}",
            w780=f"{self.base_img_url}/w780{url}",
            original=f"{self.base_img_url}/original{url}"
        )

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
        url = f"{self.base_api_url}/search/movie"

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
                MovieShort(
                    id=item["id"],
                    original_title=item["original_title"],
                    overview=item["overview"],
                    poster_path=self.__get_image(item.get("poster_path", "")),
                    release_date=item.get("release_date", ""),
                    title=item["title"],
                    vote_average=item.get("vote_average", 0.0),
                    vote_count=item.get("vote_count", 0),
                    popularity=item.get("popularity", 0),
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
        all_movies: List[MovieShort] = initial_response.results

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
        url = f"{self.base_api_url}/search/tv"

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
                TVShowShort(
                    id=item["id"],
                    original_name=item["original_name"],
                    overview=item["overview"],
                    poster_path=self.__get_image(item.get("poster_path", "")),
                    first_air_date=item.get("first_air_date", ""),
                    name=item["name"],
                    vote_average=item.get("vote_average", 0.0),
                    vote_count=item.get("vote_count", 0),
                    popularity=item.get("popularity", 0),
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
        all_shows: List[TVShowShort] = initial_response.results

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

    def get_movie_info(self, movie_id: int, language: Language) -> Movie:
        """
        Get detailed information about a movie by its ID.

        Args:
            movie_id: The ID of the movie
            language: Language code in format 'iso-639-1-ISO-3166-1' (default: "ru-RU")

        Returns:
            Movie object containing detailed movie information

        Raises:
            ValueError: If the API request fails
            requests.RequestException: If there's a network error
        """
        url = f"{self.base_api_url}/movie/{movie_id}"

        api_params = {
            "api_key": self.api_key,
            "language": language
        }

        try:
            response = requests.get(url, params=api_params)
            response.raise_for_status()
            data = response.json()

            # Преобразуем жанры в объекты Genre
            genres = [g["name"] for g in data["genres"]]

            return Movie(
                backdrop_path=self.__get_image(data.get("backdrop_path", "")),
                budget=data["budget"],
                genres=genres,
                id=data["id"],
                imdb_id=data["imdb_id"],
                origin_country=data.get("origin_country", []),
                original_language=data["original_language"],
                original_title=data["original_title"],
                overview=data["overview"],
                popularity=data["popularity"],
                poster_path=self.__get_image(data.get("poster_path", "")),
                release_date=data["release_date"],
                revenue=data["revenue"],
                runtime=data["runtime"],
                status=data["status"],
                tagline=data["tagline"],
                title=data["title"],
                vote_average=data["vote_average"],
                vote_count=data["vote_count"]
            )

        except requests.RequestException as e:
            raise ValueError(f"Failed to get movie details: {str(e)}")
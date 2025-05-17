from enum import Enum
from dataclasses import dataclass
from typing import List, Optional
from src.utils import print_class

class Language(Enum):
    RU = "ru"
    EN = "en"
    def to_api_format(self) -> str:
        return f"{self.value}-{self.value.upper()}"

@dataclass
class Image:
    w92: str
    w154: str
    w185: str
    w342: str
    w500: str
    w780: str
    original: str
    def __str__(self):
        return print_class(self)

@dataclass
class MovieShort:
    id: int
    title: str
    original_title: Image
    overview: str
    poster_path: Optional[Image]
    release_date: str
    vote_average: float
    vote_count: int
    popularity: int
    def __str__(self):
        return print_class(self)

@dataclass
class Movie:
    backdrop_path: Optional[Image]
    budget: int
    genres: List[str]
    id: int
    imdb_id: str
    origin_country: List[str]
    original_language: str
    original_title: str
    overview: str
    popularity: float
    poster_path: Optional[Image]
    release_date: str
    revenue: int
    runtime: int
    status: str
    tagline: str
    title: str
    vote_average: float
    vote_count: int

    def __str__(self):
        return print_class(self)

@dataclass
class TVShowShort:
    id: int
    name: str
    original_name: str
    overview: str
    poster_path: Optional[Image]
    first_air_date: str
    vote_average: float
    vote_count: int
    popularity: int

    def __str__(self):
        return print_class(self)

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
class MovieSearchResponse:
    page: int
    total_pages: int
    total_results: int
    results: List[MovieShort]


@dataclass
class TVShowSearchResponse:
    page: int
    total_pages: int
    total_results: int
    results: List[TVShowShort]

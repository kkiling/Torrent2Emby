import os
from dotenv import load_dotenv
from src.themoviedb.search_movie import SearchQuery, Language, TheMovieDBAPI

def main():
    load_dotenv()
    api_key = os.getenv("THE_MOVIE_DB_API_KEY")

    api = TheMovieDBAPI(api_key=api_key)

    search_params = SearchQuery(
        query="Атака",
        language=Language.RU,
        page=1,
        per_page=10
    )

    try:
        results = api.search_tv(search_params)
        print(f"page: {results.page}: total_results {results.total_results}")
        for movie in results.results:
            print(f"Name: {movie.name}")
            print(f"FirstAirDate: {movie.first_air_date}")
            print(f"Popularity: {movie.popularity}")
            print(f"Vote: {movie.vote_average}")
            print("---")
    except ValueError as e:
        print(f"Error: {e}")


if __name__ == '__main__':
    main()
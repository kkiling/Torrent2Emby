from src.themoviedb import SearchQuery, Language, TheMovieDBAPI
from src.options import Options


def main():
    opt = Options.from_env()
    api = TheMovieDBAPI(api_key=opt.the_movie_db_api_key)

    search_params = SearchQuery(
        query="Однажды",
        language=Language.RU,
        page=1,
        per_page=3
    )

    try:
        results = api.search_movie(search_params)
        print(f"page: {results.page}: total_results {results.total_results}")
        for movie in results.results:
            print(movie)
            print("---")

        print("First movie info")
        movie_info = api.get_movie_info(results.results[0].id, Language.RU)
        print(movie_info)
    except ValueError as e:
        print(f"Error: {e}")

if __name__ == '__main__':
    main()
from src.options import Options
from src.themoviedb import SearchQuery, Language, TheMovieDBAPI

def main():
    opt = Options.from_env()
    api = TheMovieDBAPI(api_key=opt.the_movie_db_api_key)

    search_params = SearchQuery(
        query="Атака",
        language=Language.RU,
        page=1,
        per_page=3
    )

    try:
        results = api.search_tv(search_params)
        print(f"page: {results.page}: total_results {results.total_results}")
        for tv in results.results:
            print(tv)
            print("---")

    except ValueError as e:
        print(f"Error: {e}")


if __name__ == '__main__':
    main()
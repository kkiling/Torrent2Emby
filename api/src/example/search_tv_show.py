from src.options import Options
from src.themoviedb import SearchQuery, Language, TheMovieDBAPI
from src.utils import pretty_print


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
            pretty_print(tv)
            print("---")

        print("First TV info")
        tv_info = api.get_tv_show(results.results[0].id, Language.RU)
        #pretty_print(tv_info)

        print("Episodes for season 1:")
        episodes = api.get_season_episodes(tv_info.id, tv_info.seasons[1].season_number, Language.RU)
        pretty_print(episodes)

    except ValueError as e:
        print(f"Error: {e}")


if __name__ == '__main__':
    main()
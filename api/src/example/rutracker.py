from src.options import Options
from src.rutracker.api import RutrackerAPI
from src.utils import pretty_print


def main():
    opt = Options.from_env()
    api = RutrackerAPI(
        username=opt.rutracker_user_name,
        password=opt.rutracker_password,
        cookies_dir=opt.rutracker_cookie_dir,
    )

    try:
        results = api.search_torrents("клинок рассекающий демонов")
        print(f"page: {results.page}: total_results {results.total_results}")
        for torrent in results.results:
            pretty_print(torrent)
            print("---")

        if not results.results:
            return
        magnet_link = api.get_magnet_link(results.results[0].href)
        pretty_print(magnet_link)

    except ValueError as e:
        print(f"Error: {e}")

if __name__ == '__main__':
    main()
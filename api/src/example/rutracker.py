from src.options import Options
from src.rutracker.api import RutrackerAPI
from src.utils import pretty_print


def main():
    opt = Options.from_env()
    api = RutrackerAPI(username="joefantor", password="AAL0X", cookies_dir="./rutracker")

    try:
        results = api.search_torrents("клинок рассекающий демонов бесконечный поезд")
        print(f"page: {results.page}: total_results {results.total_results}")
        for torrent in results.results:
            pretty_print(torrent)
            print("---")

        magnet_link = api.get_magnet_link(results.results[0].href)
        pretty_print(magnet_link)

    except ValueError as e:
        print(f"Error: {e}")

if __name__ == '__main__':
    main()

    # E197E8652333EBF7525D124104C5B6EB49848CA3
from src.options import Options
from src.qbittorrent.api import QBittorrentAPI
from src.utils import pretty_print


def main():
    opt = Options.from_env()
    # Пример использования
    api = QBittorrentAPI(
        username=opt.qbittorrent_username,
        password=opt.qbittorrent_password,
        cookies_dir=opt.qbittorrent_cookie_dir,
        base_api_url=opt.qbittorrent_api_url
    )

    torrent_hash = "e197e8652333ebf7525d124104c5b6eb49848ca3"
    try:
        info = api.get_torrent_info(torrent_hash)
        if info:
            pretty_print(info)
    except Exception as e:
        print(f"Ошибка: {e}")


if __name__ == '__main__':
    main()
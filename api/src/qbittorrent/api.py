import os
import pickle
from typing import Optional
import requests
from src.qbittorrent.model import TorrentAddOptions,TorrentInfo, TorrentState
from datetime import datetime

def map_torrent_info(data: dict) -> TorrentInfo:
    """
    Преобразует словарь в объект TorrentInfo
    
    Args:
        data: Словарь с данными торрента из API
        
    Returns:
        TorrentInfo: Объект с информацией о торренте
    """
    return TorrentInfo(
        hash=data['hash'],
        name=data['name'],
        category=data['category'],
        tags=data['tags'],
        content_path=data['content_path'],
        state=TorrentState(data['state']),
        added_on=datetime.fromtimestamp(data['added_on']),
        completion_on=datetime.fromtimestamp(data['completion_on']),
        eta=data['eta'],
        amount_left=data['amount_left'],
        completed=data['completed'],
        downloaded=data['downloaded'],
        uploaded=data['uploaded'],
        size=data['size'],
        total_size=data['total_size'],
        progress=data['progress'],
        dl_speed=data['dlspeed'],
        up_speed=data['upspeed']
    )


class QBittorrentAPI:
    def __init__(
            self,
            username: str,
            password: str,
            cookies_dir: str,
            base_api_url: str,
    ):
        self.username = username
        self.password = password
        self.cookies_dir = cookies_dir
        self.base_api_url = base_api_url.rstrip("/")
        self.api_url = f"{self.base_api_url}/api/v2"

    def __save_cookies(self, session, filename):
        os.makedirs(self.cookies_dir, exist_ok=True)
        with open(os.path.join(self.cookies_dir, filename), 'wb') as f:
            pickle.dump(session.cookies, f)

    def __load_cookies(self, session, filename):
        try:
            with open(os.path.join(self.cookies_dir, filename), 'rb') as f:
                session.cookies.update(pickle.load(f))
            return True
        except FileNotFoundError:
            return False

    def __check_auth(self, session) -> bool:
        """Check if current session is authenticated"""
        response = session.get(f"{self.api_url}/app/version")
        return response.status_code == 200

    def __try_login(self, session) -> bool:
        """Attempt to login and save cookies"""
        login_data = {
            'username': self.username,
            'password': self.password,
        }
        resp = session.post(f'{self.api_url}/auth/login', data=login_data)
        if resp.ok and self.__check_auth(session):
            self.__save_cookies(session, "qbittorrent_cookies.pkl")
            return True
        return False

    def __login(self, session):
        # Try to load and verify existing cookies
        if self.__load_cookies(session, "qbittorrent_cookies.pkl"):
            if not self.__check_auth(session):
                # If cookies are expired, try to re-login
                if not self.__try_login(session):
                    raise Exception("Re-authentication failed")
        else:
            # If no cookies exist, try to login
            if not self.__try_login(session):
                raise Exception("Authentication failed")

    def add_torrent(self, options: TorrentAddOptions) -> bool:
        """
        Add a new torrent to qBittorrent

        Args:
            options: TorrentAddOptions object containing torrent parameters

        Returns:
            bool: True if torrent was added successfully

        Raises:
            Exception: If the request fails
        """
        session = requests.Session()
        self.__login(session)

        data = {
            'urls': options.urls,
            'savepath': options.save_path,
        }
        if options.category:
            data['category'] = options.category
        if options.tags:
            data['tags'] = ','.join(options.tags)
        if options.paused:
            data['paused'] = 'true'

        response = session.post(
            f"{self.api_url}/torrents/add",
            data=data
        )

        if not response.ok:
            raise Exception(f"Failed to add torrent: {response.status_code} {response.text}")

        return True

    def get_torrent_info(self, torrent_hash: str) -> Optional[TorrentInfo]:
        """
        Получает информацию о торренте по его хешу

        Args:
            torrent_hash: Хеш торрента

        Returns:
            TorrentInfoOriginal: Информация о торренте или None, если торрент не найден

        Raises:
            Exception: При ошибке запроса
        """
        session = requests.Session()
        self.__login(session)

        response = session.get(
            f"{self.api_url}/torrents/info",
            params={'hashes': torrent_hash.lower()}
        )

        if not response.ok:
            raise Exception(f"Failed to get torrent info: {response.status_code} {response.text}")

        torrents = response.json()
        if not torrents:
            return None

        return map_torrent_info(torrents[0])
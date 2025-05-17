from typing import List
from src.rutracker.model import TorrentInfoResponse, TorrentInfo, MagnetInfo
import requests
from bs4 import BeautifulSoup
import pickle
from urllib.parse import quote
import os

class RutrackerAPI:
    def __init__(self,
                 username: str,
                 password: str,
                 cookies_dir: str,
                 base_api_url: str = "https://rutracker.org/forum/",
                 ):
        self.username = username
        self.password = password
        self.cookies_dir = cookies_dir
        self.base_api_url = base_api_url

    def __save_cookies(self, session, filename):
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
        """Проверяет валидность авторизации"""
        # Делаем запрос к странице, которая доступна только авторизованным пользователям
        test_url = f"{self.base_api_url}tracker.php"
        response = session.get(test_url)
        # Проверяем, что в ответе нет формы логина
        return 'login.php' not in response.url

    def __try_login(self, session) -> bool:
        """Выполняет логин и сохраняет куки"""
        login_data = {
            'login_username': self.username,
            'login_password': self.password,
            'login': 'вход'
        }
        resp = session.post(f'{self.base_api_url}login.php', data=login_data)
        if resp.ok and self.__check_auth(session):
            self.__save_cookies(session, "rutracker_cookies.pkl")
            return True
        return False

    def __login(self, session):
        # Пытаемся загрузить и проверить существующие куки
        if self.__load_cookies(session, "rutracker_cookies.pkl"):
            if not self.__check_auth(session):
                # Если куки протухли, пытаемся перелогиниться
                if not self.__try_login(session):
                    raise Exception("Re-authentication failed")
        else:
            # Если куков нет, пытаемся залогиниться
            if not self.__try_login(session):
                raise Exception("Authentication failed")

    def search_torrents(self, query: str) -> TorrentInfoResponse:
        session = requests.Session()
        self.__login(session)

        # Формируем поисковый запрос
        search_url = f"{self.base_api_url}tracker.php?nm={quote(query)}"
        response = session.get(search_url)

        if not response.ok:
            raise Exception(f"Search failed with status code: {response.status_code}")

        soup = BeautifulSoup(response.text, 'html.parser')
        table = soup.find('table', {'id': 'tor-tbl'})

        if not table:
            return TorrentInfoResponse(results=[], page=1,total_results=0, total_pages=0)

        # Проверяем первую строку на наличие сообщения "Не найдено"
        rows = table.find('tbody').find_all('tr')
        if rows and len(rows) > 0:
            first_row = rows[0]
            cols = first_row.find_all('td')
            if cols and cols[0].get_text(strip=True) == 'Не найдено':
                return TorrentInfoResponse(results=[], page=1, total_results=0, total_pages=0)

        results: List[TorrentInfo] = []

        # Проходим по всем строкам таблицы (кроме заголовка)
        for row in table.find('tbody').find_all('tr'):
            cols = row.find_all('td')

            # Извлекаем данные из каждого столбца
            forum = cols[2].get_text(strip=True)
            title_link = cols[3].find('a')
            title = title_link.get_text(strip=True)
            href = f"{self.base_api_url}{title_link['href']}"
            author = cols[4].get_text(strip=True)
            size = cols[5].get_text(strip=True)
            seeds = cols[6].get_text(strip=True)
            leeches = cols[7].get_text(strip=True)
            downloads = cols[8].get_text(strip=True)
            added_date = cols[9].get_text(strip=True)

            # Добавляем в результаты
            results.append(TorrentInfo(
                title=title,
                href=href,
                forum=forum,
                author=author,
                size=size,
                seeds=seeds,
                leeches=leeches,
                downloads=downloads,
                added_date=added_date
            ))

        results = sorted(results, key=lambda x: x.downloads, reverse=True)
        return TorrentInfoResponse(results=results, page=1,total_results=0, total_pages=0)


    def get_magnet_link(self, topic_url: str) -> str:
        """
        Получает магнет-ссылку для указанного топика

        Args:
            topic_url: Полный URL топика или ID топика

        Returns:
            str: Магнет-ссылка

        Raises:
            MagnetLinkNotFoundException: Если магнет-ссылка не найдена
            RutrackerAPIException: При других ошибках
        """
        session = requests.Session()
        self.__login(session)

        try:
            if not topic_url.startswith(self.base_api_url):
                raise Exception("Invalid topic URL")

            response = session.get(topic_url)

            if not response.ok:
                raise Exception(f"Failed to get topic page: {response.status_code}")

            soup = BeautifulSoup(response.text, 'html.parser')
            magnet_link = soup.find('a', class_='magnet-link')

            if not magnet_link or 'href' not in magnet_link.attrs:
                raise Exception("Magnet link not found on the page")

            magnet = magnet_link['href']

            # Извлекаем хеш из title атрибута
            hash_value = magnet_link.get('title', '')
            if not hash_value:
                # Если хеш не найден в title, пробуем извлечь его из магнет-ссылки
                import re
                hash_match = re.search(r'btih:([A-F0-9]{40})', magnet, re.IGNORECASE)
                hash_value = hash_match.group(1) if hash_match else ''

            return MagnetInfo(
                magnet=magnet,
                hash=hash_value
            )

        except requests.RequestException as e:
            raise Exception(f"Network error: {str(e)}")
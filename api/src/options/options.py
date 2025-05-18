import os
from dataclasses import dataclass
from dotenv import load_dotenv

# Constant for environment variable name
THE_MOVIE_DB_API_KEY = "THE_MOVIE_DB_API_KEY"

RUTRACKER_USERNAME = "RUTRACKER_USERNAME"
RUTRACKER_PASSWORD = "RUTRACKER_PASSWORD"
RUTRACKER_COOKIE_DIR = "RUTRACKER_COOKIE_DIR"

QBITTORRENT_USERNAME = "QBITTORRENT_USERNAME"
QBITTORRENT_PASSWORD = "QBITTORRENT_PASSWORD"
QBITTORRENT_COOKIE_DIR = "QBITTORRENT_COOKIE_DIR"
QBITTORRENT_API_URL = "QBITTORRENT_API_URL"


def load_env(key: str):
    value = os.getenv(key)
    if not value:
        raise ValueError(f"{key} environment variable must be set")
    return value

@dataclass
class Options:
    """Configuration class project options."""
    the_movie_db_api_key: str
    rutracker_user_name: str
    rutracker_password: str
    rutracker_cookie_dir: str
    qbittorrent_username: str
    qbittorrent_password: str
    qbittorrent_cookie_dir: str
    qbittorrent_api_url: str

    @classmethod
    def from_env(cls) -> 'Options':
        """
        Load project options from environment variables using dotenv.


        Returns:
            Options: Configuration instance with project options

        Raises:
            ValueError: If THE_MOVIE_DB_API_KEY environment variable is missing
        """
        # Load environment variables from .env file
        load_dotenv()

        return cls(
            the_movie_db_api_key=load_env(THE_MOVIE_DB_API_KEY),
            rutracker_user_name=load_env(RUTRACKER_USERNAME),
            rutracker_password=load_env(RUTRACKER_PASSWORD),
            rutracker_cookie_dir=load_env(RUTRACKER_COOKIE_DIR),
            qbittorrent_username=load_env(QBITTORRENT_USERNAME),
            qbittorrent_password=load_env(QBITTORRENT_PASSWORD),
            qbittorrent_cookie_dir=load_env(QBITTORRENT_COOKIE_DIR),
            qbittorrent_api_url=load_env(QBITTORRENT_API_URL)
        )
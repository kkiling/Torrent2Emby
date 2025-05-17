import os
from dataclasses import dataclass
from dotenv import load_dotenv

# Constant for environment variable name
THE_MOVIE_DB_API_KEY = "THE_MOVIE_DB_API_KEY"


@dataclass
class Options:
    """Configuration class project options."""
    the_movie_db_api_key: str

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

        the_movie_db_api_key = os.getenv(THE_MOVIE_DB_API_KEY)
        if not the_movie_db_api_key:
            raise ValueError(f"{THE_MOVIE_DB_API_KEY} environment variable must be set")

        return cls(the_movie_db_api_key=the_movie_db_api_key)
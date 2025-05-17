from src.themoviedb.search_movie import MovieSearchQuery, Language, TheMovieDBAPI

def main():
    # Example usage
    api = TheMovieDBAPI(api_key="4d5d76d2d405435cbc3d4c0e68374ff5")

    search_params = MovieSearchQuery(
        query="Начало",
        language=Language.RU,
        page=1,
        per_page=10
    )

    try:
        results = api.search_movie(search_params)
        print(f"page: {results.page}: total_results {results.total_results}")
        for movie in results.results:
            print(f"Title: {movie.title}")
            print(f"ReleaseDate: {movie.release_date}")
            print(f"Popularity: {movie.popularity}")
            print(f"Vote: {movie.vote_average}")
            print("---")
    except ValueError as e:
        print(f"Error: {e}")

if __name__ == '__main__':
    main()
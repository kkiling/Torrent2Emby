
# TheMovieDB

TheMovieDB API — это RESTful API для работы с базой данных фильмов. 
Позволяет получать информацию о фильмах, актерах, рейтингах, 
а также управлять избранным (добавлять, удалять, получать список).

[Api документация](https://developer.themoviedb.org/reference/intro/getting-started)


### Как получить API Key
Зарегистрируйтесь на [The Movie Database (TMDb)](https://www.themoviedb.org/signup).

Перейдите в [настройки API](https://www.themoviedb.org/settings/api).

Запросите API Key, заполнив форму (укажите тип использования, например, "Personal").
Получите ключ в разделе "API" вашего аккаунта.

Пример API Key:

```plaintext
abcdef1234567890abcdef1234567890
```

Используйте этот ключ в заголовке запроса:
```http
Authorization: Bearer YOUR_API_KEY
```

Или в параметрах запроса:
```http
api_key YOUR_API_KEY
```

### Лимиты
Сейчас лимит [50 запросов в секунду](https://developer.themoviedb.org/docs/rate-limiting). 
В случае превышения лимита получаем ошибку `429`. 

### Запросы
* Поиск фильмов - [Search/Movie](https://developer.themoviedb.org/reference/search-movie)
* Поиск Сериалов - [Search/Movie](https://developer.themoviedb.org/reference/search-tv)
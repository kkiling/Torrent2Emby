# Torrent2Emby

# Установка mkvmerge
https://linuxconfig.org/installation-of-mkvtoolnix-matroska-tools-on-ubuntu-linux

```bash
# Добавляем видео (трек 0), русскую аудио (трек 0 в своём файле), английские субтитры (трек 0 в своём файле)
mkvmerge -o "film.mkv" "video.mp4" \
  "audio_ru.ac3" --language 0:rus --track-name 0:"Русский" --default-track 0:yes \
  "subs_en.srt" --language 0:eng --track-name 0:"English" --default-track 0:no
```
Здесь все 0: корректны, потому что каждый файл (аудио, субтитры) содержит только одну дорожку.

### Поиск фильма
- Поиск фильма по названию
- Поиск сериала по названию
- Добавление фильма в библиотеку
- Добавление фильма в сериал
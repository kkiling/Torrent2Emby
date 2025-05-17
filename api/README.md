### Создайте виртуальное окружение:
```
python3 -m venv .venv
source .venv/bin/activate
```
### Установите проект в режиме разработки
 это способ сделать ваш Python-пакет доступным для импорта в системе, при этом продолжая редактировать его код без 
 необходимости переустанавливать после каждого изменения.
```
sudo apt install python3-poetry
pip install -e .
```

### Run api server
```bash
poetry run start
```

### Add packages
```
pip install poetry
poetry add пакет  # автоматически добавляет в pyproject.toml
```
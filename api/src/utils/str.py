import json
from dataclasses import asdict
from datetime import datetime


def pretty_print(obj):
    # Сериализация с обработкой datetime
    json_str = json.dumps(
        asdict(obj),
        default=lambda o: o.isoformat() if isinstance(o, datetime) else str(o),
        indent=4
    )
    print(json_str)


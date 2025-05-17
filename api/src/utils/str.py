from dataclasses import is_dataclass
from typing import Any


def pretty_str(
        obj: Any,
        indent: int = 0,
        _seen: set[int] | None = None,
        is_list_item: bool = False
) -> str:
    """Рекурсивно форматирует объект с правильными отступами и цветами."""
    if _seen is None:
        _seen = set()

    #if id(obj) in _seen:
    #    return "\033[90m<...>\033[0m"  # Серый цвет для циклических ссылок
    _seen.add(id(obj))

    # Цвета
    COL_KEY = "\033[36m"  # Голубой для ключей
    COL_STR = "\033[33m"  # Желтый для строк
    COL_NUM = "\033[35m"  # Фиолетовый для чисел
    COL_BOOL = "\033[32m"  # Зеленый для bool
    COL_RESET = "\033[0m"  # Сброс цвета
    COL_LIST = "\033[90m"  # Серый для элементов списка

    if is_dataclass(obj):
        fields = vars(obj)
        lines = []
        for k, v in fields.items():
            lines.append(
                f"{'  ' * indent}{COL_KEY}{k}:{COL_RESET} {pretty_str(v, indent + 1, _seen)}"
            )
        return "\n".join(lines)

    elif isinstance(obj, (list, tuple)):
        if not obj:
            return f"{COL_LIST}[]{COL_RESET}"

        lines = []
        for i, item in enumerate(obj):
            # Для первого элемента не добавляем лишний отступ
            prefix = "" if i == 0 and is_list_item else ("" * indent)
            lines.append(
                f"{prefix}{COL_LIST}-{COL_RESET} {pretty_str(item, indent + 1, _seen, True)}"
            )
            # Добавляем разделитель между элементами (кроме последнего)
            if i < len(obj) - 1:
                lines.append(f"{' ' * indent}{COL_LIST}│{COL_RESET}")
        return "\n".join(lines)

    elif isinstance(obj, str):
        return f"{COL_STR}{obj}{COL_RESET}"
    elif isinstance(obj, (int, float)):
        return f"{COL_NUM}{obj}{COL_RESET}"
    elif isinstance(obj, bool):
        return f"{COL_BOOL}{obj}{COL_RESET}"
    else:
        return str(obj)


def pretty_print(obj, indent=0):
    print(pretty_str(obj, indent))
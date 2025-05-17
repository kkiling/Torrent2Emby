from dataclasses import asdict

def print_class(self):
    return "\n".join(f"{key}: {value}" for key, value in asdict(self).items())
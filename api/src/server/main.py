from src.options import Options

def main():
    opt = Options.from_env()
    print(opt)

if __name__ == '__main__':
    main()
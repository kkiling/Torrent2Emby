module github.com/kkiling/torrent2emby

go 1.24.4

require (
	github.com/PuerkitoBio/goquery v1.10.3
	github.com/go-playground/validator/v10 v10.26.0
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/kkiling/goplatform v0.0.2
	github.com/mattn/go-sqlite3 v1.14.23
	github.com/samber/lo v1.50.0
	github.com/stretchr/testify v1.10.0
	golang.org/x/net v0.40.0
	github.com/kkiling/statemachine v1.0.0
)

require (
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/kkiling/statemachine v1.0.0 => ../pet_projects/statemachine
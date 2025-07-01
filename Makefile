.PHONY: mocks
mocks:
	@echo "\n --- 🤡 Create Mocks --- \n"
	go generate ./...

.PHONY: .test
.test:
	@echo "\n --- 🧪 Run project tests --- \n"
	go test ./...

.PHONY: test
test: .test

.PHONY: format
format:
	@echo "\n --- 🧪 Start format imports --- \n"
	smartimports_ -local "github.com/kkiling/torrent2emby/" -path . #-exclude pkg/mocks

.PHONY: bin-deps
bin-deps:
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/pav5000/smartimports/cmd/smartimports@v0.2.0

.PHONY: build run release clean

BIN := exam-transcriber
GO_DIR := src
CMD := ./cmd/exam-transcriber

build:
	cd $(GO_DIR) && go build -o ../$(BIN) $(CMD)

run:
	cd $(GO_DIR) && go run $(CMD)

release:
	./scripts/build_release.sh

clean:
	rm -f $(BIN)

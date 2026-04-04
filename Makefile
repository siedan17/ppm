.PHONY: run build clean generate

run:
	go run ./cmd/ppm --port 8080

build:
	go build -o ppm ./cmd/ppm

generate:
	sqlc generate

clean:
	rm -f ppm ppm.db

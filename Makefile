.PHONY: run build clean

run:
	go run ./cmd/ppm --port 8080

build:
	go build -o ppm ./cmd/ppm

clean:
	rm -f ppm ppm.db

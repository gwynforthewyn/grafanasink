build:
	CGO_ENABLED=1 go build

test:
	CGO_ENABLED=1 go test

.PHONY: prototype

prototype: build
	docker build -t playtechnique/gsync:0.0.0 .

protorun: prototype
	docker run -it --name gsync --rm -p 3000:3000 playtechnique/gsync:0.0.0
build:
	mkdir -p bin
	go build -o ./bin ./...

sdk: build
	./bin/pulumi-gen-short-io language nodejs

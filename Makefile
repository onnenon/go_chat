build-server:
	docker build -f build/Dockerfile.server --rm --no-cache -t 319-server:latest .

build-client:
	docker build -f build/Dockerfile.client --rm --no-cache -t 319-client:latest .

run-server:
	docker run -p 9000:9000 --name=go-server 319-server:latest

docker-clean:
	docker rm -f 319-server || true
	docker rm -f 319-client || true
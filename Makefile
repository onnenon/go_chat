build-server:
	docker build -f build/server/Dockerfile --rm --no-cache -t koozie/go-chat:server .

build-client:
	docker build -f build/client/Dockerfile --rm --no-cache -t koozie/go-chat:client .

run-server:
	docker run -p 9000:9000 --name=go-server koozie/go-chat:server

docker-clean:
	docker rm -f koozie/go-chat:server || true
	docker rm -f koozie/go-chat:client || true
.PHONY: docker-build-api
docker-build-api:
	docker build -f api.Dockerfile -t go-zero-gin-api:${version} .
	# make remove-none-images # 构建完之后删除镜像

.PHONY: remove-none-images
remove-none-images:
	docker images | awk '$$1=="<none>"' | awk '{print $$3}' | xargs docker rmi

.PHONY: docker-run-api
docker-run-api:
	docker rm -f go-zero-gin-api
	docker run -d --name go-zero-gin-api -p 9097:9097 \
		go-zero-gin-api:${version}

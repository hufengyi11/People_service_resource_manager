docker-build:
	docker build --tag hufengyi11/resource-manager-user-service:v0.0.1 .
.PHONY: docker-build

docker-push:
	docker push hufengyi11/resource-manager-user-service:v0.0.1
.PHONY: docker-push

docker-run:
	docker run -dp 8080:8080 hufengyi11/resource-manager-user-service:v0.0.1
.PHONY: docker-run
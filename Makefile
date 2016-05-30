.PHONY: all

build-app:
	docker-compose build app

test:
	@gb test

start-app: build-app
	docker-compose up -d app

generate: generate-dashboard
	gb generate

build-dashboard:
	docker build -t autoscale-dashboard -f src/autoscale/dashboard/Dockerfile src/autoscale/dashboard

generate-dashboard: build-dashboard
	mkdir -p src/autoscale/static
	docker run --rm -v "$$PWD/src/autoscale/static":"/src/static" autoscale-dashboard ember build -prod -o /src/static


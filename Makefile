.PHONY: build-app push-app publish-build

AUTOSCALE_HOST ?= 127.0.0.1

CURRENT_TAG = `date +"%Y%m%d%H%M%S"`

publish-build:
	@docker build -f Dockerfile.build -t do-autoscale-build . && \
		docker run --rm -it -v ~/.mc:/root/.mc  do-autoscale-build

build-app:
	@docker build -t bryanl/do-autoscale . && \
	  docker tag bryanl/do-autoscale bryanl/do-autoscale:${CURRENT_TAG}

push-app: build-app
	@docker push bryanl/do-autoscale

test:
	@gb test

run-app: build-app
	@docker-compose up -d app

generate: generate-dashboard
	@gb generate

build-dashboard:
	@docker build -t autoscale-dashboard -f src/autoscale/dashboard/Dockerfile src/autoscale/dashboard

generate-dashboard: build-dashboard
	@mkdir -p src/autoscale/static; \
		docker run --rm -v "$$PWD/src/autoscale/static":"/src/static" autoscale-dashboard ember build --environment="production" -o /src/static

ember-server:
	@cd src/autoscale/dashboard; ember server --proxy http://${AUTOSCALE_HOST}:8888

regen-mocks:
	@cd src/autoscale; mockery -all -inpkg
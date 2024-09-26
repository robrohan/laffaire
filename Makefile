.PHONY: build clean test

HASH=$(shell git log --pretty=format:'%h' -n 1)

include .env
export

# DOCKER_CONTAINER=$(REGION)-docker.pkg.dev/$(PROJECT_ID)/$(REPOSITORY)/$(PROJECT)
DOCKER_CONTAINER=$(REPOSITORY)/$(PROJECT)

# List all targets in thie file
list:
	@echo ""
	@echo "*^-. Go App Template .-^*"
	@echo ""
	@grep -B 1 '^[^#[:space:]].*:' Makefile
	@echo ""

# Install any needed libraries
install:
	go mod tidy

# Run all go unit tests
test:
	go test -v ./...

# Runs the localhost server
start:
	go run cmd/server/main.go

watch:
#	If you want "hot reloading"
#	sudo apt install entr
	find ./ | entr -sr 'go run cmd/server/main.go'

clean:
	rm -rf build

build: clean
	mkdir -p build
	CGO_ENABLED=1 GOOS=linux \
		go build -o build/server \
			-ldflags '-X main.build=$(HASH)' \
			cmd/server/main.go
	cp -R static build/
	cp -R templates build/
	cp -R migrations build/

docker_build: build
	docker ps ; \
	docker build -t $(DOCKER_CONTAINER):$(HASH) .

docker_push:
	docker ps ; \
	docker push $(DOCKER_CONTAINER):$(HASH)

#	Using a different enviroment variable set for prod
docker_run:
	docker ps ; \
	docker run --env-file=.env.production -p 8080:3000 $(DOCKER_CONTAINER):$(HASH)

#	Grab a base css that styles form elements with some basic style
fetch_base_css:
	curl https://raw.githubusercontent.com/robrohan/pho-ui/main/src/pho-ui.css > templates/pho-ui.css

google:
	gcloud auth application-default set-quota-project $(PROJECT_ID)
	gcloud config set project $(PROJECT_ID)
	gcloud services enable run.googleapis.com
	gcloud projects add-iam-policy-binding $(PROJECT_ID) \
		--member=serviceAccount:$(PROJECT_NUMBER)-compute@developer.gserviceaccount.com \
		--role=roles/cloudbuild.builds.builder

start_db:
	docker ps ; \
	docker run --name postgres \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-e POSTGRES_DB=$(DB_NAME) \
		-e PGDATA=$(DB_PATH) \
		-p 5432:5432 \
		-d postgres:latest

.PHONY: build clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello handlers/hello.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/world handlers/world.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/find_user handlers/find_user.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/add_user_todynamo handlers/add_user_todynamo.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/participate_group handlers/participate_group.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

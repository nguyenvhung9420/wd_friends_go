.PHONY: build clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/fetch_brief_followers handlers/fetch_brief_followers.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/find_user handlers/find_user.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/add_user_todynamo handlers/add_user_todynamo.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/participate_group handlers/participate_group.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/add_follower handlers/add_follower.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

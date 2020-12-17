BASE_DIR= github.com/roloum/users
BUILD_CMD= env GOOS=linux go build -ldflags="-s -w" -o
TEST_CMD= go test -timeout 30s

.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	${BUILD_CMD} bin/createUser cmd/lambda/handlers/create/main.go
	${BUILD_CMD} bin/notifyUser cmd/lambda/handlers/notify/main.go
	${BUILD_CMD} bin/activateUser cmd/lambda/handlers/activate/main.go

.PHONY: test
test:
	${TEST_CMD} ${BASE_DIR}/internal/user
# clean:
# 	rm -rf ./bin ./vendor Gopkg.lock
#
# deploy: clean build
# 	sls deploy --verbose
#
# gomodgen:
# 	chmod u+x gomod.sh
# 	./gomod.sh

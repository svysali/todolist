# Install all tools globally
.PHONY: dev-deps
dev-deps:
	cd; go get -u github.com/golang/mock/mockgen
	cd; go get -u golang.org/x/lint/golint
	# see https://github.com/kyleconroy/sqlc/issues/654
	tmp_dir=sqlc-$(shell date +'%s') && mkdir -p /tmp/$$tmp_dir && cd /tmp/$$tmp_dir && go mod init tmp && go mod edit -require github.com/kyleconroy/sqlc@v1.4.0 && go install github.com/kyleconroy/sqlc/cmd/sqlc && cd - # it has C dependencies that would be a headache using go mod vendor


.PHONY: build
build:
	go build -o ./dist/todolist ./main.go

.PHONY: sqlc
sqlc:
	sqlc generate

.PHONE: test
test:
	docker-compose --profile=test up --exit-code-from test

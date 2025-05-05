# test:
# 	go run ./cmd/ysplit < cfg.yml

build:
	go build -o ysplit .

run:
	go run main.go < cfg.yml

publish:
	git tag v0.1.0
	git push --tags

test:
	mkdir -p tmp
	go run main.go -json < test/fixtures/docs.yml > tmp/out.json
	jq -S . test/expected/docs.json > tmp/expected_sorted.json
	jq -S . tmp/out.json > tmp/out_sorted.json
	diff -u tmp/expected_sorted.json tmp/out_sorted.json

install:
	go install github.com/inem/ysplit@v0.1.0

.PHONY: test

ARGS = $(filter-out $@,$(MAKECMDGOALS))
%:
	@:

include *.mk
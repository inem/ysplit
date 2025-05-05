# test:
# 	go run ./cmd/ysplit < cfg.yml

build:
	go build -o ysplit .

run:
	go run main.go < cfg.yml

test:
	mkdir -p tmp
	go run main.go -json < test/fixtures/docs.yml > tmp/out.json
	jq -S . test/expected/docs.json > tmp/expected_sorted.json
	jq -S . tmp/out.json > tmp/out_sorted.json
	diff -u tmp/expected_sorted.json tmp/out_sorted.json

.PHONY: test

ARGS = $(filter-out $@,$(MAKECMDGOALS))
%:
	@:

include *.mk
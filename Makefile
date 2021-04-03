PHONY: run
run:
	go run .

PHONY: test
test: 
	go test ./...

PHONY: compare
compare: run
	@echo "No output means the files match"
	@diff ./generated_output.txt output.txt



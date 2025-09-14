ifeq ($(OS),Windows_NT)
    # Windows
    RM := del /s /q
    MKDIR := mkdir
else
    # Linux/Unix
    RM := rm -rf
    MKDIR := mkdir -p
endif

lint:
	golangci-lint run

test:
	$(RM) coverage
	$(MKDIR) coverage
	go test -race ./... -count=1 -p 1 -covermode=atomic -coverprofile=coverage/coverage.out

test.report: test
	go tool cover -html coverage/coverage.out -o coverage/coverage.html
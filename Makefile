include .env
export

## run: Run go app (with arg "runPath" or by default path = "main.go")
runPath ?= main.go
run:
	go run $(runPath)

## help: help command
help: Makefile
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
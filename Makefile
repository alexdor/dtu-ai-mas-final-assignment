# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

BINARY_NAME = dtu-ai-mas-final-assignment

LEVEL=levels/custom_levels/SASimple.lvl
ifdef level
        LEVEL=$(level)
endif

build:
	@bash -c 'go build -o $(BINARY_NAME) .'

start: build
	@bash -c ' java -jar server.jar -l $(LEVEL) -c "./$(BINARY_NAME)" -t 180'

race:
	@bash -c 'go build -o $(BINARY_NAME) -race . && java -jar server.jar -l $(LEVEL) -c "./$(BINARY_NAME)" -t 30'

start-gui: build
	@bash -c 'java -jar server.jar -l $(LEVEL) -c "./$(BINARY_NAME)" -t 180 -g'

start-debug: build
	@bash -c 'DEBUG=true java -jar server.jar -l $(LEVEL) -c "./$(BINARY_NAME)" -g'

profile: build
	@bash -c 'java -jar server.jar -l $(LEVEL) -c "./$(BINARY_NAME) -cpuprofile cpu.prof"'

runner: build
	@bash -c 'node runner.js -l $(LEVEL) -c "./$(BINARY_NAME)" -t 180 -i SA'

debug: build
	@bash -c "dlv attach --api-version 2 --headless --listen=:2345 `pgrep $(BINARY_NAME)` ./$(BINARY_NAME)"

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo


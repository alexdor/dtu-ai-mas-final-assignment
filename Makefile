# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

LEVEL=levels/custom_levels/SASimple.lvl
ifdef level
        LEVEL=$(level)
endif

start:
	@bash -c 'java -jar server.jar -l $(LEVEL) -c "go run ." -t 300'

start-gui:
	@bash -c 'java -jar server.jar -l $(LEVEL) -c "go run ." -t 300 -g'

start-debug:
	@bash -c 'go build . && DEBUG=true java -jar server.jar -l $(LEVEL) -c "./dtu-ai-mas-final-assignment" -g'

debug:
	@bash -c "dlv attach --api-version 2 --headless --listen=:2345 `pgrep dtu-ai-mas-final-assignment` ./dtu-ai-mas-final-assignment"

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo


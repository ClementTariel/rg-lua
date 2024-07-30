SHELL=/bin/sh

YELLOW='\033[1;93m'
BLUE='\033[1;34m'
RESET='\033[0m'

build:
	@echo ${YELLOW}[BUILD]${RESET}
	@echo ${BLUE}[INFO]${RESET} Install rgcore dependency && cd referee; go install ../rgcore
	@echo ${BLUE}[INFO]${RESET} Install player dependency && cd referee; go install ../player
	@echo ${BLUE}[INFO]${RESET} Build rg.exe && cd referee; go build -o ../rg.exe -a referee.go main.go

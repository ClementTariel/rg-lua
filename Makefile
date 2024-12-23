SHELL=/bin/sh

YELLOW='\033[1;93m'
BLUE='\033[1;34m'
RESET='\033[0m'

build:
	@echo ${YELLOW}[BUILD]${RESET}
	@echo ${BLUE}[INFO]${RESET} Install rgcore dependency for player && cd player; go install ../rgcore
	@echo ${BLUE}[INFO]${RESET} Build player.exe && cd player; go build -a -o ../player.exe player.go main.go
	@echo ${BLUE}[INFO]${RESET} Install rgcore dependency for referee && cd referee; go install ../rgcore
	@echo ${BLUE}[INFO]${RESET} Build referee.exe && cd referee; go build -a -o ../referee.exe referee.go main.go

download:
	@echo ${YELLOW}[DOWNLOAD]${RESET}
	@echo ${BLUE}[INFO]${RESET} Downloading golang packages for player && cd player ; go mod download
	@echo ${BLUE}[INFO]${RESET} Downloading golang packages for referee && cd referee ; go mod download

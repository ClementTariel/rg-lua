SHELL=/bin/sh

YELLOW='\033[1;93m'
BLUE='\033[1;34m'
RESET='\033[0m'

build:
	@echo ${YELLOW}[BUILD]${RESET}
	@echo ${BLUE}[INFO]${RESET} Build player.exe && cd player; go build -o ../player.exe cmd/main.go
	@echo ${BLUE}[INFO]${RESET} Build referee.exe && cd referee; go build -o ../referee.exe cmd/main.go
	@echo ${BLUE}[INFO]${RESET} Build matchmaker.exe && cd matchmaker; go build -o ../matchmaker.exe cmd/main.go

full-build:
	@echo ${YELLOW}[BUILD]${RESET}
	@echo ${BLUE}[INFO]${RESET} Build player.exe && cd player; go build -a -o ../player.exe cmd/main.go
	@echo ${BLUE}[INFO]${RESET} Build referee.exe && cd referee; go build -a -o ../referee.exe cmd/main.go
	@echo ${BLUE}[INFO]${RESET} Build matchmaker.exe && cd matchmaker; go build -a -o ../matchmaker.exe cmd/main.go

download:
	@echo ${YELLOW}[DOWNLOAD]${RESET}
	@echo ${BLUE}[INFO]${RESET} Downloading golang packages for player && cd player ; go mod download
	@echo ${BLUE}[INFO]${RESET} Downloading golang packages for referee && cd referee ; go mod download
	@echo ${BLUE}[INFO]${RESET} Downloading golang packages for matchmaker && cd matchmaker ; go mod download

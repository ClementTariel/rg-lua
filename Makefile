SHELL=/bin/sh

YELLOW='\033[1;93m'
BLUE='\033[1;34m'
RESET='\033[0m'

build:
	@echo ${YELLOW}[BUILD]${RESET}
	@echo ${BLUE}[INFO]${RESET} Build database init script && cd deployment; ./init.sh
	@echo ${BLUE}[INFO]${RESET} Build docker images && cd deployment; docker compose build

stop:
	@echo ${YELLOW}[STOP]${RESET}
	@echo ${BLUE}[INFO]${RESET} Stop docker images && cd deployment; docker compose down --remove-orphans

run: stop
	@echo ${YELLOW}[RUN]${RESET}
	@echo ${BLUE}[INFO]${RESET} Run docker images && cd deployment; docker compose up -d; cd ..
	@echo ${BLUE}[INFO]${RESET} Show logs from matchmaker; docker logs --follow `docker ps -aqf "name=matchmaker"`; cd ..
	
clear-db: stop
	@echo ${YELLOW}[CLEAR]${RESET}
	@echo ${BLUE}[INFO]${RESET} Stop docker images and clear db && cd deployment; docker volume prune --all --filter label=clear="rglua" -f

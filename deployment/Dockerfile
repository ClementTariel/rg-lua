FROM postgres:16.6

COPY ./bots /bots
# Copy the init.sql file inside the docker-entrypoint-initdb.d of the image
# The PostgreSql image will automatically execute all SQL files inside on startup
# if there is no volume already existing.
COPY ./deployment/init.sql /docker-entrypoint-initdb.d/

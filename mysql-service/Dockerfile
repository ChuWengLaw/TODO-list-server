FROM mysql:latest

ENV MYSQL_ROOT_PASSWORD=123456
ENV MYSQL_DATABASE=todo
ENV MYSQL_USER=user
ENV MYSQL_PASSWORD=123456

COPY database.sql /docker-entrypoint-initdb.d/
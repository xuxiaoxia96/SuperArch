#! /bin/bash

# rabbitmq
docker run -d -p 15672:15672 -p 5672:5672 --hostname my-rabbit --name some-rabbit -e RABBITMQ_DEFAULT_USER=admin -e RABBITMQ_DEFAULT_PASS=testadmin123ashore rabbitmq:3.11.1-management

# redis
docker run -itd --name superarch-redis -p 6379:6379 redis:6.0.16 --requirepass "testadmin123ashore"

# postgres
docker run -d -p 5432:5432 --name mypostgres -e POSTGRES_PASSWORD=testadmin123ashore -e POSTGRES_USER=admin -v pgdata:/var/lib/postgresql/data postgres:14
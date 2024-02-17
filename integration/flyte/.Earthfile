VERSION 0.8
FROM python:3
WORKDIR /code

build:
    COPY ./requirements.txt .
    RUN pip install -r requirements.txt
    COPY . .

run-tests:
    FROM earthly/dind:alpine-3.19-docker-25.0.2-r0
    COPY ./tests ./tests
    RUN apk update
    RUN apk add postgresql-client
    WITH DOCKER --compose docker-compose.yml --load app:latest=+docker
        RUN while ! pg_isready --host=localhost --port=5432; do sleep 1; done ;\
          docker run --network=default_python/part6_default app python3 ./tests/test_db_connection.py
    END

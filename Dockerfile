# Собираем в гошке
FROM golang:1.21.5 as build

ENV CODE_DIR /Applications/calendar_service/

WORKDIR ${CODE_DIR}

COPY ./ ${CODE_DIR}

RUN apt-get update
RUN apt-get -y install postgresql-client


# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go build -o car_catalog ./cmd/main.go
RUN goose -dir ./migrations postgres "postgresql://goose:password@127.0.0.1:8092/go_migrations?sslmode=disable" up

CMD [ "./car_catalog" ]
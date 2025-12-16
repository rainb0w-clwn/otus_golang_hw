# Собираем в гошке
FROM golang:1.23 AS build

WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 go test -c -o /integration-tests ./tests/integration

# На выходе тонкий образ
FROM alpine:3.23

LABEL ORGANIZATION="LittleBig"
LABEL SERVICE="sender"
LABEL MAINTAINERS="mikhail.anikin15@gmail.com"

COPY --from=build /integration-tests /app/integration-tests

CMD /app/integration-tests

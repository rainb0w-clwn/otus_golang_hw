# Собираем в гошке
FROM golang:1.23 AS build

ENV BIN_FILE /opt/scheduler/scheduler-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/scheduler/*

# На выходе тонкий образ
FROM alpine:3.23

LABEL ORGANIZATION="LittleBig"
LABEL SERVICE="scheduler"
LABEL MAINTAINERS="mikhail.anikin15@gmail.com"

ENV BIN_FILE /opt/scheduler/scheduler-app
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/scheduler/config.yml
COPY ./configs/config.yml ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}

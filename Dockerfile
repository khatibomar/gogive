FROM golang:1.17-alpine as build-env
 
ENV APP_NAME gogive-backend 
ENV CMD_PATH cmd/api/

COPY . $GOPATH/src/$APP_NAME
WORKDIR $GOPATH/src/$APP_NAME
 
RUN CGO_ENABLED=0 go build -v -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH
 
FROM alpine:3.14
 
ENV APP_NAME gogive-backend 
ARG GO_PORT=4000
ARG GO_ENV="development"


# https://stackoverflow.com/questions/35560894/is-docker-arg-allowed-within-cmd-instruction/35562189#35562189
ENV GO_PORT=${GO_PORT}
ENV GO_ENV=${GO_ENV}
ENV GOGIVE_DB_DSN=${GO_DSN}
ENV APP_PATH=./${APP_NAME}

COPY --from=build-env /$APP_NAME .
 
# Start app
CMD ${APP_PATH} -port ${GO_PORT} -env ${GO_ENV}

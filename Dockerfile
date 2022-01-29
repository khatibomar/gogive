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
ARG GO_DB_MAX_OPEN_CONNS=25
ARG GO_DB_MAX_IDLE_CONNS=25
ARG GO_DB_MAX_IDLE_TIME=15
ARG GO_RATE_LIMITER_RPS=2
ARG GO_RATE_LIMITER_BURST=4
ARG GO_RATE_LIMITER_ENABLED=true
ARG GO_SMTP_HOST=mailhog
ARG GO_SMTP_PORT=1025
ARG GO_SMTP_USERNAME=""
ARG GO_SMTP_PASSWORD=""
ARG GO_SMTP_SENDER="GoGive <no-reply@gogive.com>"
ARG CORS_TRUSTED_ORIGINS="*"

# https://stackoverflow.com/questions/35560894/is-docker-arg-allowed-within-cmd-instruction/35562189#35562189
ENV GO_PORT=${GO_PORT}
ENV GO_ENV=${GO_ENV}
ENV GO_DB_MAX_IDLE_TIME=${GO_DB_MAX_IDLE_TIME}
ENV GO_DB_MAX_IDLE_CONNS=${GO_DB_MAX_IDLE_CONNS}
ENV GO_DB_MAX_OPEN_CONNS=${GO_DB_MAX_OPEN_CONNS}
ENV GO_RATE_LIMITER_RPS=${GO_RATE_LIMITER_RPS}
ENV GO_RATE_LIMITER_BURST=${GO_RATE_LIMITER_BURST}
ENV GO_RATE_LIMITER_ENABLED=${GO_RATE_LIMITER_ENABLED}
ENV GO_SMTP_HOST=${GO_SMTP_HOST}
ENV GO_SMTP_PORT=${GO_SMTP_PORT}
ENV GO_SMTP_USERNAME=${GO_SMTP_USERNAME}
ENV GO_SMTP_PASSWORD=${GO_SMTP_PASSWORD}
ENV GO_SMTP_SENDER=${GO_SMTP_SENDER}
ENV CORS_TRUSTED_ORIGINS=${CORS_TRUSTED_ORIGINS}

ENV APP_PATH=./${APP_NAME}

COPY --from=build-env /$APP_NAME .
 
# Start app
CMD ${APP_PATH} -port ${GO_PORT} -env ${GO_ENV} -db-max-open-conns ${GO_DB_MAX_OPEN_CONNS} -db-max-idle-conns ${GO_DB_MAX_IDLE_CONNS} -db-max-idle-time ${GO_DB_MAX_IDLE_TIME} -limiter-rps=${GO_RATE_LIMITER_RPS} -limiter-burst=${GO_RATE_LIMITER_BURST} -limiter-enabled=${GO_RATE_LIMITER_ENABLED} -smtp-host=${GO_SMTP_HOST} -smtp-port=${GO_SMTP_PORT} -smtp-username=${GO_SMTP_USERNAME} -smtp-password=${GO_SMTP_PASSWORD} -smtp-sender="${GO_SMTP_SENDER}" -cors-trusted-origins="${CORS_TRUSTED_ORIGINS}"

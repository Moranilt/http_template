FROM hb.zonatelecom.ru/dhp/library/golang:1.20 AS builder

COPY . /src
WORKDIR /src

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /src/bin/sms_otp

COPY ./migrations /src/bin/migrations

FROM hb.zonatelecom.ru/dhp/library/alpine:latest

ARG PRODUCTION
ENV PRODUCTION=$PRODUCTION

ARG PORT
ENV PORT=$PORT

ARG TRACER_URL
ENV TRACER_URL=$TRACER_URL

ARG TRACER_NAME
ENV TRACER_NAME=$TRACER_NAME

COPY --from=builder /src/bin /src/bin
WORKDIR /src/bin

EXPOSE $PORT
ENTRYPOINT ["/src/bin/sms_otp"]



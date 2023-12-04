FROM golang:1.21.4-alpine3.18 AS builder

WORKDIR /usr/src/app

COPY go.mod /usr/src/app
COPY cmd /usr/src/app/cmd

RUN go build -v ./cmd/satis-hook

###

FROM composer/satis

COPY --from=builder /usr/src/app/satis-hook /usr/local/bin

COPY .docker/docker-entrypoint.sh /

ENV PATH="${PATH}:/satis/bin"

EXPOSE 80

CMD ["/docker-entrypoint.sh"]

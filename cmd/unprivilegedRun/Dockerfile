FROM golang:alpine as builder

WORKDIR /app

COPY ./ ./

RUN go mod init unpiviligedRun
RUN go mod download
RUN go build -o ./bin/unprivrun main.go

FROM alpine as runner

COPY --from=builder /app/bin/unprivrun /home/

WORKDIR /home

RUN mkdir /home/jail

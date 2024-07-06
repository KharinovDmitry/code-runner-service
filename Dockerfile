FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY ./ ./

RUN go build -o ./bin/code-runner-service cmd/app/main.go

FROM alpine AS runner

RUN apk add docker
RUN apk add g++
RUN apk add --no-cache musl-dev


COPY --from=builder /app/bin/code-runner-service /
COPY --from=builder /app/config/config.yaml config.yaml

COPY --from=builder /app/cmd/unprivilegedRun /unprivilegedRun
COPY --from=builder /app/internal/executor/python /python

ENTRYPOINT ["./code-runner-service", "-path=config.yaml"]
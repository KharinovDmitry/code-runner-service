FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY ./ ./

RUN go build -o ./bin/code-runner-service cmd/app/main.go

FROM alpine AS runner

COPY --from=builder /app/bin/code-runner-service /
COPY --from=builder /app/config/config.yaml config.yaml

ENTRYPOINT ["./article-service -path=config.yaml"]
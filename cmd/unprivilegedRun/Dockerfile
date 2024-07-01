FROM golang:alpine as builder

WORKDIR /app

COPY ./ ./

RUN go mod download
RUN go build -o ./bin/unprivrun main.go

FROM alpine as runner

COPY --from=builder /app/bin/unprivrun /home/unprivrun

WORKDIR /home

RUN mkdir /home/jail && \
    mkdir /home/jail/tmp && \
    cp -a /bin /etc /lib /sbin /usr /home/unprivrun /home/jail


RUN chroot jail adduser -h / -D unprivuser && \
    chroot jail chown -R root:root / && \
    chroot jail chmod 755 /






FROM unpivileged_run

RUN apk add g++ && \
    apk add --no-cache musl-dev

ENV FILE_NAME=""

CMD ["./unprivrun", "3000", "1000", "./tmp/${FILE_NAME}"]
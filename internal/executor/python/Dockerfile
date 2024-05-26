FROM unpivileged_run

RUN apk add python3

RUN mkdir /home/jail/usr && \
    mkdir /home/jail/usr/bin && \
    mkdir /home/jail/usr/lib && \
    mkdir /home/jail/lib && \
    cp /lib/ld-musl-x86_64.so.1 /home/jail/lib/ld-musl-x86_64.so.1 && \
    cp /usr/lib/libpython3.12.so.1.0 /home/jail/usr/lib/libpython3.12.so.1.0 && \
    cp /usr/bin/python3 /home/jail/usr/bin/python3 && \
    cp -r /usr/lib/python3.12 /home/jail/usr/lib/python3.12

ENV FILE_NAME=""

CMD ["./unprivrun", "3000", "1000", "python3", "tmp/${FILE_NAME}"]
FROM golang:1.16.4-alpine3.12 AS build

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk add --no-cache gcc g++ linux-headers  

COPY ./clipboard_archive_backend /clipboard_archive_backend

RUN cd /clipboard_archive_backend \
    && go env -w CGO_ENABLED=0 \
    && go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go build

FROM alpine:3.12
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk add --no-cache tzdata 
CMD [ "/clipboard_archive_backend" ]
COPY --from=build /clipboard_archive_backend/clipboard_archive_backend  /clipboard_archive_backend

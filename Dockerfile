FROM golang:1.16.4-alpine3.12 AS build

# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories 
RUN apk add --no-cache gcc g++ linux-headers  

COPY . .

RUN cd . \
    && go env -w CGO_ENABLED=0 \
    && go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.io,direct \
    && go build

FROM alpine:3.12
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories 
RUN apk add --no-cache tzdata 
CMD [ "/clipboard_archive_backend" ]
COPY --from=build /clipboard_archive_backend/clipboard_archive_backend  /clipboard_archive_backend

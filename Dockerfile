FROM golang:1.19-alpine AS build

WORKDIR /clipboard_archive
COPY . .
RUN    go env -w CGO_ENABLED=0 \
    && go build -v \
    && go test ./... -cover -v


FROM alpine:latest

RUN apk add --no-cache tzdata
ENTRYPOINT [ "/clipboard_archive" ]
COPY --from=build /clipboard_archive/clipboard_archive  /clipboard_archive

FROM golang:1.19-alpine AS build

WORKDIR /clipboard_archive
COPY . .
RUN    go env -w CGO_ENABLED=0 \
    && go build -v \
    && go test ./... -cover -v


FROM alpine:latest

ENTRYPOINT [ "/clipboard_archive" ]
RUN apk add --no-cache tzdata
COPY --from=build /clipboard_archive/clipboard_archive  /clipboard_archive

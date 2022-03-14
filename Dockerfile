FROM golang AS build

COPY . /

RUN    cd / \
    && go env -w CGO_ENABLED=0 \
    && go env -w GO111MODULE=on \
    && go build

FROM alpine
CMD [ "/clipboard_archive_backend" ]
RUN apk add --no-cache tzdata 
COPY --from=build /clipboard_archive_backend  /clipboard_archive_backend

FROM golang AS build
WORKDIR /
COPY . /
RUN    go env -w CGO_ENABLED=0 \
    && go env -w GO111MODULE=on 
RUN    go build -v 

FROM alpine
RUN apk add --no-cache tzdata 
CMD [ "/clipboard_archive" ]
COPY --from=build /clipboard_archive  /clipboard_archive

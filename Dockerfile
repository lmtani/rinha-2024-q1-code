FROM golang:1.22-alpine3.18

# disable cgo to avoid gcc requirement bug
ENV CGO_ENABLED=0

# Fix entr watch on Mac/Windows
ENV ENTR_INOTIFY_WORKAROUND 1

RUN apk --no-cache add entr ca-certificates

WORKDIR /app

EXPOSE 8080

CMD ["./bin/boot-dev.sh"]

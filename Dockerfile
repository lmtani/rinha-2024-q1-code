FROM golang:1.22-alpine3.18 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

FROM alpine:3.18 as production
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]

FROM golang:1.22-alpine3.18 as development
# disable cgo to avoid gcc requirement bug
ENV CGO_ENABLED=0
# Fix entr watch on Mac/Windows
ENV ENTR_INOTIFY_WORKAROUND 1
RUN apk --no-cache add entr ca-certificates
WORKDIR /app
EXPOSE 8080
CMD ["./bin/boot-dev.sh"]

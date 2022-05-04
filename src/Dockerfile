# Build it
FROM golang AS builder

WORKDIR /app

COPY /main.go .
COPY /go.mod .
COPY /go.sum .
# go get -d -v &&
RUN go mod download
RUN CGO_ENABLED=0 go build -installsuffix cgo -o snoopy .

# Final
FROM alpine:latest  as final

WORKDIR /app
COPY --from=builder /app/snoopy .

EXPOSE 9080
EXPOSE 2112

CMD ["./snoopy"]
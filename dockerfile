FROM golang:alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod ./

# Download all dependencies.
RUN go mod download 

# Copy the source from the current directory to the Working Directory inside the container
ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

CMD ["./main"]
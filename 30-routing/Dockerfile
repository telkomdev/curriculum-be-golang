# builder images from golang 1.19
FROM golang:1.19-alpine as builder

# work directory
WORKDIR /app

# copy all files in directory
COPY . .

# build binary
RUN go build -o /app/docker-go

# small images to run binary
FROM alpine:3.16.2

# work directory
WORKDIR /app

# copy binary to workdir
COPY --from=builder /app/docker-go ./

# moving binary
RUN mv ./docker-go /usr/local/bin/docker-go

# run binary
CMD ["/usr/local/bin/docker-go"]

# builder images from golang 1.19
FROM golang:1.19-alpine as builder

# work directory
WORKDIR /data

# copy all files in directory
COPY . .

# build binary
RUN go build -o ./bin/app

# small images to run binary
FROM alpine:3.16.2

# work directory
WORKDIR /data

# set env for dodumentation
ENV DOCS_LOCATION=./docs

# copy binary to workdir
COPY --from=builder /data/bin/app /data/app
COPY --from=builder /data/docs/swagger.yaml /data/docs/swagger.yaml

# run binary
CMD ["/data/app"]

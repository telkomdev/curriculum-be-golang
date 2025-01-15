FROM golang:1.19-alpine

#Change app directory
WORKDIR /app

#Copy main.go
COPY main.go .

#Copy list dependency
COPY go.mod .

#Validasi dependency
RUN go mod tidy

#Build binary
RUN go build -o /docker-go

#Run binary
CMD [ "/docker-go" ]
FROM golang:latest
COPY .. ./
RUN go build -o ./ ./app/server/Server.go
CMD ["./Server"]
EXPOSE 8080
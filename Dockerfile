FROM golang:latest
RUN mkdir /go/test
WORKDIR /go/test
COPY ./ ./
RUN go build -o ./ ./cmd/main.go
#RUN go run ./cmd/main.go
CMD ["./main"]
EXPOSE 8080
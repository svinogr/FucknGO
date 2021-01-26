FROM golang:latest
COPY ./ ./
#RUN apt-get update && apt-get install postgresql -y
RUN go build -o ./ ./app/server/Server.go
CMD ["./Server"]
EXPOSE 8080
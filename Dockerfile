FROM golang:latest
COPY ./ ./
#RUN apt-get update && apt-get install postgresql -y
RUN go build -o ./ ./app/server/main.go
CMD ["./main"]
EXPOSE 8080
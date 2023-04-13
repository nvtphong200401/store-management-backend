FROM golang:latest

WORKDIR /

COPY . .

RUN go build cmd/main.go


EXPOSE 8080

CMD ["./main"]
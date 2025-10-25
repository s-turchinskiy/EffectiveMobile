FROM golang:1.24
WORKDIR /app
COPY . .
RUN go mod tidy && cd ./cmd/subscriptions && go build -o ../../app && cd ../../
CMD ["/app/app"]
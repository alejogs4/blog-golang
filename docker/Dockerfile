FROM golang:1.15.8-alpine

WORKDIR /go/src/go-blog
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build .

CMD ["blog" ]
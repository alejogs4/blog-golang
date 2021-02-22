FROM golang:1.15.8-alpine

WORKDIR /go/src/go-blog
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080
# CMD [ "blog" ]
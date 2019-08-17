FROM golang:alpine as build
COPY . /go/src/github.com/dontfollowsean/slackspot
WORKDIR /go/src/github.com/dontfollowsean/slackspot
RUN go build -o app
CMD ["./app"]

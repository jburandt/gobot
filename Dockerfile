FROM golang:1.9

WORKDIR /go/src/app
COPY . .

RUN go get . 

ENV SLACK_TOKEN=<SLACK TOKEN>
CMD ["app"]

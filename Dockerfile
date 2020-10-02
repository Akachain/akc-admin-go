FROM golang:1.14

RUN mkdir -p /app

WORKDIR /app

ADD . /app

RUN go mod tidy & go get
RUN go mod vendor & go build

EXPOSE 4001

CMD ["/app/akc-admin-go"]

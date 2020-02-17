FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go get -d -v
RUN go install -v
RUN go get github.com/codegangsta/gin

CMD ["gin", "-i", "run" ,"main.go"]
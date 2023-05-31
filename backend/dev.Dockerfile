FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go get -d -v
RUN go install -v
RUN go install github.com/codegangsta/gin@latest

CMD ["gin", "-i", "run" ,"main.go"]
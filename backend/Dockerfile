FROM golang:latest AS build

RUN apt-get update && apt-get install ca-certificates && update-ca-certificates

WORKDIR /go/src/app
COPY . .

RUN go get -d -v
RUN go install -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch AS run
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/app/main .
CMD ["./main"]
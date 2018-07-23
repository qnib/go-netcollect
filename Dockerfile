FROM golang:alpine AS build

WORKDIR /go/src/github.com/qnib/go-netcollect
COPY ./main.go .
COPY vendor/ ./vendor/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-netcollect .

FROM alpine

COPY --from=build /go/src/github.com/qnib/go-netcollect/go-netcollect /go-netcollect
ENTRYPOINT ["/go-netcollect"]
FROM golang:1.16.0 as builder

WORKDIR /go/src/github.com/davebehr1/micro/vessel-service

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .


FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/davebehr1/micro/vessel-service/vessel-service .

CMD ["./vessel-service"]
FROM golang:1.16.0 as builder

WORKDIR /go/src/github.com/davebehr1/micro/user-cli

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o user-cli -a -installsuffix cgo .


FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/davebehr1/micro/user-cli/user-cli .

ENTRYPOINT ["./user-cli"]

CMD ["./user-cli"]

FROM golang:1.19-alpine3.16

RUN apk add --no-cache git

WORKDIR /app/ssp

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -o ./out/ssp .
EXPOSE 3000
CMD ["./out/ssp"]
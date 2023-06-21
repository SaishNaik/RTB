FROM golang:1.19-alpine3.16

RUN apk add --no-cache git

WORKDIR /app/dsp

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -o ./out/dsp .
EXPOSE 3002
CMD ["./out/dsp"]
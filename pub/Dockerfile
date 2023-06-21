FROM golang:1.19-alpine3.16

RUN apk add --no-cache git

WORKDIR /app/pub

COPY go.mod .
#COPY go.sum .
#RUN #go mod download
COPY . .

RUN go build -o ./out/pub .
EXPOSE 3001
CMD ["./out/pub"]
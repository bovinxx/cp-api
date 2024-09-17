FROM golang:latest AS BUILD

WORKDIR /app

COPY . .

RUN go build -o main ./processor/main.go 

FROM golang:latest

WORKDIR /app

COPY --from=BUILD /app/main .
COPY --from=BUILD /app/processor/config.yml .

CMD ["./main"]
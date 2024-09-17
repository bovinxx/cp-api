FROM golang:latest AS BUILD

#RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /build

COPY . .

WORKDIR /build/api/cmd/app

RUN go build -o main

FROM golang:latest

WORKDIR /app

COPY --from=BUILD /build/api/cmd/app/config.yml .
COPY --from=BUILD /build/api/cmd/app/main .

CMD ["./main"]
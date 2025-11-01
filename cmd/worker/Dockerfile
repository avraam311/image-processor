FROM golang:alpine AS build_base

WORKDIR /worker

COPY . .

RUN go mod download

RUN go build -o worker ./cmd/main.go

FROM alpine AS runner

COPY --from=build_base ./worker .
COPY ./config ./config
COPY ./.env .

CMD ["./worker"]
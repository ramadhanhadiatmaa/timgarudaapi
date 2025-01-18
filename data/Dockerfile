FROM golang:1.22.8 AS build

WORKDIR /app

COPY go.mod go.sum ./
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o data-garuda

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/data-garuda .

EXPOSE 6711

CMD [ "./data-garuda" ]
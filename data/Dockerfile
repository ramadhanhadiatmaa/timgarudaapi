FROM golang:1.22.8 AS build

WORKDIR /app

COPY go.mod go.sum ./
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o data-garuda-dev

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/data-garuda-dev .

EXPOSE 6411

CMD [ "./data-garuda-dev" ]
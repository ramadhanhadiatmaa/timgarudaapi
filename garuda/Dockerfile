FROM golang:1.22.8 AS build

WORKDIR /app

COPY go.mod go.sum ./
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o garu-garuda

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/garu-garuda .

EXPOSE 6713

CMD [ "./garu-garuda" ]
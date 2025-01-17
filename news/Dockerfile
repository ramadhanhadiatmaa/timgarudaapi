FROM golang:1.22.8 AS build

WORKDIR /app

COPY go.mod go.sum ./
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o news-garuda

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/news-garuda .

EXPOSE 6712

CMD [ "./news-garuda" ]
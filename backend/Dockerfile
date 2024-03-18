FROM golang:1.21-alpine as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY  . .
RUN go build -o main .
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/main .
COPY ../.env ./.env
EXPOSE 3000
CMD [ "./main" ]
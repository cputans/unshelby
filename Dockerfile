FROM golang:1.22-alpine as build

ENV GOOS linux

WORKDIR /app
RUN apk update && apk add --update gcc musl-dev alpine-sdk
ADD . .
RUN go build -o unshelby main.go 

FROM alpine
COPY --from=build /app/unshelby .
CMD ["./unshelby"]
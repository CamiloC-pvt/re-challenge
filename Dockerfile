# syntax=docker/dockerfile:1

FROM golang:1.22.6

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN ls

RUN go build -o /re-challenge

EXPOSE 3001

CMD [ "/re-challenge", "--port", "3001"]
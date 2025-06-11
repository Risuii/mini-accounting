#build stage
FROM golang:1.20.4-alpine AS build

ARG branch=develop
ENV env $branch

RUN apk add --no-cache bash

# Create app directory
RUN mkdir -p /usr/src/app/bridgtl-sus-be-jagat-raya
WORKDIR /usr/src/app/bridgtl-sus-be-jagat-raya

# Copying source files
COPY . /usr/src/app/bridgtl-sus-be-jagat-raya

RUN go install -tags 'sqlserver' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install github.com/google/wire/cmd/wire@latest
RUN go mod tidy

RUN go mod download

RUN go build -o /bridgtl-sus-be-jagat-raya

CMD [ "/bridgtl-sus-be-jagat-raya" ]
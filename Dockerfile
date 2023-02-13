FROM golang:1.20.0
LABEL maintainer="saqib.yawar@devopsways.com"

# Working directory
WORKDIR /app

# Fetch dependencies
COPY ./go.mod ./go.sum ./
RUN go mod download

# Copy source
COPY . .

EXPOSE 9101
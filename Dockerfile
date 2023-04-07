# syntax=docker/dockerfile:1

##
## STEP 1 - BUILD
##

# Base image
FROM golang:1.20.0 AS build
LABEL maintainer="saqib.yawar@devopsways.com"

# Working directory
WORKDIR /app

# Fetch dependencies
COPY ./go.mod ./

# Download go modules and dependencies
RUN go mod download

# Copy source files
COPY . .

# compile application
RUN go build -o ./bin/exporter main.go

##
## STEP 2 - DEPLOY
##
FROM scratch

WORKDIR /

COPY --from=build bin/exporter /exporter

EXPOSE 9101

ENTRYPOINT ["/exporter"]
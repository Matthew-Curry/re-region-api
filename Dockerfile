FROM golang:latest

LABEL maintainer="Matthew Curry <matt.curry56@gmail.com>"

WORKDIR /app

# handle dependencies using go.mod and go.sum
COPY go.mod .
COPY go.sum .
RUN go mod download

# copy remaining source code
RUN mkdir /app/src
COPY src/ /app/src

# specify needed env variables. Secret variables passed in as args during build
ARG RE_REGION_API_USER
ARG RE_REGION_API_PASSWORD
ARG RE_REGION_DB
ARG DB_PORT
ARG DB_HOST

ENV RE_REGION_API_USER $RE_REGION_API_USER
ENV RE_REGION_API_PASSWORD $RE_REGION_API_PASSWORD
ENV RE_REGION_DB $RE_REGION_DB
ENV DB_PORT $DB_PORT
ENV DB_HOST $DB_HOST
ENV PORT 8080

# build the app, expose the port, the CMD runs the build executable
RUN cd src; go build

EXPOSE 8080

CMD ["./src/src"]
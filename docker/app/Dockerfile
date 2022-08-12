# intermediate container to build the code
FROM golang:latest AS builder
LABEL maintainer="Matthew Curry <matt.curry56@gmail.com>"
WORKDIR /app
# handle dependencies using go.mod and go.sum
COPY go.mod .
COPY go.sum .
RUN go mod download
# copy remaining source code
RUN mkdir /app/src
COPY src/ /app/src
# build the app
RUN cd src; CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo

# move artifact into final container and run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/src/src .
COPY --from=builder /app/src/static/docs.yml .

EXPOSE 8080
CMD ["./src"]
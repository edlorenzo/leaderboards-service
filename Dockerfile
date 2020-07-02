# Start from the latest golang base image
FROM golang:latest
RUN apt-get update

# Add local directory
ADD . /go/src

# Change working dir
WORKDIR /go/src

# Get GO Dependencies
RUN go get -d -v .
RUN go install -v .

# Build app
RUN go build -o main .

# RUN endpoint
CMD ["/go/src/main"]
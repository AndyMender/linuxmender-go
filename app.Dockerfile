# fetch base image
# NOTE: use 'latest' tag to pull in the current golang image
# all regular images are Debian-based
FROM golang:1.13

ENV CGO_ENABLED=1

# install additional dependecies
RUN apt install -y gcc libsqlite3-dev

# make sure WORKDIR is inside Go's source dir
WORKDIR /go/src/linuxmender

COPY . .

# get all dependencies based on their usage
RUN go get -d -v ./...
RUN go install -v ./...

# prepare a Docker entry point
COPY docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]
# fetch base image
# NOTE: use 'latest' tag to pull in the current golang image
# all standard Go images are Debian-based
FROM golang:1.13

# make sure WORKDIR is inside Go's source dir
WORKDIR /go/src/linuxmender

# move tthe repo into the container
COPY . .

# get all dependencies based on their usage
RUN go get -d -v ./...
RUN go install -v ./...

# prepare a Docker entry point
COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT ["docker-entrypoint.sh"]

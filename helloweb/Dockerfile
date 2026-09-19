FROM golang:1.26-alpine AS build-stage
# The FROM setting sets the base image used for the following commands.
# We are using an alpine base image with go version indicated installed.
# Alpine is small and fast. See the tags for available images here:
#     https://hub.docker.com/_/golang
# -----------------
# The following lines copy the contents of local directories into
# /go/src/* directories in the image. We assume we are in the webserver
# directory.
ADD digits    /go/src/helloweb/digits
ADD hyphenate /go/src/helloweb/hyphenate
ADD longDiv   /go/src/helloweb/longDiv
ADD longmult  /go/src/helloweb/longmult
ADD prime     /go/src/helloweb/prime
ADD quadratic /go/src/helloweb/quadratic
ADD series    /go/src/helloweb/series
ADD webserver /go/src/helloweb/webserver
ADD go.mod    /go/src/helloweb/go.mod

# The go install command behaves almost identically to go build, but
# instead of leaving the executable in the current directory, or a directory
# specified by the -o flag, it places the executable into the
# $GOPATH/bin directory.
#
# To find where your $GOPATH directory is located, run the following command:
#     go env GOPATH
# that is /google/gopath in cloud console. In the image it appears to be
#     /go
RUN echo $(go env) # debug
RUN echo $(pwd)    # debug
WORKDIR /go/src/helloweb
RUN go test -v ./...

# GIT_COMMIT is passed in via --build-arg (see build.sh) and baked into the
# binary as main.buildETag, so RootHandler serves a stable ETag tied to the
# source revision instead of falling back to epoch seconds (see issue #55).
ARG GIT_COMMIT=""
RUN go install -ldflags "-X main.buildETag=${GIT_COMMIT}" helloweb/webserver

# Now switch the image to the one we want to run in. The last image
# specified is the one that the built container uses.
# Using alpine because it is minimal at only 5MB in size! A full list
# of official images is found here:
# https://hub.docker.com/search?q=&type=image&categories=Operating%20Systems&operating_system=linux&image_filter=official
# Note that alpine is not officially supported for golang :fingers-crossed:
FROM alpine:latest AS target-image

# We copy our binary /go/bin/webserver to "/webserver".
# index.html is embedded into the binary (see //go:embed in webserver.go),
# so the runtime image no longer needs a copy of it on disk.
# The '--from=build-stage' means take the files from the
# first FROM statement image that we labeled as build-stage
# above.
RUN mkdir webserver
COPY --from=build-stage /go/bin/webserver webserver/.
RUN echo $(ls /webserver/.) # debug
ENV PORT=8080

# CMD defines the command that the container executes by default.
# We run the webserver from its location "/webserver".
WORKDIR /webserver
RUN echo $(ls .) # debug
RUN echo $(pwd)  # debug
CMD ["./webserver"]
# Comment out CMD to go interactive.

# ADD versus COPY commands. Seems they do the same thing for
# local files. But ADD can also copy files from a tar ball
# or URL.

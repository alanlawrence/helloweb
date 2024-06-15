FROM golang:1.22.4-alpine AS build-stage
# The FROM setting sets the base image used for the following commands.
# We are using an alpine base image with go version indicated installed.
# Alpine is small and fast. See the tags for available images here:
#     https://hub.docker.com/_/golang
# -----------------
# The following lines copy the contents of local directories into
# /go/src/* directories in the image. We assume we are in the webserver
# directory.
ADD digits    /go/src/helloweb/digits
ADD longDiv   /go/src/helloweb/longDiv
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
RUN go install helloweb/webserver

# Now switch the image to the one we want to run in. The last image
# specified is the one that the built container uses.
# Using alpine because it is minimal at only 5MB in size! A full list
# of official images is found here:
# https://hub.docker.com/search?q=&type=image&categories=Operating%20Systems&operating_system=linux&image_filter=official
# Note that alpine is not officially supported for golang :fingers-crossed:
FROM alpine:latest AS target-image

# We copy our binary /go/bin/webserver to "/webserver".
# And we copy our index.html file to the same directory
# since webserver expects it to be in its working dir.
# The '--from=build-stage' means take the files from the
# first FROM statement image that we labeled as build-stage
# above.
RUN mkdir webserver
COPY --from=build-stage /go/bin/webserver webserver/.
COPY --from=build-stage /go/src/helloweb/webserver/index.html webserver/.
RUN echo $(ls /webserver/.) # debug
ENV PORT 8080

# CMD defines the command that the container executes by default.
# We run the webserver from its location "/webserver" which we set as
# the working dir so it can read the index.html file there.
WORKDIR /webserver
RUN echo $(ls .) # debug
RUN echo $(pwd)  # debug
CMD ["./webserver"]
# Comment out CMD to go interactive.

# ADD versus COPY commands. Seems they do the same thing for
# local files. But ADD can also copy files from a tar ball
# or URL.

FROM golang:1.8-alpine
# The following lines copy the contents of local directories into
# /go/src/* directories in the image. We assume we are in the webserver
# directory.
ADD digits    /go/src/helloweb/digits
ADD longDiv   /go/src/helloweb/longDiv
ADD quadratic /go/src/helloweb/quadratic
ADD webserver /go/src/helloweb/webserver

# The go install command behaves almost identically to go build, but 
# instead of leaving the executable in the current directory, or a directory 
# specified by the -o flag, it places the executable into the 
# $GOPATH/bin directory.
#
# To find where your $GOPATH directory is located, run the following command:
#     go env GOPATH
# that is /google/gopath in cloud console. In the image it appears to be
#     /go
# We install the packages first, in hierachical order. Then the binary which
# depends on them.
RUN go install helloweb/digits
RUN go install helloweb/longDiv
RUN go install helloweb/quadratic
RUN go install helloweb/webserver 

# Don't know what this line does.
FROM alpine:latest

# We copy our binary /go/bin/webserver
# to the current directory in the image which is the root dir "/".
# And we copy our index.html file to the same directory
# since webserver expects it to be in its working dir.
# The '--from=0' means build-stage 0. Not sure what it means.
COPY --from=0 /go/bin/webserver .
COPY --from=0 /go/src/helloweb/webserver/index.html .
ENV PORT 8080

# CMD defines the process that the container executes by default.
# We run the webserver from its location "/" which also becomes
# the working dir.
CMD ["./webserver"]

# ADD versus COPY commands. Seems they do the same thing for
# local files. But ADD can also copy files from a tar ball 
# or URL.

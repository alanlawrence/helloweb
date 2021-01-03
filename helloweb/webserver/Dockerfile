FROM golang:1.8-alpine
# This line copies the contents of the current directory into a 
# /go/src/hello-app directory in the image. 
ADD . /go/src/hello-app

# The go install command behaves almost identically to go build, but 
# instead of leaving the executable in the current directory, or a directory 
# specified by the -o flag, it places the executable into the 
# $GOPATH/bin directory.
#
# To find where your $GOPATH directory is located, run the following command:
#     go env GOPATH
# that is /google/gopath in cloud console.
RUN go install hello-app

# Don't know what this line does.
FROM alpine:latest
# This line seems to be copying the contents of /go/bin/hello-app
# to the current directory (but where? A: to the root dir"/" in the image)
# The '--from=0' means build-stage 0. Not sure what it means.
COPY --from=0 /go/bin/hello-app .
COPY --from=0 /go/src/hello-app/index.html .
ENV PORT 8080

# CMD defines the process that the container executes by default.
CMD ["./hello-app"]

# ADD versus COPY commands. Seems they do the same thing for
# local files. But ADD can also copy files from a tar ball 
# or URL.

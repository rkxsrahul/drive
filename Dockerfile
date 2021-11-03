FROM golang:1.10


# Set go bin which doesn't appear to be set already.
ENV GOBIN /go/bin

# build directories
ADD . /go/src/git.xenonstack.com/util/drive-portal
WORKDIR /go/src/git.xenonstack.com/util/drive-portal

#Go dep!
# RUN go get -u github.com/golang/dep/...
# RUN dep ensure -update

RUN go install git.xenonstack.com/util/drive-portal
ENTRYPOINT /go/bin/drive-portal

EXPOSE 9301

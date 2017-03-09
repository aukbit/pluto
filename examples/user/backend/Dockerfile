FROM buildpack-deps:jessie-scm

# gcc for cgo
RUN apt-get update && apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
	&& rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.7.1
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 43ad621c9b014cde8db17393dc108378d37bc853aa351a6c74bf6432c1bbd182


RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

# My Golang project
# Download project dependencies
RUN go get github.com/aukbit/pluto

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/aukbit/pluto/examples/user/backend

# Build the pluto-sample-frontend command inside the container.
RUN go install github.com/aukbit/pluto/examples/user/backend

# Run the pluto-sample-frontend command by default when the container starts.
ENTRYPOINT /go/bin/backend -db_addr=192.168.99.100 -grpc_port=:65065

# Document that the service listens on port 65060.
EXPOSE 65065

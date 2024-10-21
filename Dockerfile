# Support setting various labels on the final image
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""

# Build Geth in a stock Go builder container
FROM golang:1.22

RUN apt update && apt-get install -y gcc musl-dev git nano

# Get dependencies - will also be cached if we won't change go.mod/go.sum
COPY go.mod /go-sintrop/
COPY go.sum /go-sintrop/
RUN cd /go-sintrop && go mod download

ADD . /go-sintrop
RUN cd /go-sintrop && go run build/ci.go install -static ./cmd/geth

WORKDIR /go-sintrop 

RUN make all
RUN cp build/bin/geth /usr/local/bin

RUN apt install ca-certificates

EXPOSE 8545 8546 30303 30303/udp
ENTRYPOINT ["/bin/bash"]

# Add some metadata labels to help programatic image consumption
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""

LABEL commit="$COMMIT" version="$VERSION" buildnum="$BUILDNUM"

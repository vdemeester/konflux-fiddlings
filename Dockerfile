# ARG GO_BUILDER=brew.registry.redhat.io/rh-osbs/openshift-golang-builder:v1.22
ARG GO_BUILDER=docker.io/library/golang:1.24
ARG RUNTIME=registry.access.redhat.com/ubi9/ubi-minimal:latest@sha256:daa61d6103e98bccf40d7a69a0d4f8786ec390e2204fd94f7cc49053e9949360

FROM $GO_BUILDER AS builder
COPY . .
RUN go build -o /bin/foo .

FROM $RUNTIME
COPY --from=builder /bin/foo /bin/foo
ENTRYPOINT ["/bin/foo"] 

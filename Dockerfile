# ARG GO_BUILDER=brew.registry.redhat.io/rh-osbs/openshift-golang-builder:v1.22
ARG GO_BUILDER=docker.io/library/golang:1.23
ARG RUNTIME=registry.access.redhat.com/ubi9/ubi-minimal:latest@sha256:d85040b6e3ed3628a89683f51a38c709185efc3fb552db2ad1b9180f2a6c38be

FROM $GO_BUILDER AS builder
COPY . .
RUN go build -o /bin/foo .

FROM $RUNTIME
COPY --from=builder /bin/foo /bin/foo
ENTRYPOINT ["/bin/foo"] 

# ARG GO_BUILDER=brew.registry.redhat.io/rh-osbs/openshift-golang-builder:v1.22
ARG GO_BUILDER=docker.io/library/golang:1.23
ARG RUNTIME=registry.access.redhat.com/ubi9/ubi-minimal:latest@sha256:c0e70387664f30cd9cf2795b547e4a9a51002c44a4a86aa9335ab030134bf392

FROM $GO_BUILDER AS builder
COPY . .
RUN go build -o /bin/foo .

FROM $RUNTIME
COPY --from=builder /bin/foo /bin/foo
ENTRYPOINT ["/bin/foo"] 

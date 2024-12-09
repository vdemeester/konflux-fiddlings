# ARG GO_BUILDER=brew.registry.redhat.io/rh-osbs/openshift-golang-builder:v1.22
ARG GO_BUILDER=docker.io/library/golang:1.23
ARG RUNTIME=registry.access.redhat.com/ubi9/ubi-minimal:latest@sha256:dee813b83663d420eb108983a1c94c614ff5d3fcb5159a7bd0324f0edbe7fca1

FROM $GO_BUILDER AS builder
COPY . .
RUN go build -o /bin/foo .

FROM $RUNTIME
COPY --from=builder /bin/foo /bin/foo
ENTRYPOINT ["/bin/foo"] 

# ARG GO_BUILDER=brew.registry.redhat.io/rh-osbs/openshift-golang-builder:v1.22
ARG GO_BUILDER=docker.io/library/golang:1.23
ARG RUNTIME=registry.access.redhat.com/ubi9/ubi-minimal:latest@sha256:14f14e03d68f7fd5f2b18a13478b6b127c341b346c86b6e0b886ed2b7573b8e0

FROM $GO_BUILDER AS builder
COPY . .
RUN go build -o /bin/foo .

FROM $RUNTIME
COPY --from=builder /bin/foo /bin/foo
ENTRYPOINT ["/bin/foo"] 

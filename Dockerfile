# Build image (Golang)
FROM golang:1.14-alpine AS build-stage
ENV GO111MODULE on
ENV CGO_ENABLED 0

RUN apk add --no-cache gcc git make

WORKDIR /src
ADD . .

RUN go mod download
RUN go build -o falcosidekick

# Final Docker image
FROM alpine AS final-stage
LABEL MAINTAINER "Thomas Labarussias <issif+falcosidekick@gadz.org>"

RUN apk add --no-cache ca-certificates

# Create user falcosidekick
RUN addgroup -S falcosidekick && adduser -u 1234 -S falcosidekick -G falcosidekick
# must be numeric to work with Pod Security Policies:
# https://kubernetes.io/docs/concepts/policy/pod-security-policy/#users-and-groups
USER 1234

WORKDIR ${HOME}/app
COPY --from=build-stage /src/falcosidekick .

EXPOSE 2801

ENTRYPOINT ["./falcosidekick"]

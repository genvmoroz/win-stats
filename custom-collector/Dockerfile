# image for compiling binary
ARG BUILDER_IMAGE="golang:1.23-alpine"
# here we'll run binary app
ARG RUNNER_IMAGE="alpine:3.19.1"

FROM ${BUILDER_IMAGE} as builder
### variables
# disable golang package proxying for this modules. necessary for private repos. Example: github.com/company-name
ARG GOPRIVATE=""
# ssh key for accessing private repos
ARG SSH_KEY

ENV GOPRIVATE ${GOPRIVATE}

RUN apk add --update gcc g++ openssh git make

# configure git to work with private repos
RUN mkdir -p ~/.ssh/ &&\
    echo "${SSH_KEY}" >> ~/.ssh/id_rsa && chmod 600 ~/.ssh/id_rsa && ssh-keyscan github.com >> ~/.ssh/known_hosts &&\
    git config --global url."git@github.com:".insteadOf "https://github.com/"

### copying project files
WORKDIR /app
# COPY the source code as the last step
COPY . .

# perform pre-build actions
RUN make deps
RUN make install-tools
RUN make generate
RUN make vulnerabilities_lookup
RUN make lint
RUN make test

# creates build/main files
RUN make build

FROM ${RUNNER_IMAGE}

RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories &&\
    apk update &&\
    apk add --no-cache\
    ca-certificates

COPY --from=builder /build/bin ./bin

CMD ["./bin", "s"]

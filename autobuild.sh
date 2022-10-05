cat <<EOF > Dockerfile2
FROM golang:alpine
# Install minimum necessary dependencies,
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3
RUN apk add --no-cache \$PACKAGES
# Set working directory for the build
WORKDIR /go/src/github.com/cosmos/cosmos-sdk
RUN git config --global --add safe.directory /go/src/github.com/cosmos/cosmos-sdk
EOF
docker build -t cosmos-build -f Dockerfile2 .
rm Dockerfile2

docker run -ti --rm -v $GOPATH/build_mod:/go/pkg -v $GOPATH/build_cache:/root/.cache -v $PWD:/go/src/github.com/cosmos/cosmos-sdk --name cosmos-build cosmos-build make build-linux

cp build/simd ./simd
cat <<EOF >Dockerfile3
FROM alpine:edge
# Install ca-certificates
RUN apk add --update ca-certificates
RUN apk add --no-cache jq
WORKDIR /root
# Copy over binaries from the build-env
COPY simd /usr/bin/simd
COPY scripts/init-chain.sh  ./init-chain.sh
COPY x/privacy/scripts/test.sh  ./test.sh
EXPOSE 26656 26657 1317 9090
# Run simd by default, omit entrypoint to ease using container with simcli
CMD ["simd"]
EOF
docker build -t simapp -f Dockerfile3 .
rm simd Dockerfile3

docker run -ti --rm --network host -v /tmp/testchain:/root/.simapp --name simapp simapp sh init-chain.sh
#docker run -ti --rm --network host -v /tmp/testchain:/root/.simapp --name simapp simapp simd start
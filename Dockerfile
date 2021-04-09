FROM arm64v8/golang
COPY . $GOPATH/jaeger_test
WORKDIR $GOPATH/jaeger_test
# Configure to reduce warnings and limitations as instruction from official VSCode Remote-Containers.
# See https://code.visualstudio.com/docs/remote/containers-advanced#_reducing-dockerfile-build-warnings.
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
    && apt-get -y install --no-install-recommends apt-utils 2>&1

# Verify git, process tools, lsb-release (common in install instructions for CLIs) installed.
RUN apt-get -y install git curl iproute2 procps lsb-release
RUN go get -u -v github.com/nsf/gocode
RUN go get -u -v github.com/rogpeppe/godef
RUN go get -u -v github.com/zmb3/gogetdoc
RUN go get -u -v github.com/lukehoban/go-outline
RUN go get -u -v sourcegraph.com/sqs/goreturns
RUN go get -u -v golang.org/x/tools/cmd/gorename
RUN go get -u -v github.com/tpng/gopkgs
RUN go get -u -v github.com/newhook/go-symbols
RUN go get -u -v golang.org/x/tools/cmd/guru
RUN go get -u -v github.com/cweill/gotests/...
RUN go get -u -v golang.org/x/tools/cmd/godoc
RUN go get -u -v github.com/fatih/gomodifytags
RUN go get golang.org/x/tools/gopls@latest
RUN go get github.com/go-delve/delve/cmd/dlv
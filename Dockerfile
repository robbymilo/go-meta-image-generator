FROM golang:1.17.2-buster AS build-env
WORKDIR /
COPY . .
RUN make build-linux

FROM debian:buster
RUN apt-get update && apt-get install wget ca-certificates chromium -y
RUN wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
RUN apt install ./google-chrome-stable_current_amd64.deb -y

COPY views /views

COPY --from=build-env /bin/meta-generator_linux-amd64 /usr/bin/meta-generator_linux-amd64
ENTRYPOINT ["/usr/bin/meta-generator_linux-amd64", "-views=/views"]
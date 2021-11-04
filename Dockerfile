FROM golang:1.17.2-bullseye AS build-env
WORKDIR /
COPY . .
RUN make build-linux

FROM debian:bullseye
RUN apt-get update && apt-get install unzip wget -y
RUN wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb \
  && apt install ./google-chrome-stable_current_amd64.deb -y

RUN wget https://storage.googleapis.com/chromium-browser-snapshots/Linux_x64/901912/chrome-linux.zip \
  && mkdir -p /root/.cache/rod/browser/chromium-901912 \
  && unzip chrome-linux.zip -d /root/.cache/rod/browser/chromium-901912

COPY views /views

COPY --from=build-env /bin/meta-generator_linux-amd64 /usr/bin/meta-generator_linux-amd64
ENTRYPOINT ["/usr/bin/meta-generator_linux-amd64", "-views=/views"]
FROM golang:1.17.2-bullseye AS build-env
WORKDIR /
COPY . .
RUN make build-linux

FROM debian:bullseye-slim
RUN apt-get update && apt-get install unzip wget chromium -y

RUN wget https://storage.googleapis.com/chromium-browser-snapshots/Linux_x64/901912/chrome-linux.zip \
  && mkdir -p /root/.cache/rod/browser/chromium-901912 \
  && unzip chrome-linux.zip -d /root/.cache/rod/browser/chromium-901912 \
  && rm chrome-linux.zip

COPY views /views
COPY public /public

COPY --from=build-env /bin/meta-generator_linux-amd64 /usr/bin/meta-generator_linux-amd64
ENTRYPOINT ["/usr/bin/meta-generator_linux-amd64", "-views=/views", "-public=/public"]
FROM golang:1.17.2-bullseye AS build-env
WORKDIR /
COPY . .
RUN make build-linux

FROM debian:bullseye-slim
RUN apt-get update && apt-get install unzip wget chromium -y

RUN useradd --create-home --shell /bin/bash pptuser

USER pptuser
WORKDIR /home/pptuser

RUN wget https://storage.googleapis.com/chromium-browser-snapshots/Linux_x64/901912/chrome-linux.zip \
  && mkdir -p /home/pptuser/.cache/rod/browser/chromium-901912 \
  && unzip chrome-linux.zip -d /home/pptuser/.cache/rod/browser/chromium-901912 \
  && rm chrome-linux.zip

ENV HOME=/home/pptuser

COPY views /home/pptuser/views
COPY public /home/pptuser/public

EXPOSE 3000

COPY --from=build-env /bin/meta-generator_linux-amd64 /usr/bin/meta-generator_linux-amd64
ENTRYPOINT ["/usr/bin/meta-generator_linux-amd64", "-views=/home/pptuser/views", "-public=/home/pptuser/public"]
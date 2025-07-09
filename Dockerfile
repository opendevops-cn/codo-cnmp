FROM golang:1.23.3 AS builder

WORKDIR /src
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn,direct"


COPY ./go.mod /src
COPY ./go.sum /src
RUN go mod download


COPY . /src
# RUN if [ -f ./config.yaml ]; then cp ./config.yaml ./etc; else echo "config.yaml not found, using defaults"; fi
RUN if [ -f ./config.yaml ]; then cp ./config.yaml ./etc/;elif [ -f ./etc/config.yaml ]; then echo "/etc/config.yaml already exists, skipping copy."; \
    else cp ./etc/config.example.yaml ./etc/config.yaml;fi



RUN make build

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

LABEL container_name="codo-cnmp"
ENV LANG=C.UTF-8
WORKDIR /data
COPY --from=builder /src/bin/codo-cnmp .
COPY --from=builder /src/etc/config.yaml ./etc/config.yaml
COPY --from=builder /src/migrate/yaml/ ./migrate/yaml/
COPY --from=builder /src/cert/server.*/ ./cert/
EXPOSE 8000
EXPOSE 9099
EXPOSE 8443
VOLUME /data/conf

CMD ["./codo-cnmp", "-f", "./etc/config.yaml"]

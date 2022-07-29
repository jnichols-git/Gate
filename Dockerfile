# syntax=docker/dockerfile:1-labs
FROM ubuntu:latest
# Requirement installations
RUN apt-get update
RUN apt-get install -y ca-certificates make
RUN apt-get clean
RUN apt-get install -y golang-go
RUN go version
# Copy in files
COPY . /gate-src/
WORKDIR /gate-src/
# Get certificate from "github.com" and "proxy.golang.org"
RUN openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/github.crt
RUN openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM >  ${cert_location}/proxy.golang.crt
RUN update-ca-certificates
# Get go.mod dependencies
RUN go mod download
# RUN go run ./cmd/certgen/certgen.go US Colorado Boulder DefaultOrg DefaultOrgUnit jakenichols.dev
CMD ["make", "server-run"]
EXPOSE 2719

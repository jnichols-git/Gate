# Gate

Go-based authentication server using SMTP or credentials for initial authentication, and JWTs for session verification,
hosted for you and by you. This server is lightweight and highly self-contained, having only 1 external library dependency
in [gorm](https://github.com/go-gorm/gorm), a widely-used and maintained library for database management. This is extended
to 2 external libraries if running a baremetal installation, which uses MetalLB to allow external connections.

Gate is currently in alpha. The author ([jnichols2719](https://github.com/jakenichols2719)) started this project as a 
means to practice Go and learn about security management methods, and as with any project created for practice and education, 
it has a ways to go before a stable and secure release. However, the author is dedicated to bringing this project up to par as a 
lightweight, stable, and secure authentication solution. If this sounds like a solution you could use, please feel free to write issues
or [contribute](CONTRIBUTING.md) to the project.

## Dependencies

Gate is built to run using the following tools:

- Linux/WSL
- Go version 1.13.8 or higher
- gcc: required to build gorm

## Installation/Usage

### FOR USERS

You'll need a few components to install Gate on your domain.

- Working installation of Docker/Kubernetes.
- Machine that can be reached at port 443 through a public IP address.
- Registered domain whose SSL certificate is either wildcard or can have a subdomain added.
    - Support incoming for non-subdomain operation.
    - You should direct `gate.domain` to the machine that will be running on.
- A working SMTP service, if you want to take advantage of email authentication.

Follow these steps to get started!

#### Setup

1. Set secrets for Gate
  - kubectl create secret generic gate-admin --from-literal=ADMIN_EMAIL=(email) --from-literal=ADMIN_USERNAME=(username) --from-literal=ADMIN_PASSWORD=(password)
  - kubectl create secret generic gate-smtp --from-literal=SMTP_USERNAME=(username) --from-literal=SMTP_PASSOWRD=(password)
  - kubectl create secret tls (domain)-tls --cert=(path to cert) --key=(path to key)
2. Create Kubernetes config
  - Download kubernetes/baremetal.yaml
  - Under Ingress, update the `localhost` entry to match your domain (including subdomain, so `jakenichols.dev` becomes `gate.jakenichols.dev`)
3. Run
  - kubectl apply baremetal.yaml
  - kubectl get pods
  - kubectl logs (gate pod name)
    - Your API key will be output in the logs.

You are now good to go! You can access the dashboard at `https://gate.domain/dashboard` and the API by sending requests to `https://gate.domain`.


### FOR DEVELOPERS

#### Setup

1. Gate installations have been tested on an Ubuntu machine through WSL. It should function, following these steps, in any Linux environment.
2. Install Go through `sudo apt-get install golang-go`
3. Install gccgo through `sudo apt-get install gccgo-go`
4. Clone your fork of the GitHub repository to a folder of your choice
5. Navigate to the base directory
6. Generate localhost certs through certgen by calling `go run ./cmd/certgen/certgen.go (country) (state) (locality) (organization) (organizational unit) localhost`
7. Set `Local` in `dat/config/config.yml` to `true` to override docker-based settings and listen on `localhost`.
8. Set environment variables.
    - `GATE_SSL_[KEY/CRT]`: SSL key and certificate file paths.
    - `GATE_SMTP_[USERNAME/PASSWORD]`: SMTP credentials.
    - `GATE_ADMIN_[EMAIL/USERNAME/PASSWORD]`: Admin credentials.

#### Run

Use `make server-run` to start a server. You can connect to it using `https://localhost:2719`.

## Configuration

Gate can be confgured through `gate-config.yml` or through the dashboard at `https://gate.domain/dashboard` (or `https://localhost/dashboard` if Local), once it's running.

## Project Goals

Gate has a set of goals and non-goals to reach for the project approaching a stable v1 release. As the project evolves past alpha, these
goals and non-goals may change in response to community feedback. These goals are, unless otherwise stated, not in any particular order.

### Goals

1. Gate should be secure above all. User data should be kept secure, and vital data such as passwords should be completely impractical
to attain even in the event of a security breach.
2. Setting up and running a Gate server should be simple and intuitive.
3. Servers should be fully functional on low-spec, low-cost servers for applications with lightweight needs. The majority of performance
overhead should be in the quantity and frequency of validation needed, not in the server's basic functionality.
4. Gate servers should be stable, handle errors and improper input well, and provide strong communication back to the application regarding
those errors.

### Non-Goals

1. Gate is not meant to replace, or be overall better than, industry standard authentication solutions.
2. The immediate focus of this project is not to reach a wide commercial audience.

## Getting Involved

Contributions to Gate are highly appreciated! Check [CONTRIBUTING](CONTRIBUTING.md) to learn what you can do to help.

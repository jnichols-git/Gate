# Gate

Go-based authentication server using SMTP or credentials for initial authentication, and JWTs for session verification,
hosted for you and by you. This server is lightweight and highly self-contained, having only 1 external library dependency
in [gorm](https://github.com/go-gorm/gorm), a widely-used and maintained library for database management.

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

- Working installation of Docker and docker-compose. Gate operates as a docker swarm service.
- Machine with a public IP address.
- Registered domain whose SSL certificate is either wildcard or can have a subdomain added.
    - Support incoming for non-subdomain operation.
    - You should direct `gate.domain` to the machine that will be running on.
- A working SMTP service, if you want to take advantage of email authentication.

Follow these steps to get started!

1. Configure Gate in your YAML file used for swarm deployment. This project uses `docker-compose.yml`. Required settings:
    - ports: Gate listens on port 2719. If your machine can be directly accessed at port 443, you should map 443:2719. If you're behind a network device,
    forward port 443 to port 2719 on your device, and map 2719:2719 in your file.
    - volumes: Should have a default volume mapping to directory /gate-src/dat/database.
    - secrets: The following secrets should be configured as external and set before deploying the gate service.
        - gate-ssl-key, gate-ssl-crt: Path to the key and cert file for your domain on your local machine.
        - gate-smtp-username, gate-smtp-password: Credentials for your SMTP server
        - gate-admin-email, gate-admin-username, gate-admin-password: Credentials for Gate. A server MUST be run using admin credentials. If none exist, an account will be created using email/username, printing a randomly-generated password and API key out to the terminal before exiting.
2. Initialize swarm using `docker swarm init`.
3. Initialize secrets for Gate.
4. Use `docker stack deploy -c [compose file] [stack name]` to start your application.

Your dashboard should now be accessible at `gate.domain/dashboard`, and you can make api calls through `gate.domain`.


### FOR DEVELOPERS

1. Gate installations have been tested on an Ubuntu machine through WSL. It should function, following these steps, in any Linux environment.
2. Install Go through `sudo apt-get install golang-go`
3. Install gccgo through `sudo apt-get install gccgo-go`
4. Clone your fork of the GitHub repository to a folder of your choice
5. Navigate to the base directory
6. Generate localhost certs through certgen by calling `go run ./cmd/certgen/certgen.go (country) (state) (locality) (organization) (organizational unit) localhost`

You should be set up to contribute from there. You can run `make server-run` to start a server and begin configuring your `localhost` testing environment.

## Configuration

Once basic settings have been established, Gate can be configured from the dashboard at `https://gate.domain/dashboard`.

## Project Goals

Gate has a set of goals and non-goals to reach for the project approaching a stable v1 release. As the project evolves past alpha, these
goals and non-goals may change in response to community feedback. These goals are, unless otherwise stated, not in any particular order.

### Goals

1. Gate should be secure above all. User data should be kept secure, and vital data such as passwords should be completely impractical
to attain even in the event of a security breach.
2. Setting up and running a Gate server should be simple and intuitive.
3. Servers should be fully functional on low-spec, low-cost servers for applications with lightweight needs. The majority of performance
overhead should be in the quanitity and frequency of validation needed, not in the server's basic functionality.
4. Gate servers should be stable, handle errors and improper input well, and provide strong communication back to the application regarding
those errors.

### Non-Goals

1. Gate is not meant to replace, or be overall better than, industry standard authentication solutions.
2. The immediate focus of this project is not to reach a wide commercial audience.

## Getting Involved

Contributions to Gate are highly appreciated! Check [CONTRIBUTING](CONTRIBUTING.md) to learn what you can do to help.

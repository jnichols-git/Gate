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

#### Setup

1. Create `docker-compose.yml`:
```
services:
  gate:
    image: jakenichols2719/gate
    ports:
      - 2719:2719
    volumes:
      - gate-db:/gate-src/dat/database
    secrets:
      - gate-ssl-key
      - gate-ssl-crt
      - gate-smtp-username
      - gate-smtp-password
      - gate-admin-email
      - gate-admin-username
      - gate-admin-password

secrets:
  gate-ssl-key:
    external: true
  gate-ssl-crt:
    external: true
  gate-smtp-username:
    external: true
  gate-smtp-password:
    external: true
  gate-admin-email:
    external: true
  gate-admin-username:
    external: true
  gate-admin-password:
    external: true

volumes:
  gate-db:
```
2. Create `gate-config.yml`, and configure the settings based on your providers and preferences:
```
Domain: website.com # Website domain. Used to find SSL certificates.
Address: "0.0.0.0" # IP address to listen on. Don't change this if you don't know what you're doing.
Port: 2719 # Port to listen on.
Local: false # Local run. Setting to true overrides domain/address to localhost and uses environment variables instead of docker secrets.
SMTP:
  Host: some-smtp-provider # SES provider
  Port: 587 # SMTP port; see your provider's settings
  Sender: notifications@website.com # Emails sent from this address
  TestEmail: testemail@website.com # Test email to send to
GateKey:
  UserValidTime: 1440 # Valid time for tokens for regular user authentication, in minutes
  AdminValidTime: 30 # Valid time for tokens for admin dashboard, in minutes
```

#### Run
1. Initialize swarm using `docker swarm init`.
2. Initialize secrets for Gate.
    - `gate-ssl-[key/crt]`: SSL key and certificate file paths.
    - `gate-smtp-[username/password]`: SMTP credentials.
    - `gate-admin-[email/username/password]`: Admin credentials.
3. Use `docker stack deploy -c [compose file] [stack name]` to start your application.
    - If you haven't run Gate before, the provided admin credentials will be used to create an admin account, and the terminal will output your API key. Make sure
    you use a secure password, and that you save that key. It will not be output anywhere else.

Your dashboard should now be accessible at `gate.domain/dashboard`, and you can make api calls through `gate.domain`.


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

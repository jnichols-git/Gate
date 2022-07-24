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

## Installation

WIP

## Configuration

The configuration file for Gate is located in `gate/config/config.yml`. The repo comes with an example config file and descriptions
of each field.

## Usage

WIP

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

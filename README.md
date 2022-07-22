# auth

Go-based authentication server using SMTP or credentials for initial authentication, and JWTs for session verification,
hosted for you and by you. This server is lightweight and highly self-contained, having only 1 external library dependency
in [gorm](https://github.com/go-gorm/gorm), a widely-used and maintained library for database management.

auth is currently in alpha. The author ([jnichols2719](https://github.com/jakenichols2719)) started this project as a 
means to practice Go and learn about security management methods, and as with any project created for practice and education, 
it has a ways to go before a stable and secure release. However, the author is dedicated to bringing this project up to par as a 
lightweight, stable, and secure authentication solution. If this sounds like a solution you could use, please feel free to write issues
or contribute to the project.

## Dependencies

auth is built to run using the following tools:

- Linux/WSL
- Go version 1.13.8 or higher
- gcc: required to build gorm

## Installation

WIP

## Configuration

The configuration file for auth is located in `auth/config/config.yml`. The repo comes with an example config file and descriptions
of each field.

## Usage

WIP

## Known Issues

Currently, there is not a static list of shortcomings; however, auth has not been examined by a security professional, and should not
be considered secure and ready for production environments until it has. Specific issues will be added here as they appear.

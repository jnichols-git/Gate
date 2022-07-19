# Auth Dashboard

As part of the `authserver` package, the Dashboard provides a place to check the health of your server,
download logs, and set permissions for certain users (banned, admin, etc.).

## Features

### SMTP Configuration

You can set your SMTP host, and which email you'd like to send from, through the dashboard. These changes
are saved (see Persistent Settings), so you can be sure your server is getting authentication emails out
as intended. The TestEmail field in `dat/config/config.yaml` serves as an email to send test messages to;
when you submit new details, the Dashboard will automatically send out a blank email to that address to
make sure it can do so without error. The icon next to SMTP: indicates whether your configuration is
working.

### Database Configuration

WIP

### TLS Certificate

The dashboard, on each refresh, will attempt to connect to `auth.domain` using TLS; the icon next to
TLS Certificate: indicates if this attempt is successful. The dashboard will additionally describe the error
if one occurs.

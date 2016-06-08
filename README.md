# CodeWiz
=================

A programming game where users write the AI for a wizard and battle each-other in an online arena.

# Configuration
## Environment Variables

The configuration of CodeWiz is driven through the following environment variables:

- **CODEWIZ\_DATABASE\_DSN**:

 A DSN containing the connection information for the database.

- **CODEWIZ\_DATABASE\_DRIVER**: 

 The name of the driver to use for database interactions. Currently supporting "mysql" and "sqlite3",

- **CODEWIZ\_LOG\_LEVEL**: 

 The lowest level that should be displayed in the log output. The options are "debug", "info", "warn", "error", and "fatal".

- **CODEWIZ\_PORT**: 

 The port on which the CodeWiz server should listen for requests.

- **CODEWIZ\_SESSION\_KEY**: 

 The key to use for encrypting session information sent between the client and server.

- **CODEWIZ\_SESSION\_SECURE**: 

 A flag indicating whether to allow session information to be sent over connections that are not protected by TLS/SSL. For security purposes, this should be set to "true" where possible, but can be configured to "false" for development environments where this extra security is neither available or required.
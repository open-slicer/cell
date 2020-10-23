# Configuration

Cell looks for a config file in the working directory. It can be of JSON, TOML, YAML, HCL or Java properties format; the only requirement is that the file is called `cell.extension`, where extension can of course fit the format you're using. If an extension isn't provided, YAML will be used. This manual will use YAML.

## Examples

### cell.yml
```yaml
environment: release # Enables verbose logging (debug or release)
node: 0 # The ID of the node; only change this if you're scaling up
http:
    address: :8080 # The address to bind on
database:
    postgres: postgres://username:password@localhost/cell # A pg connection string
    redis:
        address: localhost:6379 # The Redis instance's address
        password: "" # The Redis password; leave empty if you haven't set one
        db: 0 # The Redis DB to use
sentry:
    dsn: https://00000000000000000000000000000000@0000000.ingest.sentry.io/0000000 # The optional Sentry DSN to use
cors:
    allowed_origins: ["*"]
    allowed_methods: ["*"]
    allowed_headers: ["*"]
    exposed_headers: ["*"]
    allow_credentials: true
security:
    secret: rh4NaXhju914cn60CHmuMREeQG1Qdh53o4sQ9iZWVlA= # A secure 32 byte key; try `openssl rand -base64 32`
    cert_file: "" # The optional location of an SSL cert to use
    key_file: "" # The optional location of an SSL key to use
locket:
    token: extremely secure password # The password used to secure the `/lockets` endpoint
prometheus:
    token: extremely secure password # The password used to secure the `/metrics` endpoint
```

### locketd.yml
```yaml
environment: release # Enables verbose logging (debug or release)
port: 8000 # Port for locketd
security:
    rh4NaXhju914cn60CHmuMREeQG1Qdh53o4sQ9iZWVlA= # A secure 32 byte key; try `openssl rand -base64 32`
registration: # To register your locketd instance with Cell
    home: http://localhost:8080 # Where to make a request to Cell at
    token: extremely secure password # The token used to register with Cell
    host: "" # Hostname of the locket (defaults to IP if not specified)
```

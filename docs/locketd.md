# locketd

`locketd` is a WebSocket server daemon for Cell that uses Cell's Locket API. This allows it to be horizontally scalable, as instances can register with the main Cell server by making an authorised HTTP request. Redis is used to dispatch requests to `locketd`; read more in the [scaling section](scaling.md).

## API

The Locket API's main endpoint is `/api/v2/lockets`.

### Getting a Locket

As a client (not a server like `locketd`), you may request to be assigned a Locket. This is required for WebSocket operations. To do this, you can make a GET request to `/api/v2/lockets`. This requires a bearer token.

#### Response

You'll receive an application/json response body. The `data` will be a string representing the address at which you can access the Locket: `domainOrIP:port`.

### Registering a Locket

Lockets can be registered by making a `PUT` request to `/api/v2/lockets`. This expects the config value for `locket.token` to be present in the `Authorization` header.

#### Request

It expects an application/json body:

```json
{
  "port": 8080,
  "host": "some.locket.com"
}
```

##### `port`

Required. This is the TCP port that the Locket server is listening on.

##### `host`

Defaults to the request IP. If provided, this should be a domain that points to the request IP.

If there was an error looking up the domain, `errorDomainFailedLookup` will be thrown. If the lookup succeeded but the resulting IPs didn't include the request IP, `errorDomainDidntMatch` will be thrown.

#### Response

You'll receive an application/json response. Its `data` will be a `locketInterface` object:

```json
{
  "port": 8080,
  "host": "some.optional.locket.com"
}
```

#### Errors

<!-- TODO: number representations -->

`errorDomainFailedLookup` and `errorDomainDidntMatch` are assigned to this endpoint. See the body section for when it'll be thrown.
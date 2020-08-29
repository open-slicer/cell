# locketd

`locketd` is a WebSocket server daemon for Cell that uses Cell's Locket API. This allows it to be horizontally scalable, as instances can register with the main Cell server by making an authorised HTTP request. Redis is used to dispatch requests to `locketd`; read more in the [scaling section](scaling.md).

## API

The Locket API's main endpoint is `/api/v2/locket`. It expects the config value for `locket.token` to be present in the `Authorization` header.

### Registering a Locket

Lockets can be registered by making a `PUT` request to `/api/v2/locket`.

#### Body

It expects an application/json body:

```json
{
  "port": 8080,
  "host": "some.locket.com"
}
```

##### `port`

Required. This is the port that the Locket server is running on.

##### `host`

Defaults to the request IP. If provided, this should be a domain that points to the request IP.

If there was an error looking up the domain, `errorDomainFailedLookup` will be thrown. If the lookup succeeded but the resulting IPs didn't include the request IP, `errorDomainDidntMatch` will be thrown.

#### Errors

<!-- TODO: number representations -->

`errorDomainFailedLookup` and `errorDomainDidntMatch` are assigned to this endpoint. See the body section for when it'll be thrown.
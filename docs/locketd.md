# locketd

`locketd` is a WebSocket server daemon for Cell. This allows it to be horizontally scalable, as instances can register with the main Cell server by making an authorised HTTP request. Redis is used to dispatch requests to `locketd`; read more in the [scaling section](scaling.md).

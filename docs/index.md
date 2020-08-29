# Introduction

Welcome to the Cell documentation. Currently, it's expected that Cell's API will change _a lot_; it's in early stages. This manual should act as developers' documentation for **API v2**, but also a reference for contributors to Cell.

API v1 is deprecated and isn't available in Cell.

## Main ideas

The REST API exposes everything outgoing that doesn't require a persistent connection. This allows clients without the ability to hold long-running connections to use Cell. Contrary to this, [WebSocket](https://en.wikipedia.org/wiki/WebSocket) is used for incoming data (e.g. events); this is what allows realtime communication.

For example, to send a message you'd make a REST POST request to `/api/v2/channels/:id/messages`. This is then able to be received by GETting `/api/ws` and waiting for an `EVT_MESSAGE_CREATE` event. The payload will also contain the message body, so making another request to get it shouldn't be a requirement.

## Encryption

PGP is used to encrypt messages. This must be done on the client; otherwise, it'd deprecate the entire reason this exists! The only thing Cell has to do about encryption is message validation and public key storage. On top of this, [bcrypt](https://en.wikipedia.org/wiki/Bcrypt) is used to hash passwords.

This will be covered in more detail.

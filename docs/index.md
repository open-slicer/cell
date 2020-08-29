# Introduction

Welcome to the Cell documentation. Currently, Cell's API is expected to change _a lot_; it's in early stages. This manual should act as developers' documentation for **API v2**, available at the `/api/v2` endpoint.

API v1 is deprecated and isn't available in Cell.

## Main ideas

Everything outgoing that doesn't require a persistent connection is exposed over a REST API. This allows clients without the ability to hold long-running connections to use Cell. Contrary to this, [WebSocket](https://en.wikipedia.org/wiki/WebSocket) is used for incoming data (e.g. events); this is what allows realtime communication.

For example, to send a message you'd make a REST POST request to `/api/v2/channels/:id/messages`. This could then be received by GETting `/api/ws` and waiting for an `EVT_MESSAGE_CREATE` event. The payload will also contain the message body, so making another request to get it shouldn't be a requirement.

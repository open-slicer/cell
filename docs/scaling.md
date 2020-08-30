# Scaling

This covers how Cell scales. Most of this is Cell-specific and not related to the Slicer specification.

## WebSocket

Firstly, users GET `/api/v2/lockets` (excluding the version). This endpoint should inform clients of a suitable WS server (Locket), where they'll be able to establish a connection. WS servers should check if the server permitted the action in the last minute before allowing them to connect.

!!! note
    This is specific to Cell; it's not part of the Slicer specification.

Cell runs separate processes for WebSocket connections. The main process dispatches data to be sent over WS with Redis's pub/sub system, keyed by the user ID. All processes handling connections for the user receive the request and send the message over to each client; you're of course able to have more than one concurrent connection to Cell.

!!! warning
    This is unimplemented.

On top of this, each client (connection) has a unique ID. This allows the server to send connection-specific packets.
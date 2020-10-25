<div align="center">
  <img src="./assets/Logo.svg" width="200">
  <p>Cell is the reference implementation of the Slicer server.</p>
  <img src="https://goreportcard.com/badge/github.com/open-slicer/cell">
  <img src="https://github.com/open-slicer/cell/workflows/Go%20Build/badge.svg">
  <a href="http://www.apache.org/licenses/LICENSE-2.0.html"><img alt="Apache 2.0" src="https://owo.gg/alex/sourcebadges/-/raw/master/SVG/CC_Apache%202.0.svg" height="20"></a>
</div>

## what

Cell implements the Slicer protocol in Go, consisting of two main components:

- `cell` itself; and
- `locketd`.

`cell` is the basis for all operations. It has no in-memory state, so each instance is autonomous from clients. `locketd` instances register with `cell`; they don't actually do any work by themselves. `cell` broadcasts messages to all `locketd` instances keyed by the recipients' user IDs. They then pick up these messages, sending them to clients as required.

## but why

boredom

oh but also, IRC is old and insecure. Modern platforms are fast. The latter usually can't be said about newer encrypted platforms.

## stuff to do i guess

- [x] Users
  - [x] Create
  - [x] Get
  - [x] Auth
    - [x] Login
    - [x] Refresh
- [ ] Channels
  - [x] Create
  - [x] Get
  - [ ] Announce (ws)
  - [ ] Invites
    - [x] Create
    - [ ] Get
      - [x] By name
      - [ ] All by channel
    - [x] Accept
    - [ ] Announce (ws)
- [x] Lockets
  - [x] Node
  - [x] Register
  - [x] Get/rotate
- [ ] Administration
  - [ ] Metrics
    - [x] Prometheus
    - [ ] Some other instance-specific stuff could be done.
  - [ ] Account management, etc.
- [ ] Kubernetes


## license

<a href="http://www.apache.org/licenses/LICENSE-2.0.html"><img alt="Apache 2.0" src="https://owo.gg/alex/sourcebadges/-/raw/master/SVG/CC_Apache%202.0.svg" height="100"></a>

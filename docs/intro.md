<div align="center">
  <img src="./img/cell.png" width="200">
  <p>Cell is the reference implementation of the Slicer server.</p>
  <img src="https://goreportcard.com/badge/github.com/open-slicer/cell">
  <img src="https://github.com/open-slicer/cell/workflows/Go%20Build/badge.svg">
</div>

## Roadmap

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

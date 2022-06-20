# schemacheck
A CLI utility written in [go](go.dev) that validates `json` and `yaml` files
against a `schema`.

## Install
There are a few different methods to install `schemacheck`.

### Via `go` (Recommended)
* Run `go install github.com/adrielp/schemacheck`

### Mac/Linux during local development
* Clone down this repository and run `make build`
* Install a binary for your platform from `dist/bin` locally to a path


### Windows
There's a binary for that, but it's not directly supported or tested because #windows

## Getting Started
### Prereqs
* Have [make](https://www.gnu.org/software/make/) installed
* Have [GoReleaser](https://goreleaser.com/) installed

### Instructions
* Clone down this repository
* Run commands in the [Makefile](./Makefile) like `make build`

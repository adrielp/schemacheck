# schemacheck
A CLI utility written in [go](go.dev) that validates `json` and `yaml` files
against a `schema`.

## Usage 
`schemacheck` is meant to be used against one schema and one or more `yaml` or
`json` files. 

After installation, you can run it like:
```
schemacheck --schema myschema.json --file myjson.json --file myyaml.yaml .......
```

You can get the usage at any time by running:
```
schemacheck --help
```

You can also call this CLI from other command line utililties like `find`.
```
find . -type f -name "*.json" -exec ./dist/bin/schemacheck -s test_data/schema.json -f {} \+
```

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

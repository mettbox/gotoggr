# Gotoggr

> Report your daily toggl tracked tasks

## Requirements

A Toggl API key is required. Generate one in your Toggl profile at
<https://track.toggl.com/profile> and store it with `gotoggr --set-token <your-token>`.

## Install

```sh
go install .
```

This builds the binary and puts it in `$(go env GOPATH)/bin`.

To run it from anywhere, add that directory to your `PATH` — in `~/.zshrc` (zsh, the macOS default):

```sh
export PATH=$PATH:$(go env GOPATH)/bin
```

## Usage

Set the Toggl API key

```sh
gotoggr --set-token <your-token>
```

Show the stored API key

```sh
gotoggr --show-token
```

Get latest time entries

```sh
gotoggr
```

## Development

Run without installing:

```sh
go run .
```

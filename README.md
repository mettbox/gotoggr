# Gotoggr

> Report your daily toggl tracked tasks

## Usage

Setup Toggl API Token

```sh
gotoggr --token <your-token>
```

Get latest time entries

```sh
gotoggr
```

## Build and Install

Build the Go app:

```sh
go build -o gotoggr
```

Make it available in your bash:

```sh
# Add GOPATH/bin directory to your PATH environment variable so you can run Go programs anywhere.
export GOPATH=$HOME/go
export PATH=$PATH:$(go env GOPATH)/bin
```

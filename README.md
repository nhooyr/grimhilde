# grimhilde

Mirror, mirror, on the wall, who’s the fairest URL of them all?

## Install

```
go get nhooyr.io/grimhilde
```

## Usage

First modify the [config_example.json](./cmd/grimhilde/config_example.json) to suit your needs.

Then rename it to `config.json` and run:

```
gcloud app deploy --version=1 cmd/grimhilde
```

Unfortunate to have the main command under `cmd` but its necessary as otherwise app engine goes crazy.

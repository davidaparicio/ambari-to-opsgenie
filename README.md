# Ambari to Opsgenie

## Upupa epops

This Go server forwards [Ambari](https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/index.md) alerts to [OpsGenie](https://www.opsgenie.com/).

## Installation

```bash

```

## Usage

```bash

```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Development

**For an optimal developer experience, it is recommended to install [Nix](https://nixos.org/download.html) and [direnv](https://direnv.net/docs/installation.html).**

_Alternatively, install [Go](https://go.dev/dl/) on your computer then run `make deps` to install the rest of the dependencies._

Run the test suite:

```shell
make test
```

Run linters:

```shell
make lint # pass -j option to run them in parallel
```

Some linter violations can automatically be fixed:

```shell
make fmt
```

## License
[//]: # (https://www.makeareadme.com/)

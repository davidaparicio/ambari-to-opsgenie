# Ambari to Opsgenie

## Upupa epops

This Go server forwards [Ambari](https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/index.md) alerts to [OpsGenie](https://www.opsgenie.com/). 

Some extensions have been added like [Blinky](https://getblinky.io/) or [xbar](https://github.com/matryer/xbar) integration (for MacOS).

## Installation

```bash
age-keygen -o secrets/age.key > secrets/public_age.key 2>&1
```

### Configuration

Edit with correct information and generate the new configuration file with
```bash
export SOPS_AGE_RECIPIENTS=$(<secrets/public_age.key)
sops --encrypt --age ${SOPS_AGE_RECIPIENTS} configs/config.yaml > configs/config.enc.yaml
```

## Usage

First, for in-flight decryption purposes, you need to specify your AGE private key, only, if it's not in ```secrets/age.key``` with this command

```bash
export SOPS_AGE_KEY_FILE="/TOCHANGETO/Your/Path/Of/Your/AGE/Private.key"
```


For Mac OS, [xbar](https://github.com/matryer/xbar)+[Blinky](https://getblinky.io/) integrations inside the script ```/scripts/xbar_alert.5m.sh```

or manually with
```bash
go run cmd/xbar/main.go
```

If needed, you can run a very simple mockserver of Ambari
```bash
go run cmd/mockserver/main.go
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[//]: # (https://www.makeareadme.com/)

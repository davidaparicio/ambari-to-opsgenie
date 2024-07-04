# Ambari to Opsgenie

<p align="center">
<img src="assets/img/a2o_logo.jpg" alt="Ambari-to-Opsgenie logo" title="Ambari-to-Opsgenie logo" />
</p>

[![Docs](https://img.shields.io/badge/docs-current-brightgreen.svg)](https://pkg.go.dev/github.com/davidaparicio/ambari-to-opsgenie)
[![Go Report Card](https://goreportcard.com/badge/davidaparicio/ambari-to-opsgenie)](https://goreportcard.com/report/davidaparicio/ambari-to-opsgenie)
[![Github](https://img.shields.io/static/v1?label=github&logo=github&color=E24329&message=main&style=flat-square)](https://github.com/davidaparicio/ambari-to-opsgenie)
[![GitLab](https://img.shields.io/static/v1?label=gitlab&logo=gitlab&color=green&message=mirrored&style=flat-square)](https://gitlab.com/davidaparicio/ambari-to-opsgenie)
[![Froggit](https://img.shields.io/static/v1?label=froggit&logo=froggit&color=red&message=no&style=flat-square)](https://lab.frogg.it/davidaparicio/ambari-to-opsgenie)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/davidaparicio/ambari-to-opsgenie/blob/main/LICENSE.md)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fdavidaparicio%2Fambari-to-opsgenie.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fdavidaparicio%2Fambari-to-opsgenie?ref=badge_shield)
[![Maintenance](https://img.shields.io/maintenance/yes/2024.svg)]()
[![Twitter](https://img.shields.io/twitter/follow/dadideo.svg?style=social)](https://twitter.com/intent/follow?screen_name=dadideo)

[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=davidaparicio_ambari-to-opsgenie&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=davidaparicio_ambari-to-opsgenie)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=davidaparicio_ambari-to-opsgenie&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=davidaparicio_ambari-to-opsgenie)

This Go server forwards [Ambari](https://github.com/apache/ambari/blob/trunk/ambari-server/docs/api/v1/index.md) alerts to [OpsGenie](https://www.opsgenie.com/). 

Some extensions have been added like [Blinky](https://getblinky.io/) or [xbar](https://github.com/matryer/xbar) integration (for MacOS).

## Installation

```bash
brew install sops age && mkdir secrets
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

## Contribute

Works on my machine - and yours ! Spin up pre-configured, standardized dev environments of this repository, by clicking on the button below.

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#/https://github.com/davidaparicio/gokvs)

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
Licensed under the MIT License, Version 2.0 (the "License"). You may not use this file except in compliance with the License.
You may obtain a copy of the License [here](https://choosealicense.com/licenses/mit/).

If needed some help,  there are a ["Licenses 101" by FOSSA](https://fossa.com/blog/open-source-licenses-101-mit-license/), a [Snyk explanation](https://snyk.io/learn/what-is-mit-license/)
of MIT license and a [French conference talk](https://www.youtube.com/watch?v=8WwTe0vLhgc) by [Jean-Michael Legait](https://twitter.com/jmlegait) about licenses.


[//]: # (https://www.makeareadme.com/)

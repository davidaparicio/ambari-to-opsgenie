# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2023-08-11
### Changed
- Refactoring and fit to the [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- Use with wisdom, the best data structure for the [config](https://david-yappeter.medium.com/golang-pass-by-value-vs-pass-by-reference-e48aac8b2716)[²](https://go101.org/article/pointer.html)

## [0.0.4] - 2023-07-30
### Changed
- Upgrade all dependencies to the latest version
- Read carefully the next migration [go.mozilla.org/sops/v3](https://github.com/getsops/sops/issues/1246)->[github.com/getsops/sops](https://github.com/getsops/sops)

## [0.0.3] - 2021-01-08
### Added
- Encrypt password with [go.mozilla.org/sops/v3](github.com/mozilla/sops)
- Enable GoReleaser
- Unit tests

## [0.0.2] - 2021-01-08
### Changed
- Upgrade all dependencies to the latest version

## [0.0.1] - 2020-11-04
### Changed
- Migrate Go Modules to v2+
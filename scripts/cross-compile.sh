#! /usr/bin/env bash

echo "================================================="
echo "Build locally without changing the version/tag..."
echo "currently for 64bits arch / MacOS(Darwin)+Linux.."
echo "if you need to change it, check .goreleaser.yaml "
echo "================================================="
goreleaser release --snapshot --skip-publish --rm-dist

# Some requirements / operations performed on daparici's MacOS laptop
# brew install goreleaser/tap/goreleaser

# brew install sops && sops --version
# brew install age && age --version
# age-keygen -o key.txt
# sops --encrypt --age age1c2gv7m4zr737lwmek8m8cfzygfdv6s5cpr96839y9s93ugrkv5dqvc6yv0 config.yaml > test.enc.yaml

# Test with GPG (and testing keys)
# gpg --import $GOPATH/pkg/mod/go.mozilla.org/sops/v3@v3.7.1/pgp/sops_functional_tests_key.asc 
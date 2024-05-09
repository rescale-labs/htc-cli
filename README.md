# htc-storage-cli

This is the storage cli written in go to be used for HTC Nextflow jobs.


## Push a new version

To deploy, update the `VERSION` variable in the makefile then run `make push-dev` to push a new version to dev. To push to prod run `make push-prod`.

## Archive a Linux Distribution Agnostic Binary

To archive the binary run `make archive`. This will create 2 tar archives `htccli.linux-amd64.tar.gz` and `htccli.linux-arm64.tar.gz` which contain the distribution agnostic binaries for the 2 architectures.

## Testing

To test run `make test`.

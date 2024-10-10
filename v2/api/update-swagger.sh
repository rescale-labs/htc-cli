#!/bin/sh -e

DIR=$(cd $(dirname $0); pwd)

cd $DIR
curl -o swagger.json -H "Accept: application/json" \
    https://htc.rescale.com/q/openapi
# Write sorted json and patched, sorted json so we can diff them.
jsonnet -e "import 'swagger.json'" > swagger-sorted.json
jsonnet swagger.jsonnet -o swagger-patched.json

echo "\nDifferences from original and swagger-patched.json:"
diff -u swagger-sorted.json swagger-patched.json || true # diff exits 1 for differences

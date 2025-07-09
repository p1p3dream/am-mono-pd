#!/bin/sh -e

name=$1

if [ -z "$2" ]; then
    output_dir=${ABODEMINE_WORKSPACE}/etc/ssl
else
    output_dir=$2
fi

yq eval -o=json ${ABODEMINE_WORKSPACE}/infra/cfssl/${name}.yaml | jq -Mc > ${ABODEMINE_WORKSPACE}/infra/cfssl/${name}.json
cfssl gencert -initca ${ABODEMINE_WORKSPACE}/infra/cfssl/${name}.json | cfssljson -bare ${output_dir}/${name}
rm ${ABODEMINE_WORKSPACE}/infra/cfssl/${name}.json

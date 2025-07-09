#!/bin/sh -e

ca=$1
profile=$2
config=$3

if [ -z "$4" ]; then
    output_dir=${ABODEMINE_WORKSPACE}/etc/ssl
else
    output_dir=$4
fi

yq eval -o=json --indent 0 ${ABODEMINE_WORKSPACE}/infra/cfssl/${ca}-config.yaml > ${ABODEMINE_WORKSPACE}/infra/cfssl/${ca}-config.json
yq eval -o=json --indent 0 ${ABODEMINE_WORKSPACE}/infra/cfssl/${config}.yaml > ${ABODEMINE_WORKSPACE}/infra/cfssl/${config}.json

cfssl gencert \
	-ca ${output_dir}/${ca}.pem \
	-ca-key ${output_dir}/${ca}-key.pem \
	-config ${ABODEMINE_WORKSPACE}/infra/cfssl/${ca}-config.json \
	-profile=${profile} \
	${ABODEMINE_WORKSPACE}/infra/cfssl/${config}.json \
	| cfssljson -bare ${output_dir}/${config}-${profile}

rm ${ABODEMINE_WORKSPACE}/infra/cfssl/${ca}-config.json ${ABODEMINE_WORKSPACE}/infra/cfssl/${config}.json

cat ${output_dir}/${config}-${profile}.pem ${output_dir}/${ca}.pem > ${output_dir}/${config}-${profile}-chain.pem

openssl pkcs12 \
	-export \
	-in ${output_dir}/${config}-${profile}-chain.pem \
	-inkey ${output_dir}/${config}-${profile}-key.pem \
	-out ${output_dir}/${config}-${profile}.pfx \
	-passout pass:

cfssl certinfo -cert ${output_dir}/${config}-${profile}.pem

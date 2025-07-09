# Make

```sh
ABODEMINE_BUILD_ID=$(uuidgen) \
ABODEMINE_NAMESPACE=local \
make -C ${ABODEMINE_WORKSPACE}/infra/opentofu/projects/dev/buckets build package configure release
```

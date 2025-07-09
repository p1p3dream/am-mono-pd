# Make

```sh
ABODEMINE_BUILD_ID=$(uuidgen) \
ABODEMINE_NAMESPACE=... \
make -C ${ABODEMINE_WORKSPACE}/infra/opentofu/projects/main build package configure release
```

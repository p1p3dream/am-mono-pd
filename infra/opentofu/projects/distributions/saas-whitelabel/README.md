# Make

```sh
ABODEMINE_BUILD_ID=$(uuidgen) \
ABODEMINE_NAMESPACE=... \
make -C ${ABODEMINE_WORKSPACE}/infra/opentofu/projects/distributions/saas-whitelabel build package configure release
```

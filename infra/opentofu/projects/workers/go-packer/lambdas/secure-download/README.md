# Make

```sh
ABODEMINE_BUILD_ID=$(uuidgen) \
ABODEMINE_NAMESPACE=... \
make -C ${ABODEMINE_WORKSPACE}/infra/opentofu/projects/workers/go-packer/lambdas/secure-download build package configure release
```

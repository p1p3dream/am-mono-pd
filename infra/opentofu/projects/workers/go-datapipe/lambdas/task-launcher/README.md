# Make

```sh
ABODEMINE_BUILD_ID=$(uuidgen) \
ABODEMINE_NAMESPACE=... \
make -C ${ABODEMINE_WORKSPACE}/infra/opentofu/projects/workers/go-datapipe/lambdas/task-launcher build package configure release
```

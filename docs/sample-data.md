# Datapipe

> These steps will clear everything from the supported datapipe tables.

Ensure you have latest migrations with:

```sh
make -C ${ABODEMINE_WORKSPACE}/infra/docker/projects/am-mono migrate-datapipe/up
```

Then:

```sh
make -C ${ABODEMINE_WORKSPACE}/assets/projects/datapipe/sample-data load
```

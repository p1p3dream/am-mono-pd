# on datapipe_testing

```sh
rm -rf /works/samples/api
mkdir -p /works/samples/api
```

```sql
\copy (
    select *
    from api_keys
) to '/works/samples/api/api_keys.txt';
```

```sql
\copy (
    select *
    from clients
) to '/works/samples/api/clients.txt';
```

```sh
zstd -T0 -18 -f /works/samples/api/*

aws s3 sync \
    /works/samples/api/ \
    s3://testing-mono-overlay-c5tfj3sl/assets/projects/api/sample-data/.tmp/ \
    --exclude "*.txt"
```

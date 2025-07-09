# on datapipe_testing

```sh
rm -rf /works/samples/saas
mkdir -p /works/samples/saas
```

```sql
\copy (
    select *
    from organizations
) to '/works/samples/saas/organizations.txt';
```

```sql
\copy (
    select *
    from roles
) to '/works/samples/saas/roles.txt';
```

```sql
\copy (
    select *
    from users
) to '/works/samples/saas/users.txt';
```

```sh
zstd -T0 -18 -f /works/samples/saas/*

aws s3 sync \
    /works/samples/saas/ \
    s3://testing-mono-overlay-c5tfj3sl/assets/projects/saas/sample-data/.tmp/ \
    --exclude "*.txt"
```

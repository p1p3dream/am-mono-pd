## RentalAVM

### Add csvkit to PATH

```sh
export PATH=${ABODEMINE_WORKSPACE}/code/python/abodemine/.venv/bin:${PATH}
```

### Create the csvstat file.

```sh
cat /works/data/dumps/rentalavm-22/RENTNECTAR_RENTALAVM_0022_001.txt \
| head -250000 \
| csvstat -t -u 3 --no-leading-zeroes --json \
> ${ABODEMINE_WORKSPACE}/assets/projects/datapipe/csvstat/attom-data/rental-avm.json
```

### Create the datafiles.

```sh
make -C ${ABODEMINE_WORKSPACE}/code/go/abodemine/tools/csvmine run \
RUN_ARGS="txt --stat ${ABODEMINE_WORKSPACE}/assets/projects/datapipe/csvstat/attom-data/rental-avm.json --table ad_df_rental_avm --struct RentalAvm --package attom_data"
```

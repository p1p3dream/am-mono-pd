> Remember to change `user@am-env-shell` with your local identity.

- First update `src/boot` with upstream:

```sh
rsync \
    -av \
    --exclude=\.git \
    --delete \
    am-env-shell:/works/src/mono/infra/docker/projects/am-env/ \
    ${HOME}/works/abodemine/src/boot/infra/docker/projects/am-env/

rsync \
    -av \
    am-env-shell:/works/src/mono/code/make/core.mk \
    ${HOME}/works/abodemine/src/boot/code/make/core.mk
```

- Run the setup again:

```sh
ABODEMINE_WORKSPACE=${HOME}/works/abodemine/src/boot \
gmake -C ${HOME}/works/abodemine/src/boot/infra/docker/projects/am-env up
```

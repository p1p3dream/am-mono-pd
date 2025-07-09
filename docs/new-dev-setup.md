# Make env

- Download https://github.com/abodeminehq/boot as zip.
- Execute in the host (macOS) terminal:

```sh
mkdir -p ${HOME}/works/abodemine/src

unzip -d ${HOME}/works/abodemine/src ${HOME}/boot-trunk.zip

mv ${HOME}/works/abodemine/src/boot-trunk ${HOME}/works/abodemine/src/boot

find ${HOME}/works/abodemine/src/boot -type d | xargs -I{} chmod 0777 {}
find ${HOME}/works/abodemine/src/boot -type f | xargs -I{} chmod 0666 {}

docker volume create am_src
docker volume create am_works
docker network create abodemine.local
```

- Inside `${HOME}/works/abodemine/src/boot/infra/docker/projects/am-env`, create:

  - `infra/docker/projects/am-env/.env`
  - `infra/docker/projects/am-env/config.local.yaml`
  - gpg keys
  - ssh keys

- Execute the setup:

```sh
ABODEMINE_WORKSPACE=${HOME}/works/abodemine/src/boot \
make -C ${HOME}/works/abodemine/src/boot/infra/docker/projects/am-env up
```

- Add the following to `~/.ssh/config` on the host:

> Notice that you have to change **username** with the one you've set
> on you `infra/docker/projects/am-env/config.local.yaml` file.

```
Host am-env-shell
    HostName 127.0.0.1
    User username
    Port 32533
    IdentityFile ~/works/abodemine/src/boot/infra/docker/projects/am-env/services/shell/.ssh/id_ed25519
    IdentitiesOnly yes
```

- Test the connection:

```shell
ssh am-env-shell
```

## Clone mono repo

```sh
cd /works/src
git clone git@github.com:abodeminehq/mono.git
```

## AWS SSO

Follow the steps in [aws-sso](aws-sso.md).

## Start Docker env

Copy `${ABODEMINE_WORKSPACE}/infra/docker/projects/am-mono/.env.example` to
`${ABODEMINE_WORKSPACE}/infra/docker/projects/am-mono/.env` and add
a (personal) secret value to `OPENSEARCH_INITIAL_ADMIN_PASSWORD`.

> You can use `pwgen -cns 12` inside the `am-env-shell` to generate some choices.

```sh
mkdir -p ${ABODEMINE_WORKSPACE}/code/sql/databases/api/migrations
make -C ${ABODEMINE_WORKSPACE}/infra/docker/projects/am-mono up
```

## Install certs

On the host computer:

```sh
scp am-env-shell:/works/src/mono/etc/ssl/abodemine-ca.pem ${HOME}/abodemine-ca.pem
```

### macOS

```sh
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ${HOME}/abodemine-ca.pem
```

## Add entries to /etc/hosts

Open the host /etc/hosts file and add:

```
127.0.0.1 admin.abodemine.test
127.0.0.1 api.abodemine.test
127.0.0.1 app.abodemine.test
127.0.0.1 whitelabel.abodemine.test
127.0.0.1 opensearch.abodemine.test
```

## Open remote dev containers

- Select the `am-mono-dc-root-1` container and enter the folder `/works/src/mono`.
- Install the extensions listed [here](vs-code.md). Ensure they are installed
  in the devcontainer, too, not only the host.
  > Extensions just need to be installed in a single devcontainer as they share
  > the same vscode-server folder internally. But you can disable them individually per
  > workspace (devcontainer).
- Select the `am-mono-dc-go-1` container and enter the folder `/works/src/mono/code/go`.

## Check you can reach web services

- Open https://app.abodemine.test with you browser.

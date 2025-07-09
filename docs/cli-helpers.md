# Reset gpg

```sh
gpgconf --kill gpg-agent && gpg-connect-agent updatestartuptty /bye >/dev/null && killall pinentry
echo "" | gpg --clear-sign
```

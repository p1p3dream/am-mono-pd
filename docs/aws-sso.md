## Configure

> This is required only once per (dev env) install. Future accesses should just login directly.

```sh
aws configure sso --no-browser --use-device-code
```

- Start URL: `https://d-9067f32245.awsapps.com/start`
- Region: `us-east-1`
- CLI region: `us-east-2`
- CLI output format: `json`

> The device code URL must be inserted in an authenticated browser session with AWS.

> Note that the LOGIN region is `us-east-1` because that's where Identity Center was installed;
> still, all RELEVANT assets are located in the `us-east-2` region.

With that done, you won't be required to enter AWS keys/secrets manually anywhere.
Just set the profile name where applicable and login when necessary.

See below how to list your profiles.

## List local profiles

```sh
aws configure list-profiles
```

## Login

> If you haven't done so before, [configure](#Configure) your SSO session.

Once every few hours you'll be required to login again:

```sh
aws sso login --no-browser --use-device-code --profile the_desired_profile_name
```

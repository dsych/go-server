1. Set-up the env `sudo ./configure-env.sh`
2. To set-up the certificates run `sudo ./generate-keys.sh`
3. Install binary using `go install`
4. Run `go-server`

  ## NOTE
* You need to set `GO_USERNAME`, `GO_PASSWORD`, `GO_DATABASE` and `GO_HOST` to corresponding mysql credentials, so that Go server can login into your database.

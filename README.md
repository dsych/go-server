1. Set `CERTS_PATH` env varible that points to where cert and key is going to be stored e.g. `export CERTS_PATH=$HOME/.certs`
2. To set-up the environment run `sudo ./generate-keys.sh`
3. Install binary using `go install`
4. Restart Apache `sudo service httpd restart`
5. Add exception for port `1443` in firewall `sudo iptables -I INPUT 2 -p tcp -m state --state NEW -m tcp --dport 1443 -j ACCEPT`
6. Run `go-server`

* To log in you need to pass `user` and `password` as query parameters to login route. List of login and passwords could be found in the code.

  
  ## NOTE
* You need to set `GO_USERNAME`, `GO_PASSWORD`, `GO_DATABASE` and `GO_HOST` to corresponding mysql credentials, so that Go server can login into your database.

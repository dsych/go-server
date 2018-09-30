1. To generate certificates run `sudo generate-keys.sh`
2. Install binary using `go install`
3. Restart Apache `sudo service httpd restart`
4. Add exception for port `1443` in firewall `sudo iptables -I INPUT 2 -p tcp -m state --state NEW -m tcp --dport 1443 -j ACCEPT`
4. Run `go-server`

* To log in you need to pass `user` and `password` as query parameters to login route. List of login and passwords could be found in the code.
* Routes:
  * `/login`
  * `content/hello` - protected by login
  * `content/logout` - protected by login
  * `/hello`

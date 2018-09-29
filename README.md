1. To generate certificates run `generate-keys.sh`
2. Install binary using `go install`
3. Run `go-server`

* To log in you need to pass `user` and `password` as query parameters to login route. List of login and passwords could be found in the code.
* Routes:
  * `/login`
  * `content/hello` - protected by login
  * `content/logout` - protected by login
  * `/hello`

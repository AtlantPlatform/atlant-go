#AuthCenter

##DNSAuth

DNS authorization means that for checking presence
of node write/sync permissions we are requesting TXT records of specific domains

Permissions are stored in the following form:

```
$ dig node-test.atlant-dev.io TXT | grep TXT | grep -v \;

node-test.atlant-dev.io. 96	IN	TXT	"14V8BTKR9MjKhfqgT4ybBjSb7kZHXmwvgBba7ujhN6ecTXJji:sync"
node-test.atlant-dev.io. 119	IN	TXT	"14V8BbA8ipE7jqE9a4CfLfTzwVKayV5GJsjP4c9gWVB5ZDSww:write,sync"
node-test.atlant-dev.io. 119	IN	TXT	"14V8BTCcnRXc3m28j5eEkjoGi5w31MjAB6yKXAumbuA3spsUK:write"
````

This is default authorization method. Domains can be set up with `testnet-auth-domains` parameter (env. `AN_TESTNET_DOMAINS`)

[WARNING] A tricky moment is that your DNS resolver might not support multiple TXT records. As a workaround, using `8.8.8.8` and `8.8.4.4` in you `/etc/resolv.conf` is recomended.

## UrlAuth

URL authorization means that for checking presence
of node write/sync permissions we are requesting some HTTP URL.

Sending GET request to this URL should return same list of nodes with their permissions

```
$ curl -q http://localhost:30100/

14V8BTKR9MjKhfqgT4ybBjSb7kZHXmwvgBba7ujhN6ecTXJji:sync
14V8BbA8ipE7jqE9a4CfLfTzwVKayV5GJsjP4c9gWVB5ZDSww:write,sync
14V8BTCcnRXc3m28j5eEkjoGi5w31MjAB6yKXAumbuA3spsUK:write
```

This method is good for local networks and testing environments.
URLs can be set up with `testnet-auth-urls` parameter (env. `AN_TESTNET_URLS`)
#AuthCenter

##DNSAuth

DNS authorization means that for checking presence
of node write/sync permissions we are requesting TXT records of specific domains

Permissions are stored in the following form:

```
$ dig node-test.atlant-dev.io TXT | grep TXT | grep -v \;

node-test.atlant-dev.io. 96	IN	TXT	"14V8BdHqHhExw4645xB3Xa2iheBrjYCMr7StXWUA9hBTqp8cM:sync"
node-test.atlant-dev.io. 96	IN	TXT	"14V8Bds64aUZJx6ag2TUXozS78Sko6fJ8kbkHF4bgvv9zgR6j:sync"
node-test.atlant-dev.io. 96	IN	TXT	"14V8BVs2FyU5qREKd68SgPqccrChiWX2uKdeeMtUhGfqJZjyK:sync"
node-test.atlant-dev.io. 96	IN	TXT	"14V8BTKR9MjKhfqgT4ybBjSb7kZHXmwvgBba7ujhN6ecTXJji:sync"
node-test.atlant-dev.io. 119	IN	TXT	"14V8BbA8ipE7jqE9a4CfLfTzwVKayV5GJsjP4c9gWVB5ZDSww:write,sync"
node-test.atlant-dev.io. 119	IN	TXT	"14V8BTCcnRXc3m28j5eEkjoGi5w31MjAB6yKXAumbuA3spsUK:write"
node-test.atlant-dev.io. 119	IN	TXT	"14V8BjQJ2E2xyQPC3FCFA5Aa9cMxBGpCkeDTPsYLxDqb7Xfne:write"
node-test.atlant-dev.io. 119	IN	TXT	"14V8BXzirg4YbiMda1e9b42RpJx5oWFwHmR2gXngzmFVbJR7Y:write"
node-test.atlant-dev.io. 119	IN	TXT	"14V8BSuoSQycJmwAqfbXb48CFrcAVhHwnGw1gx1cxrev1agh3:write"
````

[WARNING] A tricky moment is that your DNS resolver might not support multiple TXT records. As a workaround, using `8.8.8.8` and `8.8.4.4` in you `/etc/resolv.conf` is recomended.

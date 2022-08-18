# vault-up

simple util used for measuring uptime

* `/` (returns "yes\n" and HTTP 200)
* `/for` (returns current uptime of process in seconds)

## options
```
Usage:
      --bind-address string          ip:port where http requests are served (default ":8080")
      --default-response string      what to respond when receiving requests on '/' (default "yes\n")
```

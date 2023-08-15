# vault-up

simple app used for measuring vault uptime

* `/` (returns "yes\n" and HTTP 200)
* `/for` (returns current uptime of process in seconds)

## Options
```
Usage:
      --bind-address string          ip:port where http requests are served (default ":8080")
      --default-response string      what to respond when receiving requests on '/' (default "yes\n")
```

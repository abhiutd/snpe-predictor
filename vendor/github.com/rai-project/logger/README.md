# Logger [![Build Status](https://travis-ci.org/rai-project/logger.svg?branch=master)](https://travis-ci.org/rai-project/logger)

## Config

example

```
logger:
  level: debug
  hooks:
    - kenisis
    - syslog
    - logz
    - graylog
    - airbrake
logz:
  token: LOGZ_TOKEN
graylog:
   address: ...
   port: 12201
airbrake:
   id: ...
   api_key: ...
```

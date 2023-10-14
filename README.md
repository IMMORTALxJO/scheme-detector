# scheme-detector
![Coverage](https://img.shields.io/badge/Coverage-99.4%25-brightgreen)

[![DeepSource](https://deepsource.io/gh/IMMORTALxJO/scheme-detector.svg/?label=active+issues&show_trend=true&token=VZ2SYgG49PAWTLKYzI-vb-1A)](https://deepsource.io/gh/IMMORTALxJO/scheme-detector/?ref=repository-badge)

Detect different protocols and engines from the current environment variables.

```
$ env
...
DATABASE_URI=pgsql://user:pg_pass@pg.example.com/example
EXAMPLE_API_PASS=apipass
EXAMPLE_API_URL=https://api.example.com
EXAMPLE_API_USER=apiuser
...
$ ./schemedetector
[
  {
    "engine": "pgsql",
    "port": "5432",
    "host": "pg.example.com",
    "username": "user",
    "password": "pg_pass"
  },
  {
    "engine": "https",
    "port": "443",
    "host": "api.example.com",
    "username": "apiuser",
    "password": "apipass"
  }
]
```

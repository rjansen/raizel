# raizel [![Build Status](https://travis-ci.org/rjansen/raizel.svg?branch=master)](https://travis-ci.org/rjansen/raizel) [![Coverage Status](https://codecov.io/gh/rjansen/raizel/branch/master/graph/badge.svg)](https://codecov.io/gh/rjansen/raizel) [![Go Report Card](https://goreportcard.com/badge/github.com/rjansen/raizel)](https://goreportcard.com/report/github.com/rjansen/raizel)

A persistence helper library

# dependencies
## tools (you must provide the installation)
- [Docker](https://www.docker.com/)

## libraries
- [gocql](github.com/gocql/gocql)
- [mysql-go-sql-driver](github.com/go-sql-driver/mysql)
- [psql-libpq](github.com/lib/pq)
- [zap](https://github.com/uber-go/zap)
- [viper](github.com/spf13/viper)
- [l](github.com/rjansen/l)
- [migi](github.com/rjansen/migi)

# tests and coverage
- run unit tests: `make docker.test`
- run coverage: `make docker.coverage.text`
- run html coverage: `make docker.coverage.html`

# raizel usage
Find some samples in the test files. A better usage section will be avaiable soon ...

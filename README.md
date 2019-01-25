# raizel [![Build Status](https://travis-ci.org/rjansen/raizel.svg?branch=master)](https://travis-ci.org/rjansen/raizel) [![Coverage Status](https://codecov.io/gh/rjansen/raizel/branch/master/graph/badge.svg)](https://codecov.io/gh/rjansen/raizel) [![Go Report Card](https://goreportcard.com/badge/github.com/rjansen/raizel)](https://goreportcard.com/report/github.com/rjansen/raizel)

A persistence helper library that supports Firestore and PostgreSQL

##### Cassandra and MySQL support will be avaiable soon

# dependencies
### tools (you must provide the installation)
- [Docker](https://www.docker.com/)

### libraries
- [firestore](https://godoc.org/cloud.google.com/go/firestore)
- [gocql](https://github.com/gocql/gocql)
- [mysql-go-sql-driver](https://github.com/go-sql-driver/mysql)
- [psql-libpq](https://github.com/lib/pq)
- [l](https://github.com/rjansen/l)
- [migi](https://github.com/rjansen/migi)
- [yggdrasil](https://github.com/rjansen/yggdrasil)

# tests and coverage
- run unit tests: `make docker.test`
- run coverage: `make docker.coverage.text`
- run html coverage: `make docker.coverage.html`

# raizel usage
Find some samples in the test files. A better usage section will be avaiable soon ...

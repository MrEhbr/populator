# populator

[![Go](https://github.com/MrEhbr/populator/actions/workflows/go.yml/badge.svg)](https://github.com/MrEhbr/populator/actions/workflows/go.yml)
[![License](https://img.shields.io/badge/license-Apache--2.0%20%2F%20MIT-%2397ca00.svg)](https://github.com/MrEhbr/populator/blob/master/COPYRIGHT)
[![GitHub release](https://img.shields.io/github/release/MrEhbr/populator.svg)](https://github.com/MrEhbr/populator/releases)
[![codecov](https://codecov.io/gh/MrEhbr/populator/branch/master/graph/badge.svg)](https://codecov.io/gh/MrEhbr/populator)
![Made by Alexey Burmistrov](https://img.shields.io/badge/made%20by-Alexey%20Burmistrov-blue.svg?style=flat)

This package was created for testing purposes, to give the ability to seed a database with records from simple .yaml files. Populator respects the order in files, so you can handle foreign_keys just by placing them in the right order.

## Install

### Using go

```console
go get -u github.com/MrEhbr/populator
```

## Example usage

```golang
package foo_test

import (
 "database/sql"
 "os"
 "testing"

 "github.com/MrEhbr/populator"
 "github.com/MrEhbr/populator/postgres"
)

var fixtures = `
- table: users
  rows:
    - id: 1
      name: "foo"
      settings:
        foo: "bar"
    - id: 2
      name: "bar"
      settings:
        - foo: "bar"
`

func TestFoo(t *testing.T) {
 db, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN"))
 if err != nil {
  t.Fatalf("Unable to connect to database: %v", err)
 }
 defer db.Close()

 engine := postgres.New(postgres.NewSQL(db), postgres.DisableForeignKeyCheck())
 seeder := populator.New(
  populator.WithEngine(engine),
  populator.WithParser(populator.YAMLParse),
 )

 if err := seeder.From(fixtures); err != nil {
  t.Fatalf("Unable to seed database: %v", err)
 }
}

```

## License

© 2020 [Alexey Burmistrov]

Licensed under the [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0) ([`LICENSE`](LICENSE)). See the [`COPYRIGHT`](COPYRIGHT) file for more details.

`SPDX-License-Identifier: Apache-2.0`

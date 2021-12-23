package populator_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/MrEhbr/populator/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/pkg/errors"
)

func TestPostgres(t *testing.T) {
	if strings.Contains(os.Getenv("OS"), "macos") {
		t.Skip("not supported on mac os runner")
	}

	var (
		db          *pgxpool.Pool
		databaseUrl string
	)
	{
		// uses a sensible default on windows (tcp/http) and linux/osx (socket)
		pool, err := dockertest.NewPool("")
		if err != nil {
			if errors.Is(err, docker.ErrInvalidEndpoint) {
				t.Skip("docker endpoint not found")
			}

			t.Fatalf("Could not connect to docker: %s", err)
		}
		if _, err := pool.Client.Info(); err != nil {
			if errors.Is(err, docker.ErrConnectionRefused) {
				t.Skip("docker not running")
			}
		}

		// pulls an image, creates a container based on it and runs it
		resource, err := pool.RunWithOptions(&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "11",
			Env: []string{
				"POSTGRES_PASSWORD=secret",
				"POSTGRES_USER=user_name",
				"POSTGRES_DB=test",
				"listen_addresses = '*'",
			},
		}, func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		})
		if err != nil {
			t.Fatalf("Could not start resource: %s", err)
		}

		databaseUrl = fmt.Sprintf("postgres://user_name:secret@%s/test?sslmode=disable", resource.GetHostPort("5432/tcp"))

		t.Logf("Connecting to database on url: %s", databaseUrl)

		resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

		// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
		pool.MaxWait = 120 * time.Second
		if err = pool.Retry(func() error {
			db, err = pgxpool.Connect(context.Background(), databaseUrl)
			if err != nil {
				return err
			}
			return db.Ping(context.Background())
		}); err != nil {
			t.Fatalf("Could not connect to postgres: %s", err)
		}

		t.Cleanup(func() {
			if err := pool.Purge(resource); err != nil {
				t.Fatalf("Could not purge resource: %s", err)
			}
		})
	}

	testWithDriver(postgres.NewPgx(db), t)

	stdDB, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		t.Fatalf("Unable to connect to database: %v", err)
	}

	testWithDriver(postgres.NewSQL(stdDB), t)
}

func testWithDriver(conn postgres.Conn, t *testing.T) {
	pgScheme := `
drop table if exists role;
drop table if exists users;
create table if not exists users
(
    id   serial constraint users_pk primary key,
    name text,
    settings JSON
);
create table if not exists role
(
    name    text,
    user_id int constraint role_user_id_fk references users(id),
    attrs int[]
);`
	prepareFn := func(conn postgres.Conn) func() error {
		return func() error {
			if _, err := conn.Exec(context.Background(), pgScheme); err != nil {
				return fmt.Errorf("failed to create schema: %w", err)
			}

			return nil
		}
	}

	t.Run("pgx driver", func(t *testing.T) {
		t.Run("simple", func(t *testing.T) {
			engine := postgres.New(conn)

			cases := []testCase{
				{
					name: "positive",
					fixtures: []string{
						`- table: users
  rows:
    - id: 1
      name: "foo"
      settings:
        foo: "bar"
    - id: 2
      name: "bar"
      settings:
        - foo: "bar"
- table: role
  rows:
  - user_id: 1
    name: "test_foo"
    attrs: [1,2,3]
  - user_id: 2
    name: "test_bar"`,
					},
					wantErr: false,
					prepare: prepareFn(conn),
				},
				{
					name: "table not exists",
					fixtures: []string{
						`- table: unknown
  rows:
  - id: 1
    name: "foo"
  - id: 2
    name: "bar"`,
					},
					wantErr: true,
					prepare: prepareFn(conn),
				},
			}
			testWithEngine(engine, cases, t)
		})

		t.Run("disabled forign key check", func(t *testing.T) {
			engine := postgres.New(conn, postgres.DisableForeignKeyCheck())

			cases := []testCase{
				{
					name: "positive",
					fixtures: []string{
						`- table: role
  rows:
  - user_id: -1
    name: "test_foo"
  - user_id: -2
    name: "test_bar"`,
					},
					wantErr: false,
					prepare: prepareFn(conn),
				},
			}
			testWithEngine(engine, cases, t)
		})
	})
}

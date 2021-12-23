package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/MrEhbr/populator"
	"github.com/MrEhbr/populator/helpers"
	"github.com/jackc/pgtype"
)

type Conn interface {
	BeginFunc(ctx context.Context, f func(Conn) error) (err error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (Rows, error)
}

type Rows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...interface{}) error
}

type Engine struct {
	conn                   Conn
	disableForeignKeyCheck bool
}

func New(conn Conn, options ...Option) *Engine {
	engine := &Engine{conn: conn}

	for _, o := range options {
		o(engine)
	}

	return engine
}

func (pg *Engine) Build(fixtures populator.Fixtures) ([]populator.Command, error) {
	cmds := make([]populator.Command, 0)
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	for _, f := range fixtures {
		for _, v := range f.Rows {
			fields, err := prepareFields(v)
			if err != nil {
				return nil, fmt.Errorf("failed to prepare fields: %w", err)
			}
			cmd := populator.Command{}
			cmd.Query, cmd.Args = psql.Insert(quoteKeyword(f.Table)).SetMap(fields).MustSql()
			cmds = append(cmds, cmd)
		}
	}

	if pg.disableForeignKeyCheck {
		alterCmds := make([]populator.Command, 0)
		for _, table := range fixtures.Tables() {
			alterCmds = append(alterCmds, populator.Command{
				Query: fmt.Sprintf("ALTER TABLE %s DISABLE TRIGGER ALL;", quoteKeyword(table)),
			})

			cmds = append(cmds, populator.Command{
				Query: fmt.Sprintf("ALTER TABLE %s ENABLE TRIGGER ALL;", quoteKeyword(table)),
			})
		}
		cmds = append(alterCmds, cmds...)
	}

	return cmds, nil
}

func (pg *Engine) Exec(cmds []populator.Command) error {
	return pg.conn.BeginFunc(context.Background(), func(conn Conn) error {
		for _, cmd := range cmds {
			_, err := conn.Exec(context.Background(), cmd.Query, cmd.Args...)
			if err != nil {
				return fmt.Errorf("failed to execute %q with args %v :%w", cmd.Query, cmd.Args, err)
			}
		}
		return nil
	})
}

func quoteKeyword(s string) string {
	parts := strings.Split(s, ".")
	for i, p := range parts {
		parts[i] = fmt.Sprintf(`"%s"`, p)
	}
	return strings.Join(parts, ".")
}

// nolint: exhaustive
func prepareFields(fields map[string]interface{}) (map[string]interface{}, error) {
	for k, v := range fields {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Map:
			var vv pgtype.JSON
			raw, err := json.Marshal(toJSON(v))
			if err != nil {
				return nil, err
			}
			_ = vv.Set(raw)

			fields[k] = vv
		case reflect.Slice:
			fields[k] = toJSON(v)
		case reflect.String:
			if t, err := helpers.TryStrToDate(v.(string)); err == nil {
				fields[k] = t
			} else if b, err := helpers.TryHexStringToBytes(v.(string)); err == nil {
				fields[k] = b
			}
		}
	}
	return fields, nil
}

func toJSON(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		for i, el := range v {
			v[i] = toJSON(el)
		}
		return v
	case map[interface{}]interface{}:
		newMap := make(map[string]interface{}, len(v))
		for k, e := range v {
			newMap[k.(string)] = toJSON(e)
		}

		return newMap
	}

	return v
}

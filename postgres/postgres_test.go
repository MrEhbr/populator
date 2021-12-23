package postgres

import (
	"testing"

	"github.com/MrEhbr/populator"
	"github.com/matryer/is"
)

func TestPgx_Build(t *testing.T) {
	t.Run("check generated commands", func(t *testing.T) {
		engine := New(nil)
		cmds, err := engine.Build(populator.Fixtures{
			{
				Table: "foo",
				Rows: []map[string]interface{}{
					{"a": "str", "b": []int{1, 2, 3}, "c": map[string]string{"key": "val"}},
				},
			},
			{
				Table: "foo.bar",
				Rows: []map[string]interface{}{
					{"a": "str", "b": []int{1, 2, 3}, "c": map[string]string{"key": "val"}},
				},
			},
		})
		is := is.New(t)
		is.NoErr(err)
		is.Equal(len(cmds), 2) // must be 2 commands
		is.Equal(`INSERT INTO "foo" (a,b,c) VALUES ($1,$2,$3)`, cmds[0].Query)
		is.Equal(3, len(cmds[0].Args))

		is.Equal(`INSERT INTO "foo"."bar" (a,b,c) VALUES ($1,$2,$3)`, cmds[1].Query)
		is.Equal(3, len(cmds[1].Args))
	})

	t.Run("check alter commands", func(t *testing.T) {
		engine := New(nil, DisableForeignKeyCheck())
		cmds, err := engine.Build(populator.Fixtures{
			{
				Table: "foo",
				Rows: []map[string]interface{}{
					{"a": "str", "b": []int{1, 2, 3}, "c": map[string]string{"key": "val"}},
				},
			},
			{
				Table: "foo.bar",
				Rows: []map[string]interface{}{
					{"a": "str", "b": []int{1, 2, 3}, "c": map[string]string{"key": "val"}},
				},
			},
		})
		is := is.New(t)
		is.NoErr(err)
		is.Equal(len(cmds), 6) // must be 6 commands
		is.Equal(`ALTER TABLE "foo" DISABLE TRIGGER ALL;`, cmds[0].Query)
		is.Equal(nil, cmds[0].Args)

		is.Equal(`ALTER TABLE "foo"."bar" DISABLE TRIGGER ALL;`, cmds[1].Query)
		is.Equal(nil, cmds[1].Args)

		is.Equal(`ALTER TABLE "foo" ENABLE TRIGGER ALL;`, cmds[len(cmds)-2].Query)
		is.Equal(nil, cmds[len(cmds)-2].Args)

		is.Equal(`ALTER TABLE "foo"."bar" ENABLE TRIGGER ALL;`, cmds[len(cmds)-1].Query)
		is.Equal(nil, cmds[len(cmds)-1].Args)
	})
}

package populator

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type (
	Fixture struct {
		Table string
		Rows  []map[string]interface{}
	}
	Fixtures []Fixture

	Command struct {
		Query string
		Args  []interface{}
	}

	Engine interface {
		Build(Fixtures) ([]Command, error)
		Exec(cmds []Command) error
	}

	Parser func(r io.Reader) (Fixtures, error)

	Populator struct {
		parser Parser
		engine Engine
	}
)

func (f Fixtures) Tables() []string {
	tables := make([]string, 0)
	seen := make(map[string]struct{})
	for _, v := range f {
		if _, ok := seen[v.Table]; !ok {
			tables = append(tables, v.Table)
			seen[v.Table] = struct{}{}
		}
	}

	return tables
}

func New(options ...Option) *Populator {
	p := &Populator{
		parser: YAMLParse,
	}

	for _, o := range options {
		o(p)
	}

	return p
}

func (p *Populator) Load(file ...string) error {
	var fixtures Fixtures
	for _, fName := range file {
		f, err := os.Open(fName)
		if err != nil {
			return fmt.Errorf("failed to open file %q: %w", fName, err)
		}

		defer func() {
			_ = f.Close()
		}()

		v, err := p.parser(f)
		if err != nil {
			return fmt.Errorf("failed to parse file %q: %w", fName, err)
		}

		fixtures = append(fixtures, v...)
	}

	cmds, err := p.engine.Build(fixtures)
	if err != nil {
		return fmt.Errorf("failed to build commands: %w", err)
	}

	if err := p.engine.Exec(cmds); err != nil {
		return fmt.Errorf("failed to execute commands: %w", err)
	}

	return nil
}

func (p *Populator) From(content string) error {
	data, err := p.parser(strings.NewReader(content))
	if err != nil {
		return fmt.Errorf("failed to parse: %w", err)
	}

	cmds, err := p.engine.Build(data)
	if err != nil {
		return fmt.Errorf("failed to build commands: %w", err)
	}
	if err := p.engine.Exec(cmds); err != nil {
		return fmt.Errorf("failed to execute commands: %w", err)
	}

	return nil
}

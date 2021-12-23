package populator_test

import (
	"io/ioutil"
	"path"
	"reflect"
	"strconv"
	"testing"

	"github.com/MrEhbr/populator"
)

type testCase struct {
	name     string
	fixtures []string
	wantErr  bool
	prepare  func() error
}

func testWithEngine(engine populator.Engine, testCase []testCase, t *testing.T) {
	populator := populator.New(populator.WithEngine(engine), populator.WithParser(populator.YAMLParse))
	t.Run("From", func(t *testing.T) {
		for _, tt := range testCase {
			tt := tt
			if err := tt.prepare(); err != nil {
				t.Errorf("tt.prepare() error = %v ", err)
			}
			for _, f := range tt.fixtures {
				if err := populator.From(f); (err != nil) != tt.wantErr {
					t.Errorf("Case: %s, Populator.From() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				}
			}
		}
	})

	t.Run("Load", func(t *testing.T) {
		for _, tt := range testCase {
			tt := tt
			fNames := make([]string, 0)
			dir := t.TempDir()
			for i, f := range tt.fixtures {
				fname := path.Join(dir, strconv.Itoa(i)+".yaml")
				if err := ioutil.WriteFile(fname, []byte(f), 0600); err != nil {
					t.Fatalf("failed to write temp file: %s", err)
				}

				fNames = append(fNames, fname)
			}
			if err := tt.prepare(); err != nil {
				t.Errorf("tt.prepare() error = %v ", err)
			}
			if err := populator.Load(fNames...); (err != nil) != tt.wantErr {
				t.Errorf("Case: %s, Populator.Load() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		}
	})
}

func TestFixtures_Tables(t *testing.T) {
	tests := []struct {
		name string
		f    populator.Fixtures
		want []string
	}{
		{
			name: "no duplicates",
			f:    populator.Fixtures{{Table: "foo"}, {Table: "bar"}},
			want: []string{"foo", "bar"},
		},
		{
			name: "with duplicates",
			f:    populator.Fixtures{{Table: "foo"}, {Table: "bar"}, {Table: "bar"}},
			want: []string{"foo", "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Tables(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fixtures.Tables() = %v, want %v", got, tt.want)
			}
		})
	}
}

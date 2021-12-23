package populator

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func JSONParse(r io.Reader) (Fixtures, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read")
	}

	var fixtures Fixtures
	if err := json.Unmarshal(data, &fixtures); err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}

	return fixtures, nil
}

func YAMLParse(r io.Reader) (Fixtures, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}

	var fixtures Fixtures
	if err := yaml.Unmarshal(data, &fixtures); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return fixtures, nil
}

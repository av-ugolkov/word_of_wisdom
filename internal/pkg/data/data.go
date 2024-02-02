package data

import (
	"fmt"
	"strings"
)

type Key string

const (
	Challenge Key = "challenge"
	Response  Key = "response"
	Grant     Key = "grant"
)

type Data struct {
	Key   Key
	Value string
}

func Parse(str string) (*Data, error) {
	str = strings.TrimSpace(str)
	d := strings.Split(str, "::")
	if len(d) != 2 {
		return nil, fmt.Errorf("pkg.data.Parse - invalid data: %s", str)
	}
	return &Data{Key: Key(d[0]), Value: d[1]}, nil
}

func (d *Data) String() string {
	return fmt.Sprintf("%s::%s", d.Key, d.Value)
}

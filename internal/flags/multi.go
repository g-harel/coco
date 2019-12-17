package flags

import (
	"flag"
	"fmt"
)

// MultiFlag is a custom flag type which supports passing in
// the same flag multiple times. All instances are
// accumulated into a slice of strings.
type multiFlag []string

var _ flag.Value = &multiFlag{}

func (m *multiFlag) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

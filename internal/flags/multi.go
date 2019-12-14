package flags

import (
	"fmt"
)

type multiFlag []string

func (m *multiFlag) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *multiFlag) Set(value string) error {
	for i := 0; i < len(*m); i++ {
		if value == (*m)[i] {
			return nil
		}
	}
	*m = append(*m, value)
	return nil
}

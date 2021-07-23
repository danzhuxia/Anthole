package anthole

import "fmt"

var (
	pName   = `Anthole`
	version = `0.0.1`

	Version = func() string {
		return fmt.Sprintf("%s-%s", pName, version)
	}
)

package ucerr

import "fmt"

var (
	InvalidArgument = fmt.Errorf("invalid argument")
	NotFound        = fmt.Errorf("not found")
	AlreadyExists   = fmt.Errorf("already exists")
)

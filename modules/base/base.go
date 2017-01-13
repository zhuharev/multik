package base

import (
	"github.com/ungerik/go-dry"
)

func Sites() ([]string, error) {
	return dry.ListDirDirectories("sites")
}

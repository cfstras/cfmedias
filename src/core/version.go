package core

import "fmt"

type version struct {
	major   int
	minor   int
	patch   int
	comment string
}

var currentVersion version

func init() {
	currentVersion = version{0, 1, 0, "pre"}
}

func (v version) String() string {
	return fmt.Sprintf("%d.%d.%d-%s", v.major, v.minor, v.patch, v.comment)
}

func Version() string {
	return currentVersion.String()
}

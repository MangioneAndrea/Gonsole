package gonsole

import (
	"testing"
)

func TestCli_Confirm(t *testing.T) {
	Cli{}.Confirm("hello")
}

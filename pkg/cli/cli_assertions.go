package cli

import (
	"strings"
	"testing"

	"github.com/cirss/go-cli/pkg/util"
)

func (pc *ProgramContext) AssertExitCode(t *testing.T, commandLine string, expected int) {
	invokedName := strings.Fields(commandLine)[0]
	if invokedName != pc.Name {
		t.Fatal("Unexpected program name: ", invokedName)
	}
	actual := pc.Invoke(commandLine)
	util.IntEquals(t, actual, expected)
}

func (pc *ProgramContext) AssertSuccess(t *testing.T, commandLine string) {
	pc.AssertExitCode(t, commandLine, 0)
}

package tests

import (
	"fmt"
	"testing"

	"github.com/cirss/gocli/go/cli"
)

func main() {

}

func TestCLI(t *testing.T) {

	programContext := cli.NewProgramContext("main", main)
	fmt.Printf("%v\n", programContext)
}

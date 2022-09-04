package wasm

import (
	//	"github.com/errata-ai/vale/v2/internal/check"
	//	"github.com/errata-ai/vale/v2/internal/core"
	"fmt"

	//"github.com/errata-ai/vale/v2/internal/core"
	"github.com/errata-ai/vale/v2/internal/core"
	"github.com/errata-ai/vale/v2/internal/lint"
	// "github.com/errata-ai/vale/v2/internal/nlp"
	// "github.com/jdkato/prose/tag"
	// "github.com/pterm/pterm"
)

func main() {

	runTest([]string{"Test"})
}

func runTest(args []string) error {

	cfg, err := core.NewConfig(&core.CLIFlags{})
	if err != nil {
		panic(err)
	}
	cfg.GBaseStyles = []string{"Vale"}

	cfg.MinAlertLevel = 0
	cfg.GBaseStyles = []string{"Test"}
	cfg.Flags.InExt = ".txt" // default value

	linter, err := lint.NewLinter(cfg)
	if err != nil {
		return err
	}

	linted, err := linter.LintString("This is just a string to be linted")
	if err != nil {
		return err
	}

	fmt.Println(linted)

	/*
		    linted, err := linter.Lint([]string{args[1]}, "*")
			if err != nil {
				return err
			}
	*/

	/*
		    PrintJSONAlerts(linted)
			return nil
	*/
	return nil

}

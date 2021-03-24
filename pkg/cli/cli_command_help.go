package cli

import (
	"fmt"
	"io"
)

func Help(cc *CommandContext) (err error) {
	if len(cc.Args) < 2 {
		cc.ShowProgramUsage(cc.OutWriter)
		cc.ShowProgramCommands(cc.OutWriter)
		return
	}
	commandName := cc.Args[1]
	if commandName == "help" {
		return
	}
	if c, exists := cc.Commands.Lookup(commandName); exists {
		cc.Descriptor = c
		cc.Args = []string{commandName, "help"}
		c.Handler(cc)
	} else {
		fmt.Fprintf(cc.ErrWriter, "%s help: unrecognized %s command: %s\n\n",
			cc.programContext.Name, cc.programContext.Name, commandName)
		cc.ShowProgramUsage(cc.OutWriter)
		cc.ShowProgramCommands(cc.OutWriter)
		err = NewCLIError(nil)
	}
	return
}

func (cc *CommandContext) ShowHelpIfRequested(w io.Writer) bool {
	if len(cc.Args) > 1 && cc.Args[1] == "help" {
		cc.ShowCommandDescription(w)
		cc.ShowCommandUsage(w)
		return true
	}
	return false
}

func (cc *CommandContext) ShowCommandDescription(w io.Writer) {
	fmt.Fprintf(w, "%s %s: %s\n",
		cc.programContext.Name, cc.Descriptor.Name, cc.Descriptor.Description)
}

func (cc *CommandContext) ShowCommandUsage(w io.Writer) {
	fmt.Fprintf(w, "\nusage: %s %s [<flags>]\n\n", cc.programContext.Name, cc.Descriptor.Name)
	fmt.Fprint(w, "flags:\n")
	cc.Flags.SetOutput(w)
	cc.Flags.PrintDefaults()
	// fmt.Fprintln(w)
}

func (cc *CommandContext) ShowProgramUsage(w io.Writer) {
	fmt.Fprintf(w, "usage: %s <command> [<flags>]\n\n", cc.programContext.Name)
}

func (cc *CommandContext) ShowProgramCommands(w io.Writer) {
	fmt.Fprint(w, "commands:\n")
	for _, sc := range cc.Commands.commandList {
		fmt.Fprintf(cc.OutWriter, "  %-7s  - %s\n", sc.Name, sc.Summary)
	}
	fmt.Fprint(w, "\nflags:\n")
	cc.Flags.PrintDefaults()
	fmt.Fprintf(w,
		"\nSee '%s help <command>' for help with one of the above commands.\n",
		cc.programContext.Name)
	return
}

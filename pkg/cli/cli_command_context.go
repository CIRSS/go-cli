package cli

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type NullWriter struct {
	w io.Writer
}

func (nw NullWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

type CommandContext struct {
	programContext *ProgramContext
	Commands       *CommandSet
	Descriptor     *CommandDescriptor
	Args           []string
	Flags          *flag.FlagSet
	Quiet          *bool
	Silent         *bool
	InReader       io.Reader
	OutWriter      io.Writer
	ErrWriter      io.Writer
	Properties     map[string]interface{}
	Providers      map[string]func(cc *CommandContext) interface{}
}

func NewCommandContext(pc *ProgramContext, commands *CommandSet) (cc *CommandContext) {

	cc = new(CommandContext)
	cc.programContext = pc
	cc.Commands = commands

	cc.Flags = pc.InitFlagSet()
	cc.Flags.Usage = func() {}

	cc.InReader = pc.InReader
	cc.ErrWriter = pc.ErrWriter
	cc.OutWriter = pc.OutWriter

	cc.Properties = make(map[string]interface{})
	cc.Providers = make(map[string]func(cc *CommandContext) interface{})

	cc.Quiet = cc.Flags.Bool("quiet", false, "Discard normal command output")
	cc.Silent = cc.Flags.Bool("silent", false, "Discard normal and error command output")

	return
}

func (cc *CommandContext) AddProvider(resource string, f func(cc *CommandContext) interface{}) {
	cc.Providers[resource] = f
}

func (cc *CommandContext) Resource(resourceName string) (resource interface{}) {
	provider, exists := cc.Providers[resourceName]
	if !exists {
		panic("No resource provider for " + resourceName)
	}
	resource = provider(cc)
	return
}

func (cc *CommandContext) ParseCommandFlags() (err error) {

	var errBuffer = new(strings.Builder)

	cc.Flags.SetOutput(errBuffer)
	if err = cc.Flags.Parse(cc.Args[1:]); err != nil {
		err = NewCLIError(err)
		fmt.Fprintf(cc.ErrWriter, "%s %s: %s",
			cc.programContext.Name, cc.Descriptor.Name, errBuffer.String())
		cc.ShowCommandUsage(cc.ErrWriter)

		return
	}

	if *cc.Quiet {
		cc.OutWriter = NullWriter{}
	}

	if *cc.Silent {
		cc.OutWriter = NullWriter{}
		cc.ErrWriter = NullWriter{}
	}

	return
}

func (cc *CommandContext) ParseFlags() (helpShown bool, err error) {

	if cc.ShowHelpIfRequested(cc.OutWriter) {
		helpShown = true
		return
	}

	err = cc.ParseCommandFlags()
	if err != nil {
		return
	}

	if len(cc.Flags.Args()) > 0 {
		fmt.Fprintf(cc.ErrWriter, "%s %s: unused argument: %s\n",
			cc.programContext.Name, cc.Descriptor.Name, cc.Flags.Args()[0])
		cc.ShowCommandUsage(cc.ErrWriter)
		err = NewCLIError(nil)
		return
	}

	return
}

func (cc *CommandContext) InvokeCommand(args []string) {

	if len(args) < 2 {
		fmt.Fprintf(cc.programContext.ErrWriter, "%s: no command given\n\n",
			cc.programContext.Name)
		cc.ShowProgramUsage(cc.OutWriter)
		cc.ShowProgramCommands(cc.OutWriter)
		cc.programContext.ExitIfNonzero(1)
		return
	}

	commandName := args[1]
	descriptor, exists := cc.Commands.Lookup(commandName)
	cc.Descriptor = descriptor
	if !exists {
		fmt.Fprintf(cc.programContext.ErrWriter, "%s: unrecognized command: %s\n\n",
			cc.programContext.Name, commandName)
		cc.ShowProgramUsage(cc.OutWriter)
		cc.ShowProgramCommands(cc.OutWriter)
		cc.programContext.ExitIfNonzero(1)
		return
	}

	cc.Args = args[1:]
	err := cc.Descriptor.Handler(cc)

	if err != nil {
		switch err.(type) {
		case CLIError:
			break
		default:
			fmt.Fprintf(cc.ErrWriter, "%s %s: %s\n",
				cc.programContext.Name, cc.Descriptor.Name, err.Error())
		}
		cc.programContext.ExitIfNonzero(1)
		return
	}
}

func (cc *CommandContext) ReadFileOrStdin(filePath string) (bytes []byte, err error) {
	var r io.Reader
	if filePath == "-" {
		r = cc.InReader
	} else {
		r, _ = os.Open(filePath)
	}
	return ioutil.ReadAll(r)
}

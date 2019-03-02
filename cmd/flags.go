package cmd

import (
	"errors"
)

var errMissingFlagId = errors.New("missing flag block ID")

func init() {
	register(&flagsCmd{})
}

type flagsCmd struct {
	Add    addFlagsCmd `command:"add" description:"Add a thread flag"`
	List   lsFlagsCmd  `command:"ls" description:"List thread flags"`
	Get    getFlagsCmd `command:"get" description:"Get a thread flag"`
	Ignore rmFlagsCmd  `command:"ignore" description:"Ignore a thread flag"`
}

func (x *flagsCmd) Name() string {
	return "flags"
}

func (x *flagsCmd) Short() string {
	return "Manage thread flags"
}

func (x *flagsCmd) Long() string {
	return `
Flags are added as blocks in a thread, which target
another block, usually a file(s).
Use this command to add, list, get, and ignore flags.
`
}

type addFlagsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Block  string        `required:"true" short:"b" long:"block" description:"Thread block ID. Usually a file(s) block."`
}

func (x *addFlagsCmd) Usage() string {
	return `

Adds a flag to a thread block.`
}

func (x *addFlagsCmd) Execute(args []string) error {
	setApi(x.Client)

	res, err := executeJsonCmd(POST, "blocks/"+x.Block+"/flags", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type lsFlagsCmd struct {
	Client ClientOptions `group:"Client Options"`
	Block  string        `required:"true" short:"b" long:"block" description:"Thread block ID. Usually a file(s) block."`
}

func (x *lsFlagsCmd) Usage() string {
	return `

Lists flags on a thread block.`
}

func (x *lsFlagsCmd) Execute(args []string) error {
	setApi(x.Client)

	res, err := executeJsonCmd(GET, "blocks/"+x.Block+"/flags", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type getFlagsCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *getFlagsCmd) Usage() string {
	return `

Gets a thread flag by block ID.`
}

func (x *getFlagsCmd) Execute(args []string) error {
	setApi(x.Client)
	if len(args) == 0 {
		return errMissingFlagId
	}

	res, err := executeJsonCmd(GET, "blocks/"+args[0]+"/flag", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

type rmFlagsCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *rmFlagsCmd) Usage() string {
	return `

Ignores a thread flag by its block ID.
This adds an "ignore" thread block targeted at the flag.
Ignored blocks are by default not returned when listing. 
`
}

func (x *rmFlagsCmd) Execute(args []string) error {
	setApi(x.Client)
	return callRmBlocks(args)
}

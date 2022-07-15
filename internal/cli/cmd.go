package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"
	"unsafe"
)

type Intercept func(Ctx) error

type Cmd struct {
	Name     string
	Version  string
	About    string
	Example  string
	Wildcard bool
	Pos      []string
	Flags    []Flag
	cmds     []*Cmd
	Run      func(Ctx) error

	set *flag.FlagSet
	// root   *Cmd
	parent string

	Intercept Intercept
}

func (c *Cmd) Init() {
	c.init("", c, c.Flags)
}

func (c *Cmd) RunCtx(x context.Context) error {
	// find the index of the first flag
	args := os.Args
	i, help := seperate(args)

	// find the subcommand. offset is the number of
	// position arguments of the subcommand.
	cmd, offset, err := find(c, args[:i], help)
	if err != nil {
		return err
	}
	if cmd == nil {
		return errors.New("no such command")
	}

	// print the help information of this subcommand
	// if it has no user defined function
	if cmd.Run == nil {
		return cmd.help(os.Stdout)
	}

	// parse flags if we've got any
	if i < len(args) {
		if err := cmd.set.Parse(args[i:]); err != nil {
			return err
		}
	}

	// see if we are missing any flags
	if err := cmd.validate(); err != nil {
		return err
	}

	// build the context that goes through the entire program
	context := &ctx{
		Context: x,
		cmd:     cmd,
		pos:     args[i-offset : i],
	}

	// run interceptor if provided
	if c.Intercept != nil {
		if err := c.Intercept(context); err != nil {
			return err
		}
	}

	// execute the user defined function
	return cmd.Run(context)
}

func (c *Cmd) Register(cmd *Cmd) {
	c.cmds = append(c.cmds, cmd)
}

// finds the subcommand matching the given arguments from top-down.
func find(root *Cmd, args []string, help bool) (found *Cmd, offset int, err error) {
	i := 0
	n := len(args)
	// use a dummy node to normalize the searching process.
	next := &Cmd{cmds: []*Cmd{root}}

search:
	for _, cmd := range next.cmds {
		if cmd.Name == args[i] || cmd.Wildcard {
			found = cmd
			offset = len(cmd.Pos)
			// see if we exceed the given arguments.
			if offset > n-i {
				return nil, 0, fmt.Errorf("%s needs %d position arguments", cmd.Name, offset)
			}
			// we got no position argument to skip if
			// only with the help flag.
			if !help {
				i += offset
			}
			i++
			next = found
			// but this is damn fast!
			goto search
		}
	}
	return
}

func (c *Cmd) validate() error {
	for _, flag := range c.Flags {
		if flag.Invalid() {
			return fmt.Errorf("option `%s` is required", flag.Key())
		}
	}
	return nil
}

func (c *Cmd) help(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 8, 1, '\t', tabwriter.AlignRight)

	if len(c.About) != 0 {
		fmt.Fprintf(tw, "About:\n    %s %s\n\n", c.About, c.Version)
	}

	sep := ""
	if len(c.parent) > 0 {
		sep = " "
	}
	fmt.Fprintf(tw, "Usage:\n    %s%s%s", c.parent, sep, c.Name)

	if len(c.cmds) > 0 {
		fmt.Fprint(tw, " <command>")
	}

	for _, arg := range c.Pos {
		fmt.Fprintf(tw, " <%s>", arg)
	}

	fmt.Fprint(tw, " [options...]\n")

	if len(c.Example) > 0 {
		fmt.Fprintf(tw, "\nExamples:\n    %s\n", c.Example)
	}

	if c.cmds != nil {
		fmt.Fprint(tw, "\nCommands:\n")
		for _, cmd := range c.cmds {
			fmt.Fprintf(tw, "    %s\t%s\n", cmd.Name, cmd.About)
		}
	}

	if c.Flags != nil {
		fmt.Fprint(tw, "\nOptions:\n")
		for _, flag := range c.Flags {
			fmt.Fprintf(tw, "    --%s\t%s\t(default %v) \n", flag.Key(), flag.Help(), flag.Var())
		}
	}

	return tw.Flush()
}

func (c *Cmd) init(parent string, root *Cmd, globalFlags []Flag) {
	set := flag.NewFlagSet(c.Name, flag.ExitOnError)

	if c.Flags != nil {
		sort.Slice(c.Flags, func(i, j int) bool {
			return c.Flags[i].Key() < c.Flags[j].Key()
		})

		for _, flag := range c.Flags {
			switch f := flag.(type) {
			case *BoolFlag:
				set.BoolVar(&f.Value, f.Name, f.Value, f.Usage)
			case *IntFlag:
				set.IntVar(&f.Value, f.Name, f.Value, f.Usage)
			case *UintFlag:
				set.UintVar(&f.Value, f.Name, f.Value, f.Usage)
			case *FloatFlag:
				set.Float64Var(&f.Value, f.Name, f.Value, f.Usage)
			case *StringFlag:
				set.StringVar(&f.Value, f.Name, f.Value, f.Usage)
			case *StringSliceFlag:
				set.Var((*StringSlice)(unsafe.Pointer(&f.Value)), f.Name, f.Usage)
			case *IntSliceFlag:
				set.Var((*IntSlice)(unsafe.Pointer(&f.Value)), f.Name, f.Usage)
			case *UintSliceFlag:
				set.Var((*UintSlice)(unsafe.Pointer(&f.Value)), f.Name, f.Usage)
			case *FloatSliceFlag:
				set.Var((*FloatSlice)(unsafe.Pointer(&f.Value)), f.Name, f.Usage)
			case *StringMapFlag:
				f.Value = StringMap{}
				set.Var((*StringMap)(unsafe.Pointer(&f.Value)), f.Name, f.Usage)
			case *MapFlag:
				f.Value = Map{}
				set.Var((*Map)(unsafe.Pointer(&f.Value)), f.Name, f.Usage)
			default:
				fmt.Printf("unsupported flag %s, ignored\n", f)
			}
		}

	}

	set.Usage = func() {
		c.help(os.Stdout)
	}
	c.set = set
	c.parent = parent

	for _, cmd := range c.cmds {
		cmd.Flags = append(cmd.Flags, globalFlags...)
		cmd.init(fmt.Sprintf("%s%s", parent, c.Name), root, globalFlags)
	}
}

func seperate(args []string) (int, bool) {
	var help bool
	for i, arg := range args {
		if arg[0] == '-' {
			switch arg[1] {
			case '-':
				help = arg[2:] == "help"
			case 'h':
				help = true
			}
			return i, help
		}
	}
	return len(args), help
}

package model

import (
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type Config struct {
	DSN    string
	Driver string
	Folder string

	Verbose bool
	Quiet   bool
}

func NewConfig() *Config {
	return &Config{}
}

// NewFlagSet is used for command flags.
func NewFlagSet(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.SetInterspersed(true)
	return fs
}

func (c *Config) ParseFlags() ([]string, error) {
	flagSet := NewFlagSet("Config")
	flagSet.StringVar(&c.DSN, "db-dsn", os.Getenv("DB_DSN"), "Database DSN")
	flagSet.StringVar(&c.Driver, "db-driver", os.Getenv("DB_DRIVER"), "Database Driver")
	flagSet.StringVarP(&c.Folder, "folder", "f", "output", "Folder with outputs")
	flagSet.BoolVarP(&c.Verbose, "verbose", "v", false, "Folder with outputs")
	flagSet.BoolVarP(&c.Quiet, "quiet", "q", false, "Quiet output")

	k, u := filterKnownArgs(flagSet, os.Args[1:])

	err := flagSet.Parse(k)
	if err != nil {
		return nil, err
	}

	result := flagSet.Args()

	return append(result, u...), nil
}

func (c *Config) GetDSN() string {
	if c.DSN == "" {
		return ":memory:"
	}
	return c.DSN
}

func (c *Config) GetDriver() string {
	if c.Driver == "" {
		return "sqlite"
	}
	return c.Driver
}

// filterKnownArgs separates known flags from unknown ones
func filterKnownArgs(flagSet *pflag.FlagSet, args []string) (knownArgs, unknownArgs []string) {
	isKnownFlag := func(f string) bool {
		if flagSet.Lookup(f) != nil {
			return true
		}
		if len(f) == 1 && flagSet.ShorthandLookup(f) != nil {
			return true
		}
		return false
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if isFlag := strings.HasPrefix(arg, "-"); !isFlag {
			knownArgs = append(knownArgs, arg)
			continue
		}

		flagName := strings.TrimPrefix(arg, "-")
		if strings.Contains(flagName, "=") {
			flagName = strings.SplitN(flagName, "=", 2)[0]
		}

		if isKnownFlag(flagName) {
			knownArgs = append(knownArgs, arg)
			if !strings.Contains(arg, "=") && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				knownArgs = append(knownArgs, args[i+1])
				i++
			}
			continue
		}

		unknownArgs = append(unknownArgs, arg)
		if !strings.Contains(arg, "=") && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
			unknownArgs = append(unknownArgs, args[i+1])
			i++
		}
	}
	return
}

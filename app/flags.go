package app

import "flag"

type runTimeFlags struct {
	configYamlAddress string
}

func NewRunTimeFlags() *runTimeFlags {
	flags := new(runTimeFlags)

	configAddress := flag.String("config", "./config.yaml", "config file address")
	flag.Parse()

	flags.configYamlAddress = *configAddress

	return flags
}

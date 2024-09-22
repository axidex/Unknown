package shell

import (
	"github.com/axidex/Unknown/config"
)

type CLIBuilder interface {
	StartScanCommand(pathToSources string, reportFile string) []string
}

func CreateCLIBuilder(config config.ShellCommand) CLIBuilder {
	return &CLIBuilderImpl{
		config: config,
	}
}

type CLIBuilderImpl struct {
	config config.ShellCommand
}

func (builder *CLIBuilderImpl) StartScanCommand(pathToSources string, reportFile string) []string {

	command := []string{
		builder.config.Binary,
		"detect",
		"--source",
		pathToSources,
		"-r",
		reportFile,
	}
	for _, arg := range builder.config.AdditionalCommands {
		command = append(command, arg)
	}

	return command
}

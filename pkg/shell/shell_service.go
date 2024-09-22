package shell

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/axidex/Unknown/pkg/logger"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Service interface {
	Scan(ctx context.Context, scan Scan) (string, error)
}

type ServiceCLI struct {
	logger  logger.Logger
	builder CLIBuilder
}

func CreateServiceCLI(
	logger logger.Logger,
	builder CLIBuilder,
) Service {
	return &ServiceCLI{
		logger:  logger,
		builder: builder,
	}
}

func (service *ServiceCLI) Scan(ctx context.Context, scan Scan) (string, error) {
	service.logger.Infof("Got scan")

	pathToReport := path.Join(scan.ScanFolder, fmt.Sprintf("%s_report", scan.TaskId))

	startScanCommand := service.builder.StartScanCommand(scan.ScanFolder, pathToReport)
	service.logger.Infof("Running command: %s", strings.Join(startScanCommand, " "))
	err := service.launchCommand(ctx, startScanCommand)
	if err != nil {
		return "", nil
	}

	return pathToReport, nil
}

func (service *ServiceCLI) launchCommand(ctx context.Context, command []string) error {
	service.logger.Infof("Launching command - %s", strings.Join(command, " "))

	cmd := exec.CommandContext(ctx, command[0], command[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	go service.processOutput(bufio.NewScanner(stdout), service.logger.Infof)
	go service.processOutput(bufio.NewScanner(stderr), service.logger.Infof)

	if err = cmd.Wait(); err != nil {
		return err
	}

	code := cmd.ProcessState.ExitCode()
	service.logger.Infof("Command return code - %d", code)
	if code != 0 {
		return errors.New(fmt.Sprintf("return code not zero - %d", code))
	}

	return nil
}

// Aux functions
func (service *ServiceCLI) removeFolder(workdir string) {
	service.logger.Infof("Removing tmp folder: %s", workdir)
	err := os.RemoveAll(workdir)
	if err != nil {
		service.logger.Errorf("Error when deleting tmp folder %s | %s", workdir, err)
	}
}

func (service *ServiceCLI) processOutput(scanner *bufio.Scanner, logFunc func(string, ...interface{})) {
	for scanner.Scan() {
		logFunc(scanner.Text())
	}
}

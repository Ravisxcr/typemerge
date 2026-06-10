package typst

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Compiler interface {
	Compile(context.Context, string, string) error
}

type CommandCompiler struct {
	Binary string
}

func (c CommandCompiler) Compile(ctx context.Context, sourcePath, destinationPath string) error {
	command := exec.CommandContext(ctx, c.Binary, "compile", sourcePath, destinationPath)
	output, err := command.CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(output))
		if message == "" {
			return fmt.Errorf("compile %s: %w", sourcePath, err)
		}
		return fmt.Errorf("compile %s: %w: %s", sourcePath, err, message)
	}
	return nil
}

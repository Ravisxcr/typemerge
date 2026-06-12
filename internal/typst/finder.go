package typst

import (
	"fmt"
	"os/exec"
)

func Find(binary string) (string, error) {
	path, err := exec.LookPath(binary)
	if err != nil {
		return "", fmt.Errorf("find Typst executable %q: %w", binary, err)
	}
	return path, nil
}

package render

import "fmt"

func ValidateSource(source []byte) error {
	if len(source) == 0 {
		return fmt.Errorf("rendered document is empty")
	}
	return nil
}

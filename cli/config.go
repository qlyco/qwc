package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"slices"

	"github.com/qlyco/qwc/internal/res"
)

func applyConfig(configPath string) error {
	source, err := os.ReadFile(configPath)

	if err != nil {
		return err
	}

	settings := map[string]json.RawMessage{}
	err = json.Unmarshal(source, &settings)

	if err != nil {
		return fmt.Errorf(res.FileErrorFormat, configPath, err)
	}

	typed := map[string]bool{}

	flag.Visit(func(typedFlag *flag.Flag) {
		typed[typedFlag.Name] = true
	})

	for name, value := range settings {
		if flag.Lookup(name) == nil || slices.Contains(res.CommandFlags, name) {
			return fmt.Errorf(res.UnknownSettingMessage, name, configPath)
		}

		if typed[name] {
			continue
		}

		err = flag.Set(name, settingText(value))

		if err != nil {
			return fmt.Errorf(res.InvalidSettingMessage, value, name, configPath)
		}
	}

	return nil
}

func settingText(value json.RawMessage) string {
	var text string

	err := json.Unmarshal(value, &text)

	if err != nil {
		return string(value)
	}

	return text
}

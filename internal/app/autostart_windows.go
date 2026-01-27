//go:build windows

package app

import (
	"golang.org/x/sys/windows/registry"
)

const runKeyPath = `Software\Microsoft\Windows\CurrentVersion\Run`
const runValueName = "gammatray"

func SetAutoStart(enabled bool, exePath string) error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, runKeyPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	if enabled {
		return k.SetStringValue(runValueName, exePath)
	}
	_ = k.DeleteValue(runValueName)
	return nil
}

func GetAutoStart() (bool, string, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKeyPath, registry.QUERY_VALUE)
	if err != nil {
		return false, "", nil
	}
	defer k.Close()

	v, _, err := k.GetStringValue(runValueName)
	if err != nil {
		return false, "", nil
	}
	return true, v, nil
}

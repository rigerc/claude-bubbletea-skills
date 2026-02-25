package screens

import "scaffold/config"

// BackMsg signals that the current screen wants to go back.
type BackMsg struct{}

// SettingsSavedMsg carries the updated config after the user submits the form.
type SettingsSavedMsg struct {
	Cfg config.Config
}

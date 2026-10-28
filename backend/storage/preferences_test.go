package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConnectionSortPreferenceDefaultsForExistingConfig(t *testing.T) {
	p := &PreferencesStorage{storage: &localStorage{ConfPath: filepath.Join(t.TempDir(), "preferences.yaml")}}
	if err := os.WriteFile(p.storage.ConfPath, []byte("behavior:\n  aside_width: 300\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if p.GetPreferences().Behavior.ConnectionSortOrder != "original" {
		t.Fatal("existing configuration must retain its original connection order")
	}
}

func TestConnectionSortPreferencePersistsAcrossReload(t *testing.T) {
	for _, order := range []string{"original", "asc", "desc"} {
		t.Run(order, func(t *testing.T) {
			filename := filepath.Join(t.TempDir(), "preferences.yaml")
			p := &PreferencesStorage{storage: &localStorage{ConfPath: filename}}
			preferences := p.DefaultPreferences()
			preferences.Behavior.WindowWidth = 1400
			preferences.Behavior.WindowHeight = 900
			preferences.Behavior.WindowMaximised = true
			if err := p.SetPreferences(&preferences); err != nil {
				t.Fatal(err)
			}
			if err := p.UpdatePreferences(map[string]any{"behavior.connectionSortOrder": order}); err != nil {
				t.Fatal(err)
			}
			reloaded := &PreferencesStorage{storage: &localStorage{ConfPath: filename}}
			behavior := reloaded.GetPreferences().Behavior
			if behavior.ConnectionSortOrder != order {
				t.Fatal("connection sort mode was not restored after reloading preferences")
			}
			if behavior.WindowWidth != 1400 || behavior.WindowHeight != 900 || !behavior.WindowMaximised {
				t.Fatal("changing the sort mode overwrote the saved window state")
			}
		})
	}
}

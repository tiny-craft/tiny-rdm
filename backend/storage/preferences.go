package storage

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
	"sync"
)

type PreferencesStorage struct {
	storage *localStorage
	mutex   sync.Mutex
}

func NewPreferences() *PreferencesStorage {
	return &PreferencesStorage{
		storage: NewLocalStore("preferences.yaml"),
	}
}

func (p *PreferencesStorage) DefaultPreferences() map[string]any {
	return map[string]any{
		"general": map[string]any{
			"language":           "auto",
			"font":               "",
			"font_size":          14,
			"use_sys_proxy":      false,
			"use_sys_proxy_http": false,
			"check_update":       true,
			"skip_version":       "",
			"aside_width":        300,
		},
		"editor": map[string]any{
			"font":      "",
			"font_size": 14,
		},
	}
}

func (p *PreferencesStorage) getPreferences() (ret map[string]any) {
	b, err := p.storage.Load()
	if err != nil {
		ret = p.DefaultPreferences()
		return
	}

	if err = yaml.Unmarshal(b, &ret); err != nil {
		ret = p.DefaultPreferences()
		return
	}
	return
}

func (p *PreferencesStorage) flatPreferences(data map[string]any, prefix string) map[string]any {
	flattened := make(map[string]any)
	for key, value := range data {
		newKey := key
		if prefix != "" {
			newKey = prefix + "." + key
		}

		if nested, ok := value.(map[string]any); ok {
			nestedFlattened := p.flatPreferences(nested, newKey)
			for k, v := range nestedFlattened {
				flattened[k] = v
			}
		} else {
			flattened[newKey] = value
		}
	}
	return flattened
}

// GetPreferences Get preferences from local
func (p *PreferencesStorage) GetPreferences() (ret map[string]any) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	pref := p.getPreferences()
	ret = p.flatPreferences(pref, "")
	return
}

func (p *PreferencesStorage) Value(keys ...string) any {
	kv := p.getPreferences()
	var ok bool
	var a any
	total := len(keys)
	for i, key := range keys {
		if a, ok = kv[key]; !ok {
			return nil
		}
		if i == total-1 {
			// last key, return value
			return a
		}
		if kv, ok = a.(map[string]any); !ok {
			return nil
		}
	}
	return nil
}

func (p *PreferencesStorage) setPreferences(pf map[string]any, key string, value any) error {
	keyPath := strings.Split(key, ".")
	if len(keyPath) <= 0 {
		return fmt.Errorf("invalid key path(%s)", key)
	}
	var node any = pf
	for _, k := range keyPath[:len(keyPath)-1] {
		if subNode, ok := node.(map[string]any); ok {
			node = subNode[k]
		} else {
			return fmt.Errorf("invalid key path(%s)", key)
		}
	}

	if subNode, ok := node.(map[string]any); ok {
		subNode[keyPath[len(keyPath)-1]] = value
	}

	return nil
}

func (p *PreferencesStorage) savePreferences(pf map[string]any) error {
	b, err := yaml.Marshal(&pf)
	if err != nil {
		return err
	}

	if err = p.storage.Store(b); err != nil {
		return err
	}
	return nil
}

// SetPreferences assign value to key path, the key path use "." to indicate multiple level
func (p *PreferencesStorage) SetPreferences(key string, value any) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	pf := p.getPreferences()
	if err := p.setPreferences(pf, key, value); err != nil {
		return err
	}
	return p.savePreferences(pf)
}

// SetPreferencesN set multiple key path and value
func (p *PreferencesStorage) SetPreferencesN(values map[string]any) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	pf := p.getPreferences()
	for path, v := range values {
		if err := p.setPreferences(pf, path, v); err != nil {
			return err
		}
	}

	return p.savePreferences(pf)
}

func (p *PreferencesStorage) RestoreDefault() map[string]any {
	pf := p.DefaultPreferences()
	p.savePreferences(pf)
	return p.flatPreferences(pf, "")
}

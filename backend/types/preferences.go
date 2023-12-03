package types

import "tinyrdm/backend/consts"

type Preferences struct {
	Behavior PreferencesBehavior `json:"behavior" yaml:"behavior"`
	General  PreferencesGeneral  `json:"general" yaml:"general"`
	Editor   PreferencesEditor   `json:"editor" yaml:"editor"`
}

func NewPreferences() Preferences {
	return Preferences{
		Behavior: PreferencesBehavior{
			AsideWidth:   consts.DEFAULT_ASIDE_WIDTH,
			WindowWidth:  consts.DEFAULT_WINDOW_WIDTH,
			WindowHeight: consts.DEFAULT_WINDOW_HEIGHT,
		},
		General: PreferencesGeneral{
			Theme:       "auto",
			Language:    "auto",
			FontSize:    consts.DEFAULT_FONT_SIZE,
			ScanSize:    consts.DEFAULT_SCAN_SIZE,
			CheckUpdate: true,
		},
		Editor: PreferencesEditor{
			FontSize:    consts.DEFAULT_FONT_SIZE,
			ShowLineNum: true,
		},
	}
}

type PreferencesBehavior struct {
	AsideWidth      int  `json:"asideWidth" yaml:"aside_width"`
	WindowWidth     int  `json:"windowWidth" yaml:"window_width"`
	WindowHeight    int  `json:"windowHeight" yaml:"window_height"`
	WindowMaximised bool `json:"windowMaximised" yaml:"window_maximised"`
}

type PreferencesGeneral struct {
	Theme           string `json:"theme" yaml:"theme"`
	Language        string `json:"language" yaml:"language"`
	Font            string `json:"font" yaml:"font,omitempty"`
	FontSize        int    `json:"fontSize" yaml:"font_size"`
	ScanSize        int    `json:"scanSize" yaml:"scan_size"`
	UseSysProxy     bool   `json:"useSysProxy" yaml:"use_sys_proxy,omitempty"`
	UseSysProxyHttp bool   `json:"useSysProxyHttp" yaml:"use_sys_proxy_http,omitempty"`
	CheckUpdate     bool   `json:"checkUpdate" yaml:"check_update"`
	SkipVersion     string `json:"skipVersion" yaml:"skip_version,omitempty"`
}

type PreferencesEditor struct {
	Font        string `json:"font" yaml:"font,omitempty"`
	FontSize    int    `json:"fontSize" yaml:"font_size"`
	ShowLineNum bool   `json:"showLineNum" yaml:"show_line_num,omitempty"`
}

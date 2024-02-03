package types

import "tinyrdm/backend/consts"

type Preferences struct {
	Behavior PreferencesBehavior  `json:"behavior" yaml:"behavior"`
	General  PreferencesGeneral   `json:"general" yaml:"general"`
	Editor   PreferencesEditor    `json:"editor" yaml:"editor"`
	Cli      PreferencesCli       `json:"cli" yaml:"cli"`
	Decoder  []PreferencesDecoder `json:"decoder" yaml:"decoder,omitempty"`
}

func NewPreferences() Preferences {
	return Preferences{
		Behavior: PreferencesBehavior{
			AsideWidth:   consts.DEFAULT_ASIDE_WIDTH,
			WindowWidth:  consts.DEFAULT_WINDOW_WIDTH,
			WindowHeight: consts.DEFAULT_WINDOW_HEIGHT,
		},
		General: PreferencesGeneral{
			Theme:        "auto",
			Language:     "auto",
			FontSize:     consts.DEFAULT_FONT_SIZE,
			ScanSize:     consts.DEFAULT_SCAN_SIZE,
			KeyIconStyle: 0,
			CheckUpdate:  true,
		},
		Editor: PreferencesEditor{
			FontSize:    consts.DEFAULT_FONT_SIZE,
			ShowLineNum: true,
			ShowFolding: true,
		},
		Cli: PreferencesCli{
			FontSize:    consts.DEFAULT_FONT_SIZE,
			CursorStyle: "block",
		},
		Decoder: []PreferencesDecoder{},
	}
}

type PreferencesBehavior struct {
	AsideWidth      int  `json:"asideWidth" yaml:"aside_width"`
	WindowWidth     int  `json:"windowWidth" yaml:"window_width"`
	WindowHeight    int  `json:"windowHeight" yaml:"window_height"`
	WindowMaximised bool `json:"windowMaximised" yaml:"window_maximised"`
	WindowPosX      int  `json:"windowPosX" yaml:"window_pos_x"`
	WindowPosY      int  `json:"windowPosY" yaml:"window_pos_y"`
}

type PreferencesGeneral struct {
	Theme           string   `json:"theme" yaml:"theme"`
	Language        string   `json:"language" yaml:"language"`
	Font            string   `json:"font" yaml:"font,omitempty"`
	FontFamily      []string `json:"fontFamily" yaml:"font_family,omitempty"`
	FontSize        int      `json:"fontSize" yaml:"font_size"`
	ScanSize        int      `json:"scanSize" yaml:"scan_size"`
	KeyIconStyle    int      `json:"keyIconStyle" yaml:"key_icon_style"`
	UseSysProxy     bool     `json:"useSysProxy" yaml:"use_sys_proxy,omitempty"`
	UseSysProxyHttp bool     `json:"useSysProxyHttp" yaml:"use_sys_proxy_http,omitempty"`
	CheckUpdate     bool     `json:"checkUpdate" yaml:"check_update"`
	SkipVersion     string   `json:"skipVersion" yaml:"skip_version,omitempty"`
}

type PreferencesEditor struct {
	Font        string   `json:"font" yaml:"font,omitempty"`
	FontFamily  []string `json:"fontFamily" yaml:"font_family,omitempty"`
	FontSize    int      `json:"fontSize" yaml:"font_size"`
	ShowLineNum bool     `json:"showLineNum" yaml:"show_line_num"`
	ShowFolding bool     `json:"showFolding" yaml:"show_folding"`
}

type PreferencesCli struct {
	FontFamily  []string `json:"fontFamily" yaml:"font_family,omitempty"`
	FontSize    int      `json:"fontSize" yaml:"font_size"`
	CursorStyle string   `json:"cursorStyle" yaml:"cursor_style,omitempty"`
}

type PreferencesDecoder struct {
	Name       string   `json:"name" yaml:"name"`
	Enable     bool     `json:"enable" yaml:"enable"`
	Auto       bool     `json:"auto" yaml:"auto"`
	DecodePath string   `json:"decodePath" yaml:"decode_path"`
	DecodeArgs []string `json:"decodeArgs" yaml:"decode_args,omitempty"`
	EncodePath string   `json:"encodePath" yaml:"encode_path"`
	EncodeArgs []string `json:"encodeArgs" yaml:"encode_args,omitempty"`
}

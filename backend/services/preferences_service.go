package services

import (
	"encoding/json"
	"github.com/adrg/sysfont"
	"net/http"
	"sort"
	"strings"
	"sync"
	"tinyrdm/backend/consts"
	storage2 "tinyrdm/backend/storage"
	"tinyrdm/backend/types"
	"tinyrdm/backend/utils/coll"
)

type preferencesService struct {
	pref          *storage2.PreferencesStorage
	clientVersion string
}

var preferences *preferencesService
var oncePreferences sync.Once

func Preferences() *preferencesService {
	if preferences == nil {
		oncePreferences.Do(func() {
			preferences = &preferencesService{
				pref:          storage2.NewPreferences(),
				clientVersion: "",
			}
		})
	}
	return preferences
}

func (p *preferencesService) GetPreferences() (resp types.JSResp) {
	resp.Data = p.pref.GetPreferences()
	resp.Success = true
	return
}

func (p *preferencesService) SetPreferences(pf types.Preferences) (resp types.JSResp) {
	err := p.pref.SetPreferences(&pf)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	return
}

func (p *preferencesService) UpdatePreferences(value map[string]any) (resp types.JSResp) {
	err := p.pref.UpdatePreferences(value)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Success = true
	return
}

func (p *preferencesService) RestorePreferences() (resp types.JSResp) {
	defaultPref := p.pref.RestoreDefault()
	resp.Data = map[string]any{
		"pref": defaultPref,
	}
	resp.Success = true
	return
}

type FontItem struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (p *preferencesService) GetFontList() (resp types.JSResp) {
	finder := sysfont.NewFinder(nil)
	fontSet := coll.NewSet[string]()
	var fontList []FontItem
	for _, font := range finder.List() {
		if len(font.Family) > 0 && !strings.HasPrefix(font.Family, ".") && fontSet.Add(font.Family) {
			fontList = append(fontList, FontItem{
				Name: font.Family,
				Path: font.Filename,
			})
		}
	}
	sort.Slice(fontList, func(i, j int) bool {
		return fontList[i].Name < fontList[j].Name
	})
	resp.Data = map[string]any{
		"fonts": fontList,
	}
	resp.Success = true
	return
}

func (p *preferencesService) SetAppVersion(ver string) {
	if !strings.HasPrefix(ver, "v") {
		p.clientVersion = "v" + ver
	} else {
		p.clientVersion = ver
	}
}

func (p *preferencesService) GetAppVersion() (resp types.JSResp) {
	resp.Success = true
	resp.Data = map[string]any{
		"version": p.clientVersion,
	}
	return
}

func (p *preferencesService) SaveWindowSize(width, height int) {
	if width >= consts.DEFAULT_WINDOW_WIDTH && height >= consts.DEFAULT_WINDOW_HEIGHT {
		p.UpdatePreferences(map[string]any{
			"behavior.windowWidth":  width,
			"behavior.windowHeight": height,
		})
	}
}

func (p *preferencesService) GetWindowSize() (width, height int) {
	data := p.pref.GetPreferences()
	width, height = data.Behavior.WindowWidth, data.Behavior.WindowHeight
	if width <= 0 {
		width = consts.DEFAULT_WINDOW_WIDTH
	}
	if height <= 0 {
		height = consts.DEFAULT_WINDOW_HEIGHT
	}
	return
}

type latestRelease struct {
	Name    string `json:"name"`
	TagName string `json:"tag_name"`
	Url     string `json:"url"`
	HtmlUrl string `json:"html_url"`
}

func (p *preferencesService) CheckForUpdate() (resp types.JSResp) {
	// request latest version
	res, err := http.Get("https://api.github.com/repos/tiny-craft/tiny-rdm/releases/latest")
	if err != nil || res.StatusCode != http.StatusOK {
		resp.Msg = "network error"
		return
	}

	var respObj latestRelease
	err = json.NewDecoder(res.Body).Decode(&respObj)
	if err != nil {
		resp.Msg = "invalid content"
		return
	}

	// compare with current version
	resp.Success = true
	resp.Data = map[string]any{
		"version":  p.clientVersion,
		"latest":   respObj.TagName,
		"page_url": respObj.HtmlUrl,
	}
	return
}

package ui

import (
    "fmt"
    "go-musicfox/ds"
)

type SearchResultMenu struct {
    menus      []MenuItem
    offset     int
    limit      int
    searchType SearchType
    keyword    string
    result     interface{}
}

func NewSearchResultMenu(searchType SearchType) *SearchResultMenu {
    return &SearchResultMenu{
        offset:     0,
        limit:      100,
        searchType: searchType,
    }
}

func (m *SearchResultMenu) MenuData() interface{} {
    return m.result
}

func (m *SearchResultMenu) BeforeBackMenuHook() Hook {
    return func(model *NeteaseModel) bool {
        if model.searchModel.wordsInput.Value() != "" {
            model.searchModel.wordsInput.SetValue("")
        }

        return true
    }
}

func (m *SearchResultMenu) IsPlayable() bool {
    playableMap := map[SearchType]bool{
        StSingleSong: true,
        StAlbum:      false,
        StSinger:     false,
        StPlaylist:   false,
        StUser:       false,
        StLyric:      true,
        StRadio:      false,
    }

    if playable, ok := playableMap[m.searchType]; ok {
        return playable
    }

    return false
}

func (m *SearchResultMenu) ResetPlaylistWhenPlay() bool {
    return false
}

func (m *SearchResultMenu) GetMenuKey() string {
    return fmt.Sprintf("search_result_%d_%s", m.searchType, m.keyword)
}

func (m *SearchResultMenu) MenuViews() []MenuItem {
    return m.menus
}

func (m *SearchResultMenu) SubMenu(_ *NeteaseModel, index int) IMenu {
    switch resultWithType := m.result.(type) {
    case []ds.Song:
        return nil
    case []ds.Album:
        if index >= len(resultWithType) {
            return nil
        }
        return NewAlbumDetailMenu(resultWithType[index].Id)
    case []ds.Playlist:
        if index >= len(resultWithType) {
            return nil
        }
        return NewPlaylistDetailMenu(resultWithType[index].Id)
    case []ds.Artist:
        if index >= len(resultWithType) {
            return nil
        }
        return NewArtistDetailMenu(resultWithType[index].Id)
    case []ds.User:
        if index >= len(resultWithType) {
            return nil
        }
        return NewUserPlaylistMenu(resultWithType[index].UserId)
    case []ds.DjRadio:
       if index >= len(resultWithType) {
           return nil
       }
       return NewDjRadioDetailMenu(resultWithType[index].Id)
    }

    return nil
}

func (m *SearchResultMenu) BeforePrePageHook() Hook {
    // Nothing to do
    return nil
}

func (m *SearchResultMenu) BeforeNextPageHook() Hook {
    // Nothing to do
    return nil
}

func (m *SearchResultMenu) BeforeEnterMenuHook() Hook {
    return func(model *NeteaseModel) bool {
        if model.searchModel.wordsInput.Value() == "" {
            // 显示搜索页面
            SearchHandle(model, m.searchType)
            return false
        }

        m.result = model.searchModel.result

        switch resultWithType := m.result.(type) {
        case []ds.Song:
            m.menus = GetViewFromSongs(resultWithType)
        case []ds.Album:
            m.menus = GetViewFromAlbums(resultWithType)
        case []ds.Playlist:
            m.menus = GetViewFromPlaylists(resultWithType)
        case []ds.Artist:
            m.menus = GetViewFromArtists(resultWithType)
        case []ds.User:
            m.menus = GetViewFromUsers(resultWithType)
        case []ds.DjRadio:
            m.menus = GetViewFromDjRadios(resultWithType)
        }

        return true
    }
}

func (m *SearchResultMenu) BottomOutHook() Hook {
    // 暂时只取前100
    return nil
}

func (m *SearchResultMenu) TopOutHook() Hook {
    // Nothing to do
    return nil
}

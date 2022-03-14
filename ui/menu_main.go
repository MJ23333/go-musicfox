package ui

type MainMenu struct {
	menus    []MenuItem
	menuList []IMenu
}

func NewMainMenu() *MainMenu {
	mainMenu := new(MainMenu)
	mainMenu.menus = []MenuItem{
		{Title: "日推"},
		{Title: "歌单推荐"},
		{Title: "我的歌单"},
		{Title: "私人FM"},
		{Title: "专辑列表"},
		{Title: "搜索"},
		{Title: "排行榜"},
		{Title: "精选歌单"},
		{Title: "热门歌手"},
		{Title: "云盘"},
		{Title: "主播电台"},
		{Title: "帮助"},
		{Title: "检查更新"},
	}
	mainMenu.menuList = []IMenu{
		NewDailyRecommendSongsMenu(),
		NewDailyRecommendPlaylistMenu(),
		NewUserPlaylistMenu(CurUser),
		NewPersonalFmMenu(),
		NewAlbumListMenu(),
		NewSearchTypeMenu(),
		NewRanksMenu(),
		NewHighQualityPlaylistsMenu(),
		NewHotArtistsMenu(),
		NewCloudMenu(),
		NewRadioDjTypeMenu(),
		NewHelpMenu(),
		NewCheckUpdateMenu(),
	}

	return mainMenu
}

func (m *MainMenu) MenuData() interface{} {
	return nil
}

func (m *MainMenu) IsPlayable() bool {
	return false
}

func (m *MainMenu) ResetPlaylistWhenPlay() bool {
	return false
}

func (m *MainMenu) GetMenuKey() string {
	return "main_menu"
}

func (m *MainMenu) MenuViews() []MenuItem {
	return m.menus
}

func (m *MainMenu) SubMenu(_ *NeteaseModel, index int) IMenu {

	if index >= len(m.menuList) {
		return nil
	}

	return m.menuList[index]
}

func (m *MainMenu) BeforePrePageHook() Hook {
	// Nothing to do
	return nil
}

func (m *MainMenu) BeforeNextPageHook() Hook {
	// Nothing to do
	return nil
}

func (m *MainMenu) BeforeEnterMenuHook() Hook {
	// Nothing to do
	return nil
}

func (m *MainMenu) BeforeBackMenuHook() Hook {
	// Nothing to do
	return nil
}

func (m *MainMenu) BottomOutHook() Hook {
	// Nothing to do
	return nil
}

func (m *MainMenu) TopOutHook() Hook {
	// Nothing to do
	return nil
}

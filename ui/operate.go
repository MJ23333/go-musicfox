package ui

import (
	"fmt"
	"go-musicfox/constants"
	"go-musicfox/db"
	"go-musicfox/ds"
	"go-musicfox/utils"
	"math"
	"strconv"

	"github.com/anhoder/netease-music/service"
	"github.com/muesli/termenv"
)

type menuStackItem struct {
	menuList      []MenuItem
	selectedIndex int
	menuCurPage   int
	menuTitle     string
	menu          IMenu
}

// 上移
func moveUp(m *NeteaseModel) {
	topHook := m.menu.TopOutHook()
	if m.doubleColumn {
		if m.selectedIndex-2 < 0 && topHook != nil {
			loading := NewLoading(m)
			loading.start()
			if res := topHook(m); !res {
				loading.complete()
				return
			}
			// 更新菜单UI
			m.menuList = m.menu.MenuViews()
			loading.complete()
		}
		if m.selectedIndex-2 < 0 {
			return
		}
		m.selectedIndex -= 2
	} else {
		if m.selectedIndex-1 < 0 && topHook != nil {
			loading := NewLoading(m)
			loading.start()
			if res := topHook(m); !res {
				loading.complete()
				return
			}
			m.menuList = m.menu.MenuViews()
			loading.complete()
		}
		if m.selectedIndex-1 < 0 {
			return
		}
		m.selectedIndex--
	}
	if m.selectedIndex < (m.menuCurPage-1)*m.menuPageSize {
		prePage(m)
	}
}

// 下移
func moveDown(m *NeteaseModel) {
	bottomHook := m.menu.BottomOutHook()
	if m.doubleColumn {
		if m.selectedIndex+2 > len(m.menuList)-1 && bottomHook != nil {
			loading := NewLoading(m)
			loading.start()
			if res := bottomHook(m); !res {
				loading.complete()
				return
			}
			m.menuList = m.menu.MenuViews()
			loading.complete()
		}
		if m.selectedIndex+2 > len(m.menuList)-1 {
			return
		}
		m.selectedIndex += 2
	} else {
		if m.selectedIndex+1 > len(m.menuList)-1 && bottomHook != nil {
			loading := NewLoading(m)
			loading.start()
			if res := bottomHook(m); !res {
				loading.complete()
				return
			}
			m.menuList = m.menu.MenuViews()
			loading.complete()
		}
		if m.selectedIndex+1 > len(m.menuList)-1 {
			return
		}
		m.selectedIndex++
	}
	if m.selectedIndex >= m.menuCurPage*m.menuPageSize {
		nextPage(m)
	}
}

// 左移
func moveLeft(m *NeteaseModel) {
	if !m.doubleColumn || m.selectedIndex%2 == 0 || m.selectedIndex-1 < 0 {
		return
	}
	m.selectedIndex--
}

// 右移
func moveRight(m *NeteaseModel) {
	if !m.doubleColumn || m.selectedIndex%2 != 0 {
		return
	}
	if bottomHook := m.menu.BottomOutHook(); m.selectedIndex >= len(m.menuList)-1 && bottomHook != nil {
		loading := NewLoading(m)
		loading.start()
		if res := bottomHook(m); !res {
			loading.complete()
			return
		}
		m.menuList = m.menu.MenuViews()
		loading.complete()
	}
	if m.selectedIndex >= len(m.menuList)-1 {
		return
	}
	m.selectedIndex++
}

// 切换到上一页
func prePage(m *NeteaseModel) {
	m.isListeningKey = false
	defer func() {
		m.isListeningKey = true
	}()

	if prePageHook := m.menu.BeforePrePageHook(); prePageHook != nil {
		loading := NewLoading(m)
		loading.start()
		if res := prePageHook(m); !res {
			loading.complete()
			return
		}
		loading.complete()
	}

	if m.menuCurPage <= 1 {
		return
	}
	m.menuCurPage--
}

// 切换到下一页
func nextPage(m *NeteaseModel) {
	m.isListeningKey = false
	defer func() {
		m.isListeningKey = true
	}()

	if nextPageHook := m.menu.BeforeNextPageHook(); nextPageHook != nil {
		loading := NewLoading(m)
		loading.start()
		if res := nextPageHook(m); !res {
			loading.complete()
			return
		}
		loading.complete()
	}
	if m.menuCurPage >= int(math.Ceil(float64(len(m.menuList))/float64(m.menuPageSize))) {
		return
	}

	m.menuCurPage++
}

// 进入菜单
func enterMenu(m *NeteaseModel) {
	m.isListeningKey = false
	defer func() {
		m.isListeningKey = true
	}()

	if m.selectedIndex >= len(m.menuList) {
		return
	}

	newTitle := m.menuList[m.selectedIndex]
	stackItem := &menuStackItem{
		menuList:      m.menuList,
		selectedIndex: m.selectedIndex,
		menuCurPage:   m.menuCurPage,
		menuTitle:     m.menuTitle,
		menu:          m.menu,
	}
	m.menuStack.Push(stackItem)

	menu := m.menu.SubMenu(m, m.selectedIndex)
	if menu == nil {
		m.menuStack.Pop()
		return
	}

	if enterMenuHook := menu.BeforeEnterMenuHook(); enterMenuHook != nil {
		loading := NewLoading(m)
		loading.start()
		if res := enterMenuHook(m); !res {
			loading.complete()
			m.menuStack.Pop() // 压入的重新弹出
			return
		}

		// 如果位于正在播放的菜单中，更新播放列表
		if menu.GetMenuKey() == m.player.playingMenuKey {
			if songs, ok := menu.MenuData().([]ds.Song); ok {
				m.player.playlist = songs
			}
		}

		loading.complete()
	}

	menuList := menu.MenuViews()

	m.menu = menu
	m.menuList = menuList
	m.menuTitle = fmt.Sprintf("%s %s", newTitle.Title, SetFgStyle(newTitle.Subtitle, termenv.ANSIBrightBlack))
	m.selectedIndex = 0
	m.menuCurPage = 1
}

// 菜单返回
func backMenu(m *NeteaseModel) {
	m.isListeningKey = false
	defer func() {
		m.isListeningKey = true
	}()

	if m.menuStack.Len() <= 0 {
		return
	}

	stackItem := m.menuStack.Pop()
	if backMenuHook := m.menu.BeforeBackMenuHook(); backMenuHook != nil {
		loading := NewLoading(m)
		loading.start()
		if res := backMenuHook(m); !res {
			loading.complete()
			m.menuStack.Push(stackItem) // 弹出的重新压入
			return
		}
		loading.complete()
	}

	stackMenu, ok := stackItem.(*menuStackItem)
	if !ok {
		return
	}

	m.menuList = stackMenu.menuList
	m.menu = stackMenu.menu
	m.menuTitle = stackMenu.menuTitle
	m.selectedIndex = stackMenu.selectedIndex
	m.menuCurPage = stackMenu.menuCurPage
}

// 空格监听
func spaceKeyHandle(m *NeteaseModel) {
	var (
		songs         []ds.Song
		inPlayingMenu = m.player.InPlayingMenu()
	)
	if inPlayingMenu && !m.menu.ResetPlaylistWhenPlay() {
		songs = m.player.playlist
	} else {
		if data, ok := m.menu.MenuData().([]ds.Song); ok {
			songs = data
		}
	}

	selectedIndex := m.selectedIndex
	if !m.menu.IsPlayable() || len(songs) == 0 || m.selectedIndex > len(songs)-1 {
		if m.player.curSongIndex > len(m.player.playlist)-1 {
			return
		}

		switch m.player.State {
		case utils.Paused:
			m.player.Resume()
		case utils.Playing:
			m.player.Paused()
		case utils.Stopped:
			_ = m.player.PlaySong(m.player.playlist[m.player.curSongIndex], DurationNext)
		}

		return
	}

	if inPlayingMenu && songs[selectedIndex].Id == m.player.playlist[m.player.curSongIndex].Id {
		switch m.player.State {
		case utils.Paused:
			m.player.Resume()
		case utils.Playing:
			m.player.Paused()
		}
	} else {
		m.player.curSongIndex = selectedIndex
		m.player.playingMenuKey = m.menu.GetMenuKey()
		m.player.playingMenu = m.menu
		m.player.playlist = songs
		if m.player.mode == PmIntelligent {
			m.player.SetPlayMode("")
		}
		_ = m.player.PlaySong(songs[selectedIndex], DurationNext)
	}

}

// 空格监听
func spaceKeyHandle1(m *NeteaseModel) {
	var (
		songs []ds.Song
	)
	// inPlayingMenu = m.player.InPlayingMenu())
	// if inPlayingMenu && !m.menu.ResetPlaylistWhenPlay() {
	songs = m.player.playlist
	if data, ok := m.menu.MenuData().([]ds.Song); ok {
		// songs = append(songs, data[m.selectedIndex])
		m.player.extra = append(m.player.extra, data[m.selectedIndex])
		utils.Notify(fmt.Sprintf("把 %s 加入临时歌单", data[m.selectedIndex].Name), fmt.Sprintf("临时歌单还有 %d 首歌", len(m.player.extra)), constants.AppGithubUrl)
	}

	// }

	// selectedIndex := m.selectedIndex
	// if !m.menu.IsPlayable() || len(songs) == 0 || m.selectedIndex > len(songs)-1 {
	// 	if m.player.curSongIndex > len(m.player.playlist)-1 {
	// 		return
	// 	}

	// 	switch m.player.State {
	// 	case utils.Paused:
	// 		m.player.Resume()
	// 	case utils.Playing:
	// 		m.player.Paused()
	// 	case utils.Stopped:
	// 		_ = m.player.PlaySong(m.player.playlist[m.player.curSongIndex], DurationNext)
	// 	}

	// 	return
	// }

	// if inPlayingMenu && songs[selectedIndex].Id == m.player.playlist[m.player.curSongIndex].Id {
	// 	switch m.player.State {
	// 	case utils.Paused:
	// 		m.player.Resume()
	// 	case utils.Playing:
	// 		m.player.Paused()
	// 	}
	// } else {
	// 	m.player.curSongIndex = selectedIndex
	// 	m.player.playingMenuKey = m.menu.GetMenuKey()
	// 	m.player.playingMenu = m.menu
	m.player.playlist = songs
	// 	if m.player.mode == PmIntelligent {
	// 		m.player.SetPlayMode("")
	// 	}
	// 	_ = m.player.PlaySong(songs[selectedIndex], DurationNext)
	// }

}

// likePlayingSong like/unlike playing song
func likePlayingSong(m *NeteaseModel, isLike bool) {
	loading := NewLoading(m)
	loading.start()
	defer loading.complete()

	if m.player.curSongIndex >= len(m.player.playlist) {
		return
	}

	if utils.CheckUserInfo(m.user) == utils.NeedLogin {
		NeedLoginHandle(m, func(m *NeteaseModel) {
			likePlayingSong(m, isLike)
		})
		return
	}

	likeService := service.LikeService{
		ID: strconv.FormatInt(m.player.playlist[m.player.curSongIndex].Id, 10),
		L:  strconv.FormatBool(isLike),
	}
	likeService.Like()

	if isLike {
		utils.Notify("已添加到我喜欢的歌曲", m.player.playlist[m.player.curSongIndex].Name, constants.AppGithubUrl)
	} else {
		utils.Notify("已从我喜欢的歌曲移除", m.player.playlist[m.player.curSongIndex].Name, constants.AppGithubUrl)
	}
}

// logout 登出
func logout() {
	table := db.NewTable()
	_ = table.DeleteByKVModel(db.User{})
	utils.Notify("登出成功", "已清理用户信息", constants.AppGithubUrl)
}

// likeSelectedSong like/unlike selected song
func likeSelectedSong(m *NeteaseModel, isLike bool) {
	loading := NewLoading(m)
	loading.start()
	defer loading.complete()

	songs, ok := m.menu.MenuData().([]ds.Song)
	if !ok || m.selectedIndex >= len(songs) {
		return
	}

	if utils.CheckUserInfo(m.user) == utils.NeedLogin {
		NeedLoginHandle(m, func(m *NeteaseModel) {
			likeSelectedSong(m, isLike)
		})
		return
	}

	likeService := service.LikeService{
		ID: strconv.FormatInt(songs[m.selectedIndex].Id, 10),
		L:  strconv.FormatBool(isLike),
	}
	likeService.Like()

	if isLike {
		utils.Notify("已添加到我喜欢的歌曲", songs[m.selectedIndex].Name, constants.AppGithubUrl)
	} else {
		utils.Notify("已从我喜欢的歌曲移除", songs[m.selectedIndex].Name, constants.AppGithubUrl)
	}
}

// trashPlayingSong 标记为不喜欢
func trashPlayingSong(m *NeteaseModel) {
	loading := NewLoading(m)
	loading.start()
	defer loading.complete()

	if m.player.curSongIndex >= len(m.player.playlist) {
		return
	}

	if utils.CheckUserInfo(m.user) == utils.NeedLogin {
		NeedLoginHandle(m, func(m *NeteaseModel) {
			trashPlayingSong(m)
		})
		return
	}

	trashService := service.FmTrashService{
		SongID: strconv.FormatInt(m.player.playlist[m.player.curSongIndex].Id, 10),
	}
	trashService.FmTrash()

	utils.Notify("已标记为不喜欢", m.player.playlist[m.player.curSongIndex].Name, constants.AppGithubUrl)
}

// trashSelectedSong 标记为不喜欢
func trashSelectedSong(m *NeteaseModel) {
	loading := NewLoading(m)
	loading.start()
	defer loading.complete()

	songs, ok := m.menu.MenuData().([]ds.Song)
	if !ok || m.selectedIndex >= len(songs) {
		return
	}

	if utils.CheckUserInfo(m.user) == utils.NeedLogin {
		NeedLoginHandle(m, func(m *NeteaseModel) {
			trashSelectedSong(m)
		})
		return
	}

	trashService := service.FmTrashService{
		SongID: strconv.FormatInt(songs[m.selectedIndex].Id, 10),
	}
	trashService.FmTrash()

	utils.Notify("已标记为不喜欢", songs[m.selectedIndex].Name, constants.AppGithubUrl)
}

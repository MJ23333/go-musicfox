# go-musicfox

**给个star✨吧**

go-musicfox是 [musicfox](https://github.com/anhoder/musicfox) 的重写版，为了解决某些问题，提升体验，因此采用go进行重写。*go初学者，顺便熟悉go，哈哈RwR*

Fork自原仓库，加入临时歌单（下一首播放**队列**），扩展网易云自带的下一首播放，可以在原播放列表的基础上插入选定的歌曲（不必再原歌单内），不限个数。临时歌单播放完成后才会继续播放原来的播放列表。加入快进快退功能（跳10s），按个人喜好修改快捷键。


### 增加通知

目前触发通知的操作有：
* 播放歌曲
* 加入临时歌单
* 自动签到成功
* 添加歌曲到喜欢
* 从喜欢移除
* 标记为不喜欢
* 登出


x && brew link --overwrite go-musicfox
```

目前仅可从源代码编译下载（懒）

## 使用

```sh
$ musicfox
```

| 按键 | 作用 | 备注 |
| --- | --- | --- |
| h/H/LEFT | 左 |  |
| l/L/RIGHT | 右 |  |
| k/K/UP | 上 |  |
| j/J/DOWN | 下 | |
| Esc | 退出程序 | |
| space | 暂停/播放 | |
| [ | 上一曲 | |
| ] | 下一曲 | |
| - | 减小音量 | |
| = | 加大音量 | |
| . | 快进10s | |
| , | 后退10s | |
| n/N/ENTER | 进入选中的菜单 | |
| b/B/q | 返回上级菜单 | |
| w/W | 退出并退出登录 | |
| p | 切换播放方式 | |
| P | 心动模式(仅在歌单中时有效) | |
| r/R | 重新渲染UI | Windows调整窗口大小后，没有事件触发，可以使用该方法手动重新渲染 |
| l | 喜欢当前播放歌曲 | |
| a | 将当前选中歌曲加入临时歌单 | |

## Todo

* 修改从网页下载歌曲的方法（目前暂时存为a.mp3）（主要是因为快进功能需要使用`IO.Seek()`功能）

## 感谢

感谢以下项目及其贡献者们（不限于）：

* [bubbletea](https://github.com/charmbracelet/bubbletea)
* [beep](https://github.com/faiface/beep)
* [musicbox](https://github.com/darknessomi/musicbox)
* [NeteaseCloudMusicApi](https://github.com/Binaryify/NeteaseCloudMusicApi)
* [NeteaseCloudMusicApiWithGo](https://github.com/sirodeneko/NeteaseCloudMusicApiWithGo)
* [gcli](https://github.com/gookit/gcli)
* [go-musicfox](https://github.com/anhoder/go-musicfox)
* ...

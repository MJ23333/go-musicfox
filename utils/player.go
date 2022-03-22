package utils

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
)

// State 播放器状态
type State uint8

const (
	Stopped State = iota
	Paused
	Playing
)

// SongType 歌曲类型
type SongType uint8

const (
	Mp3 SongType = iota
	Wav
	Ogg
	Flac
)

type UrlMusic struct {
	Url      string
	Type     SongType
	Duration time.Duration
}

type Player struct {
	State      State
	Progress   float64
	CurMusic   UrlMusic
	timepass   time.Duration
	ctrl       *beep.Ctrl
	seeker     beep.StreamSeeker
	volume     *effects.Volume
	Timer      *Timer
	timechange time.Duration
	SampleRate beep.SampleRate
	Normal     float64
	timeChan   chan time.Duration
	done       chan struct{}
	musicChan  chan UrlMusic
	speedy     *beep.Resampler
}

func NewPlayer() *Player {
	player := new(Player)
	player.timeChan = make(chan time.Duration)
	player.done = make(chan struct{})
	player.musicChan = make(chan UrlMusic)
	player.ctrl = &beep.Ctrl{
		Paused: false,
	}
	player.timechange = 0
	player.timepass = 0
	player.volume = &effects.Volume{
		Base:   2,
		Silent: false,
	}
	// player.seeker=&beep.StreamSeeker
	player.speedy = beep.ResampleRatio(4, 1, player.volume)
	player.Normal = player.speedy.Ratio()
	go func() {
		player.listen()
	}()

	return player
}

// listen 开始监听
func (p *Player) listen() {
	p.SampleRate = 44100
	sampleRate := p.SampleRate
	err := speaker.Init(sampleRate, sampleRate.N(time.Second/30))
	if err != nil {
		panic(err)
	}
	done := make(chan bool)

	var (
		streamer, oldStreamer beep.StreamSeekCloser
	)

	for {
		select {
		case <-done:
			p.State = Stopped
			p.pushDone()
			break
		case p.CurMusic = <-p.musicChan:
			var (
				resp   *http.Response
				format beep.Format
			)

			speaker.Clear()

			// 关闭旧响应body
			if resp != nil {
				_ = resp.Body.Close()
			}

			// 关闭旧计时器
			if p.Timer != nil {
				p.Timer.Stop()
			}
			p.Progress = 0

			resp, err = http.Get(p.CurMusic.Url)
			if err != nil {
				p.pushDone()
				break
			}
			rc, err1 := os.Create("a" + ".mp3")
			if err1 != nil {
				panic(err1)
			}
			defer rc.Close()

			// b, err := io.ReadAll(resp.Body)
			_, err = io.Copy(rc, resp.Body)
			rc.Close()
			rc, err = os.Open("a" + ".mp3")
			if err != nil {
				panic(err)
			}
			defer rc.Close()
			oldStreamer = streamer
			switch p.CurMusic.Type {
			case Mp3:
				streamer, format, err = mp3.Decode(rc)
			case Wav:
				streamer, format, err = wav.Decode(rc)
			case Ogg:
				streamer, format, err = vorbis.Decode(rc)
			case Flac:
				streamer, format, err = flac.Decode(rc)
			default:
				p.pushDone()
				break
			}
			if err != nil {
				p.pushDone()
				break
			}

			p.State = Playing
			p.seeker = streamer
			newStreamer := beep.Resample(3, format.SampleRate, sampleRate, p.seeker)
			// fmt.Printf("%d", format.SampleRate)
			p.ctrl.Streamer = beep.Seq(newStreamer, beep.Callback(func() {
				done <- true
			}))
			p.volume.Streamer = p.ctrl

			p.ctrl.Paused = false
			// speed:=beep.ResampleRatio(4,1,p.volume)
			speaker.Play(p.volume)
			p.timechange = 0
			// 启动计时器
			p.Timer = New(Options{
				Duration:       24 * time.Hour,
				TickerInternal: 200 * time.Millisecond,
				OnRun:          func(started bool) {},
				OnPaused:       func() {},
				OnDone:         func(stopped bool) {},
				OnTick: func() {
					if p.CurMusic.Duration > 0 {
						p.Progress = float64(float64(p.Timer.Passed()+p.timechange) * 100 / float64(p.CurMusic.Duration))
						p.timepass = p.Timer.Passed() + p.timechange
					}
					select {
					case p.timeChan <- p.Timer.Passed() + p.timechange:
					default:
					}
				},
			})
			go p.Timer.Run()

			// 关闭旧Streamer，避免协程泄漏
			if oldStreamer != nil {
				_ = oldStreamer.Close()
			}
		}

	}

}

// Play 播放音乐
func (p *Player) Play(songType SongType, url string, duration time.Duration) {
	music := UrlMusic{
		url,
		songType,
		duration,
	}
	select {
	case p.musicChan <- music:
	default:
	}
}

func (p *Player) pushDone() {
	select {
	case p.done <- struct{}{}:
	default:
	}
}

// TimeChan 获取定时器
func (p *Player) TimeChan() <-chan time.Duration {
	return p.timeChan
}

// Done done chan, 如果播放完成往chan中写入struct{}
func (p *Player) Done() <-chan struct{} {
	return p.done
}

// UpVolume 调大音量
func (p *Player) UpVolume() {
	if p.volume.Volume > 0 {
		return
	}

	speaker.Lock()
	p.volume.Silent = false
	p.volume.Volume += 0.5
	speaker.Unlock()
}

// DownVolume 调小音量
func (p *Player) DownVolume() {
	if p.volume.Volume <= -5 {
		speaker.Lock()
		p.volume.Silent = true
		speaker.Unlock()
		return
	}

	speaker.Lock()
	p.volume.Volume -= 0.5
	speaker.Unlock()
}

func (p *Player) Speedup() {
	speaker.Lock()
	pos := p.seeker.Position()
	pos += p.SampleRate.N(10 * time.Second)
	err := p.seeker.Seek(pos)
	if err == nil {
		p.timechange += 10 * time.Second
	}
	speaker.Unlock()
}
func (p *Player) Slowdown() {
	speaker.Lock()
	pos := p.seeker.Position()
	pos -= p.SampleRate.N(10 * time.Second)
	err := p.seeker.Seek(pos)
	if err == nil {
		p.timechange -= 10 * time.Second
	} else {
		panic(err)
	}
	speaker.Unlock()
}

// Paused 暂停播放
func (p *Player) Paused() {
	if p.State == Paused {
		return
	}

	speaker.Lock()
	p.ctrl.Paused = true
	speaker.Unlock()
	p.State = Paused
	p.Timer.Pause()
}

// Resume 继续播放
func (p *Player) Resume() {
	if p.State == Playing {
		return
	}

	speaker.Lock()
	p.ctrl.Paused = false
	speaker.Unlock()
	p.State = Playing
	go p.Timer.Run()
}

// Close 关闭
func (p *Player) Close() {
	p.Timer.Stop()
	speaker.Clear()
}

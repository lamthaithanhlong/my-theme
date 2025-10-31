package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type audioPlayer struct {
	mu       sync.Mutex
	ctrl     *beep.Ctrl
	streamer beep.StreamSeekCloser
	playing  bool
}

func newAudioPlayer(path string) (*audioPlayer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	streamer, format, err := mp3.Decode(file)
	if err != nil {
		file.Close()
		return nil, err
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	player := &audioPlayer{
		ctrl:     &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: true},
		streamer: streamer,
		playing:  false,
	}

	speaker.Play(player.ctrl)

	return player, nil
}

func (p *audioPlayer) play() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.playing {
		return
	}

	p.streamer.Seek(0)
	p.ctrl.Paused = false
	p.playing = true
}

func (p *audioPlayer) stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.playing {
		return
	}

	p.ctrl.Paused = true
	p.streamer.Seek(0)
	p.playing = false
}

func (p *audioPlayer) state() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.playing
}

func (p *audioPlayer) close() error {
	p.stop()
	return p.streamer.Close()
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("khong the ghi phan hoi json: %v", err)
	}
}

func main() {
	player, err := newAudioPlayer("theme.mp3")
	if err != nil {
		log.Fatalf("khong the mo theme.mp3: %v", err)
	}
	defer player.close()

	player.play()

	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		player.play()
		writeJSON(w, map[string]any{"playing": true})
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		player.stop()
		writeJSON(w, map[string]any{"playing": false})
	})

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, map[string]any{"playing": player.state()})
	})

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	log.Println("Dang lang nghe tai http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server loi: %v", err)
	}
}

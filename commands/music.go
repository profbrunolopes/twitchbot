package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ChicoCodes/twitchbot/messages"
)

type gmdpPlayback struct {
	Playing bool
	Song    struct {
		Title    string
		Artist   string
		Album    string
		AlbumArt string
	}
}

func music(_ []string, notification *messages.Notification) {
	const filePath = "${HOME}/Library/Application Support/Google Play Music Desktop Player/json_store/playback.json"
	f, err := os.Open(os.ExpandEnv(filePath))
	if err != nil {
		noMusic(notification)
		return
	}
	defer f.Close()
	var playback gmdpPlayback
	err = json.NewDecoder(f).Decode(&playback)
	if err != nil {
		noMusic(notification)
		return
	}
	if !playback.Playing {
		noMusic(notification)
		return
	}
	notification.Reply(fmt.Sprintf("/me now playing: %s - %s", playback.Song.Artist, playback.Song.Title))
}

func noMusic(notification *messages.Notification) {
	notification.Reply("nenhuma música tocando no momento")
}

package main

import (
	"math/rand"
	"os/exec"
	"time"
)

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func YoutubeDownload(videoID string, ran string, ytdlp string) string {
	cmd := exec.Command(ytdlp, "--format", "mp4", "-o", ran+".mp4", "https://youtube.com/watch?v="+videoID)
	cmd.Start()
	cmd.Wait()

	return ran + ".mp4"
}

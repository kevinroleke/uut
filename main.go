package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/chris124567/zoomer/pkg/zoom"
	"github.com/go-shiori/obelisk"
)

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Downloads Youtube video, uploads it to Google Drive, returns the link.
func video(videoId string, ytdlp string) string {
	fn := YoutubeDownload(videoId, RandomString(10), ytdlp)

	f, err := os.Open(fn)
	HandleErr(err)

	defer f.Close()

	service, err := GetService()
	HandleErr(err)

	dirId := "1AA54mHWmKZJHDd2-zhi4_n7VBU7igEYp"

	file, err := CreateFile(service, fn, "video/mp4", f, dirId)
	os.Remove(fn)
	HandleErr(err)

	return GetViewLink(service, file.Id, true)
}

// Request search term using invidious link, save the page as HTML with styles and images, upload to Google Drive, return the link.
func search(term string, page string, invidiousLink string) string {
	req := obelisk.Request{
		URL: invidiousLink + "/search?q=" + term + "&page=" + page,
	}

	arc := obelisk.Archiver{EnableLog: false}
	arc.Validate()

	result, _, err := arc.Archive(context.Background(), req)
	HandleErr(err)

	// convert html to string, replace video links with prompt that supplies video idea.
	html := string(result[:])
	re := regexp.MustCompile(regexp.QuoteMeta(invidiousLink) + `/watch\?v=(.*)">`)
	s := re.ReplaceAllString(html, "javascript:prompt('Video ID', '$1')\">")

	service, err := GetService()
	HandleErr(err)

	dirId := "1AA54mHWmKZJHDd2-zhi4_n7VBU7igEYp"

	file, err := CreateFile(service, term+".html", "text/html", strings.NewReader(s), dirId)
	HandleErr(err)

	return GetViewLink(service, file.Id, false)
}

func main() {
	var meetingNumber = flag.String("meetingNumber", "", "Meeting number")
	var meetingPassword = flag.String("password", "", "Meeting password")
	var invidiousLink = flag.String("invidious", "https://inv.riverside.rocks", "Invidious Instance URL")
	var dirId = flag.String("folder", "1AA54mHWmKZJHDd2-zhi4_n7VBU7igEYp", "Google Drive Folder ID")
	var hwid = flag.String("hwid", "581bd97f-3d7b-407c-b2bd-c6f28e90f5ce", "Hardware ID for Zoom API")
	var ytdlp = flag.String("ytdlp", "/usr/local/bin/yt-dlp", "YTDLP Location")

	flag.Parse()

	if *meetingNumber == "" || *meetingPassword == "" {
		panic("Usage: ./uut --meetingNumber <zoom meeting number> --password <zoom meeting password> --folder <google drive folder ID> --invidious <optional invidious link> --hwid <optional HWID as a guid>")
	}

	apiKey := os.Getenv("ZOOM_JWT_API_KEY")
	apiSecret := os.Getenv("ZOOM_JWT_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		panic("Please supply the env variables ZOOM_JWT_API_KEY and ZOOM_JWT_API_SECRET")
	}

	_, err := os.Stat(*ytdlp)
	if os.IsNotExist(err) {
		panic("Please install ytdtl or specify its location with --ytdlp")
	}

	// For some reason this module spams logs to stdout
	log.SetOutput(ioutil.Discard)

	session, err := zoom.NewZoomSession(*meetingNumber, *meetingPassword, "uut", *hwid, "", apiKey, apiSecret)
	HandleErr(err)

	meetingInfo, cookieString, err := session.GetMeetingInfoData()
	HandleErr(err)

	websocketUrl, err := session.GetWebsocketUrl(meetingInfo, false)
	HandleErr(err)

	fmt.Println("[*] WS: " + websocketUrl)

	// Establish the websocket. On new join, send the help message. On new message, handle with below function.
	err = session.MakeWebsocketConnection(websocketUrl, cookieString, func(session *zoom.ZoomSession, message zoom.Message) error {
		switch m := message.(type) {
		case *zoom.ConferenceRosterIndication:
			for _, person := range m.Add {
				if person.ID != session.JoinInfo.UserID {
					session.SendChatMessage(zoom.EVERYONE_CHAT_ID, "!search <search terms> <@page12>  !video <video_id>   "+string(person.Dn2)+"!")
				}
			}
			return nil
		case *zoom.ConferenceChatIndication:
			return handleChatMessage(session, m, string(m.Text), *invidiousLink, *dirId, *ytdlp)
		default:
			return nil
		}
	})

	HandleErr(err)
	fmt.Println("[*] Done")
}

var MESSAGE_PREFIX string = "!"

func handleChatMessage(session *zoom.ZoomSession, body *zoom.ConferenceChatIndication, messageText string, invidiousLink string, dirId string, ytdlp string) error {
	// Only care about commands starting with prefix
	if !strings.HasPrefix(messageText, MESSAGE_PREFIX) {
		return nil
	}
	messageText = strings.TrimPrefix(messageText, MESSAGE_PREFIX)

	words := strings.Fields(messageText)
	wordsCount := len(words)
	if wordsCount < 1 {
		return errors.New("no command provided after prefix")
	}
	args := words[1:]
	argsCount := len(args)

	switch words[0] {
	case "video":
		if argsCount > 0 {
			// This is in an anonymous goroutine so multiple commands can be issued at once
			go func(session *zoom.ZoomSession, body *zoom.ConferenceChatIndication, vidId string, ytdlp string) {
				session.SendChatMessage(body.DestNodeID, "Working on it!")
				fmt.Println("[*] Downloading video: " + vidId)
				link := video(vidId, ytdlp)
				fmt.Println("[*] Finished video: " + link)
				session.SendChatMessage(body.DestNodeID, vidId+": "+link)
			}(session, body, args[0], ytdlp)
		}
	case "search":
		if argsCount > 0 {
			// Extract page request (!search <terms> @p<page num>)
			page := "1"
			if strings.HasPrefix(args[len(args)-1], "@page") {
				page = strings.Split(args[len(args)-1], "@page")[1]
				args = args[:len(args)-1]
			}

			term := strings.Join(args, " ")

			fmt.Println("[*] Searching '" + term + "' on page " + page)
			go func(session *zoom.ZoomSession, body *zoom.ConferenceChatIndication, term string, page string) {
				link := search(term, page, invidiousLink)
				fmt.Println("[*] Finished search: " + link)
				session.SendChatMessage(body.DestNodeID, "'"+term+"': "+link)
			}(session, body, term, page)
		}
	}

	return nil
}

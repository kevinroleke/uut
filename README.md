# Uut

Download and search for Youtube videos from a Zoom meeting chatbot. 

## Installation

- Download Golang and setup your environment
- Clone this repo
- Install modules with `go mod tidy`
- Build the project with `go build`
- Download YT-DLP using the below commands
```sh
sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
sudo chmod a+rx /usr/local/bin/yt-dlp
```
- Create a publicly accessible folder in google drive. Save its ID.
- Enable the Drive API in your Google account
- Run Uut, supplying the --meetingNumber, --password (meeting password) and --folder (google drive folder ID)
- Follow the oauth link and directions

## Why Google Drive?
- The Zoom file upload feature has limited support. It does not even exist on the web client. 
- Google Drive is often whitelisted in restricted networks.

## Why.. in general?
- This is for evading web filters to watch Youtube videos. It's an attempt to disguise the process of evading a filter, by using Zoom and Drive, two apps commonly used for studies.

I believe protocol steganography will eventually be the only way to access the internet freely for many people. 

## Credits
- [devtud on Medium](https://devtud.medium.com/upload-files-in-google-drive-with-golang-and-google-drive-api-d686fb62f884), for information on using the Google Drive API in Golang
- [chris124567/zoomer](https://github.com/chris124567/zoomer), for reverse engineering the Zoom meeting API and writing a Golang library.
- [go-shiori/obelisk](https://github.com/go-shiori/obelisk), for providing a way to archive a webpage with all assets in one HTML file.

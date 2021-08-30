# Casigo

"Casi" = Fast (zoom)

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
- Run casigo, supplying the --meetingNumber, --password (meeting password) and --folder (google drive folder ID)
- Follow the oauth link and directions

## Why Google Drive?
- The Zoom file upload feature has limited support. It does not even exist on the web client. 
- Google Drive is often whitelisted in restricted networks.

## Why.. in general?
- This is for evading web filters to watch Youtube videos. It's an attempt to disguise the process of evading a filter, by using Zoom and Drive, two apps commonly used for studies. 

# Slackspot 

Slackspot is a slack integration that allows users of a slack workspace to interact with a Spotify account. 

---
## Usage

### Set up:
* Create Spotify app
	* Save client ID
	* Specify auth callback url (`{SLACKSPOT_HOST}/callback`)
* Create Slack app
	* Save Slack Secret Key
	* Create incoming webhook for admin alerts
	* Add `/nowplaying` slash command
	* Add `/recentlyplayed` slash command
* Optional: Add SSL cert and key to project root for https connections. 

### Start server:
1. Clone repo
2. `cd path/to/slackspot`
3. `dep ensure`
4. `docker-compose up`

---
## Endpoints 

### `/nowplaying`
Returns a json object containing metadata about the song currently playing or 404 when there is no music playing.
```json
{
	"title": "Song Title", 
	"artist": "Song Artist(s)", 
	"images": [ 
                 { 
                    "height":640,
                    "width":640,
                    "url":"https://i.scdn.co/image/bb05317292465b8809b29c00906c1a4b6a226194"
                 },
                 { 
                    "height":300,
                    "width":300,
                    "url":"https://i.scdn.co/image/4db4370eb6b2fd2800cc428879143dfc7866180c"
                 },
                 { 
                    "height":64,
                    "width":64,
                    "url":"https://i.scdn.co/image/2bc1fbc23fa717856bd5df9adc4dbf75ed76284f"
                 }
              ] 
} 
```

### `/recentlyplayed`
Returns an array of json objects containing metadata about the most recently played songs. Size of this array is configured using the `SONG_HISTORY_LENGTH` environment variable. 
```json
[{
	"title": "Song Title", 
	"artist": "Song Artist(s)", 
	"images": [ 
                 { 
                    "height":640,
                    "width":640,
                    "url":"https://i.scdn.co/image/bb05317292465b8809b29c00906c1a4b6a226194"
                 },
                 { 
                    "height":300,
                    "width":300,
                    "url":"https://i.scdn.co/image/4db4370eb6b2fd2800cc428879143dfc7866180c"
                 },
                 { 
                    "height":64,
                    "width":64,
                    "url":"https://i.scdn.co/image/2bc1fbc23fa717856bd5df9adc4dbf75ed76284f"
                 }
              ] 
}, 
{
	"title": "Song Title", 
	"artist": "Song Artist(s)", 
	"images": [ 
                 { 
                    "height":640,
                    "width":640,
                    "url":"https://i.scdn.co/image/bb05317292465b8809b29c00906c1a4b6a226194"
                 },
                 { 
                    "height":300,
                    "width":300,
                    "url":"https://i.scdn.co/image/4db4370eb6b2fd2800cc428879143dfc7866180c"
                 },
                 { 
                    "height":64,
                    "width":64,
                    "url":"https://i.scdn.co/image/2bc1fbc23fa717856bd5df9adc4dbf75ed76284f"
                 }
              ] 
}, 
{
	"title": "Song Title", 
	"artist": "Song Artist(s)", 
	"images": [ 
                 { 
                    "height":640,
                    "width":640,
                    "url":"https://i.scdn.co/image/bb05317292465b8809b29c00906c1a4b6a226194"
                 },
                 { 
                    "height":300,
                    "width":300,
                    "url":"https://i.scdn.co/image/4db4370eb6b2fd2800cc428879143dfc7866180c"
                 },
                 { 
                    "height":64,
                    "width":64,
                    "url":"https://i.scdn.co/image/2bc1fbc23fa717856bd5df9adc4dbf75ed76284f"
                 }
              ] 
}] 
```

### `/slack`
For use with slack slash commands. 

### `/login`
Sends a link to log in to a Spotify account with the permission to read currently playing and recently played songs to the Slack Admin webhook.

### `/callback`
This endpoint is used to finish authenticating a Spotify login request. This url is configured using the `AUTH_CALLBACK` environment variable and must also be specified when creating the Spotify app. 

### `/`
Serves files located in `[PROJECT_ROOT]/static` directory

---
## Slack Slash Commands
These commands must be specified when configuring the slash commands for your slack workspace. They should be pointing to the `/slack` endpoint. 

###  `/nowplaying`
Will display the currently playing song. 

### `/recentlyplayed`
Will display the most recently played songs. 

---
## Environment Variables
`SLACK_SIGNING_SECRET`: Secret key assigned by slack to allow access to workspace

`SONG_HISTORY_LENGTH`: Number of songs to display as "recently played". Default value is 3 

`SPOTIFY_ID`: Assigned Client ID used to interact with the Spotify API. 

`SPOTIFY_SECRET`: Secret key used to interact with the Spotify API.

`CONTACT_SLACK_USER`: Used in message displayed to slack users when there is an error processing a request. When a slack user ID or channel ID is provided (`<@userID|#channelID>`) a clickable link will be shown. Default value is "an Admin". 

`SLACK_ADMIN_WEBHOOK`: Webhook used to send messages to administrators. For example a link to log in to a Spotify account on start up. This webhook can be for a channel or for a single user.

`SLACKSPOT_HOST`: Domain name or IP address of slackspot server. Default value is `http://localhost`.
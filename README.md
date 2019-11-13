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
Returns a json object containing metadata about the song currently playing or 204 (No Content) when there is no music playing.
```json
{
   "id":"7oK9VyNzrYvRFo7nQEYkWN",
   "title":"Mr. Brightside",
   "artist":[
      {
         "id":"0C0XlULifJtAgn6ZNCW2eu",
         "name":"The Killers",
         "url":"https://open.spotify.com/artist/0C0XlULifJtAgn6ZNCW2eu"
      }
   ],
   "url":"https://open.spotify.com/track/7oK9VyNzrYvRFo7nQEYkWN",
   "images":[
      {
         "height":640,
         "width":629,
         "url":"https://i.scdn.co/image/ae33d84ad5b1b47f5c7b73c63ca0f1fd4d131b84"
      },
      {
         "height":300,
         "width":295,
         "url":"https://i.scdn.co/image/e8862e61bc38f868f52bd83c5934b5d41e48500b"
      },
      {
         "height":64,
         "width":63,
         "url":"https://i.scdn.co/image/74bcf303c0186b19b873a1bbdc2a6c9b3fd0b90b"
      }
   ],
   "progress":69873,
   "duration":222586
} 
```

### `/recentlyplayed`
Returns an array of json objects containing metadata about the most recently played songs. Size of this array is configured using the `SONG_HISTORY_LENGTH` environment variable or using `length` query parameter. 
```json
[
   {
      "id":"3dFwpxh2yH7C7p9BGEKLVB",
      "title":"Goodies",
      "artist":[
         {
            "id":"2NdeV5rLm47xAvogXrYhJX",
            "name":"Ciara",
            "url":"https://open.spotify.com/artist/2NdeV5rLm47xAvogXrYhJX"
         },
         {
            "id":"4Js9eYwAf9rypNtV8pNSw9",
            "name":"Petey Pablo",
            "url":"https://open.spotify.com/artist/4Js9eYwAf9rypNtV8pNSw9"
         }
      ],
      "url":"https://open.spotify.com/track/3dFwpxh2yH7C7p9BGEKLVB",
      "images":[
         {
            "height":640,
            "width":640,
            "url":"https://i.scdn.co/image/5c0f4e43e219f387384aab1c5567776307776960"
         },
         {
            "height":300,
            "width":300,
            "url":"https://i.scdn.co/image/aea487a488073e69354ffb9cdb157e73668086f2"
         },
         {
            "height":64,
            "width":64,
            "url":"https://i.scdn.co/image/5b522043445b78a8cc4564a03f4a1d851ff52dd3"
         }
      ],
      "progress":0,
      "duration":223000
   },
   {
      "id":"6t1FIJlZWTQfIZhsGjaulM",
      "title":"Video Killed The Radio Star",
      "artist":[
         {
            "id":"057gc1fxmJ2vkctjQJ7Tal",
            "name":"The Buggles",
            "url":"https://open.spotify.com/artist/057gc1fxmJ2vkctjQJ7Tal"
         }
      ],
      "url":"https://open.spotify.com/track/6t1FIJlZWTQfIZhsGjaulM",
      "images":[
         {
            "height":624,
            "width":640,
            "url":"https://i.scdn.co/image/47fa5d5bbb5d3c1a397fda75c286a7a86d002e4e"
         },
         {
            "height":293,
            "width":300,
            "url":"https://i.scdn.co/image/cff3144f9066ec5beacfb50883ff1f32c7505c5f"
         },
         {
            "height":62,
            "width":64,
            "url":"https://i.scdn.co/image/eb9347b503acf9b58470b9d387dc66dd04339d22"
         }
      ],
      "progress":0,
      "duration":253800
   },
   {
      "id":"1yTQ39my3MoNROlFw3RDNy",
      "title":"Say You'll Be There",
      "artist":[
         {
            "id":"0uq5PttqEjj3IH1bzwcrXF",
            "name":"Spice Girls",
            "url":"https://open.spotify.com/artist/0uq5PttqEjj3IH1bzwcrXF"
         }
      ],
      "url":"https://open.spotify.com/track/1yTQ39my3MoNROlFw3RDNy",
      "images":[
         {
            "height":640,
            "width":640,
            "url":"https://i.scdn.co/image/f9dac1591869800af56100db4d69f28998cf1f06"
         },
         {
            "height":300,
            "width":300,
            "url":"https://i.scdn.co/image/dc265db757daf10739890ead6ee1c88415c3bf33"
         },
         {
            "height":64,
            "width":64,
            "url":"https://i.scdn.co/image/0433314818206d69e9aa18e91a6abe98eb6403ba"
         }
      ],
      "progress":0,
      "duration":235973
   }
] 
```

### `/slack`
For use with slack slash commands. 

### `/login`
Sends a link to log in to a Spotify account with the permission to read currently playing and recently played songs to the Slack Admin webhook.

### `/callback`
This endpoint is used to finish authenticating a Spotify login request. This url must be specified when creating the Spotify app. 

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

`SONG_HISTORY_LENGTH`: Number of songs to display as "recently played" when no query parameter is provided. Default value is 3 

`SPOTIFY_ID`: Assigned Client ID used to interact with the Spotify API. 

`SPOTIFY_SECRET`: Secret key used to interact with the Spotify API.

`CONTACT_SLACK_USER`: Used in message displayed to slack users when there is an error processing a request. When a slack user ID or channel ID is provided (`<@userID|#channelID>`) a clickable link will be shown. Default value is "an Admin". 

`SLACK_ADMIN_WEBHOOK`: Webhook used to send messages to administrators. For example a link to log in to a Spotify account on start up. This webhook can be for a channel or for a single user.

`SLACKSPOT_HOST`: Domain name or IP address of slackspot server. Default value is `http://localhost`.
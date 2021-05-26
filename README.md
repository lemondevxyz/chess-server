# chess-server
A server implementation of the game Chess. Websocket for live board updates and invites, and http for client commands/requests...

## demo
soon... :)

## orders
Orders are basically generic models that could be used for updates(server->client) *and/or* command(client->server). 

In the beginning it used to be 2 different packages, but making them as one helped achieved more consistency when modifying the fields.

### updates
Updates are sent via the Websocket connection, http pools seemed more complicated.

### commands
commands are sent via http, since http can report failure of request.

## authentication
To enable an authentication method, define the following enviroment variables: [$PLATFORM_CLIENT_ID, $PLATFORM_CLIENT_SECRET, $PLATFORM_REDIRECT]. 

$PLATFORM can be one of three: 'DISCORD', 'GOOGLE', 'GITHUB'.

Example:
```
DISCORD_CLIENT_ID=''
DISCORD_CLIENT_SECRET=''
DISCORD_REDIRECT=''
```

Note: You can test each platform's oauth credentials via /api/v1/platform/private. If you are logged in, it will show your user's information.

## matchup
Matchups are created via the invite system. A user invites an available user, and if that other user accepts then a match is created.

Invites expire after 30 seconds.

## deployment
Build the server using `build.sh`, and deploy it as a standalone executable then run it in the background(tmux or in a service file).

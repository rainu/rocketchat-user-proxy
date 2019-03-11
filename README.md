# rocket.chat User Proxy

A REST-Service which provides some basic endpoints to send [Rocket.Chat](https://rocket.chat/)-Messages 
to a **user** or a **room**. It use the [realtime-api](https://rocket.chat/docs/developer-guides/realtime-api/) from 
rocket.chat. Therefore it needs a valid rocket.chat-user.

## Installation

Docker:

```sh
docker run -p 8080:8080 \
       -e "CHAT_HOST=<chat-host>" \
       -e "CHAT_USERNAME=<chat-username>" \
       -e "CHAT_PASSWORD=<chat-user-password>" \
       rainu/rocketchat-user-proxy
```

You have to replace **&lt;chat-host&gt;**, **&lt;chat-username&gt;**, **&lt;chat-user-password&gt;** with your own values.

## Usage example

Sends a message to the user **rainu**
```sh
curl -X POST -v localhost:8080/api/v1/send/u/rainu --data 'Hello rainu!'
```

Sends a message to the room **public**
```sh
curl -X POST -v localhost:8080/api/v1/send/r/public --data 'Hello public World!'
```

Trigger the user **rainu** (sends a message and delete them immediately)
```sh
curl -X POST -v localhost:8080/api/v1/trigger/u/rainu --data 'SPAM!!!'
```

Trigger the room **public** (sends a message and delete them immediately)
```sh
curl -X POST -v localhost:8080/api/v1/trigger/r/public --data 'SPAM!!!'
```

## Documentation

### Configuration

| ENV-Variable        | CLI-Option-Name      | Default-Value | required | Description  |
| ------------------- |----------------------|:-------------:|:--------:| -------------|
| BIND_PORT           | --bind-port          | 8080          | false    | The port where the service listen on |
| CHAT_WS_URL         | --chat-ws-url        |               | true - if hostname is not set | The websocket url of the rocket.chat instance |
| CHAT_HOST           | --chat-host          |               | true - if ws-url is not set | The hostname of the rocket.chat instance |
| CHAT_USERNAME       | --chat-host          |               | true     | The username - this is the user which sends the messages  |
| CHAT_PASSWORD       | --chat-password      |               | true - if password-hash is not set | The user's password (plain) |
| CHAT_PASSWORD_HASH  | --chat-password-hash |               | true - if password is not set | The user's password hash (sha256) |


### API

| Method  | Path      | Variables     | Body |  Description  |
| ------- | --------- | ------------- | ---- | ------------- |
| POST | /api/v1/send/u/${username} | username - the recipient of the message | 1:1 the message to send | Sends a message to the given user. |
| POST | /api/v1/send/r/${room} | room - the target room of the message | 1:1 the message to send | Sends a message to the given room/channel. |
| POST | /api/v1/trigger/u/${username} | username - the recipient of the message | 1:1 the message to send | Sends a message to the given user and delete them immediately. |
| POST | /api/v1/trigger/r/${room} | room - the target room of the message | 1:1 the message to send | Sends a message to the given room/channel and delete them immediately. |

## Development setup

The following scriptlet shows how to setup the project and build from source code.

```sh
mkdir -p ./workspace/src
export GOPATH=./workspace

cd ./workspace/src
git clone git@github.com:rainu/rocketchat-user-proxy.git

cd rocketchat-user-proxy
go get ./...
go build -ldflags -s -a -installsuffix cgo ./cmd/proxy/
```

## Release History
* 0.0.3
    * Transform to go-modules project
* 0.0.2
    * Endpoint for trigger a user
    * Endpoint for trigger a whole room
* 0.0.1
    * Endpoint for sending a message to a user
    * Endpoint for sending a message to a room

## Meta

Distributed under the MIT license. See ``LICENSE`` for more information.

### Intention

I searched a **simple** and easy-to-use way to send a message to a user. So that i can use it for example in 
shell-scripts to inform a user/group of the progress of these script. But i don't want to provide the credentials and 
such stuff in my scripts. So i decided to develop this little project. Now i can simple use curl to send a message :)

## Contributing

1. Fork it (<https://github.com/rainu/rocketchat-user-proxy/fork>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request

package client

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"rocketchat-user-proxy/log"
	"rocketchat-user-proxy/messages"
	"sync"
)

type RocketChat interface {
	//Establish the connection an waits for in-/output messages
	Start() error

	//Stops the client. It blocks until all internal routines are finished
	Stop()

	//LogsIn a user with the given plain password
	LoginWithPassword(string, string)

	//LogIn a user with the given sha-256 hashed password
	LoginWithHash(string, string)

	//LogsOut the user
	Logout()

	//Send a direct message to the given recipients
	SendDirectMessage(string, ...string)

	//Send a message to the given rooms
	SendRoomMessage(string, ...string)

	//Delete a message by the given message id
	DeleteMessage(string)

	//Trigger a user :D
	TriggerUser(string, ...string)

	//Trigger a whole room :D
	TriggerRoom(string, ...string)
}

type rcClient struct {
	url string
	con *websocket.Conn

	chanIn chan interface{}

	history History

	userDictionary map[string]string
	roomDictionary map[string]string

	chanReceiverClose chan interface{}
	waitGroup         sync.WaitGroup
}

func NewRocketChat(url string) RocketChat {
	return &rcClient{
		url:            url,
		history:        NewHistory(),
		userDictionary: make(map[string]string),
		roomDictionary: make(map[string]string),
	}
}

func (rc *rcClient) Start() error {
	var err error
	rc.con, _, err = websocket.DefaultDialer.Dial(rc.url, nil)
	if err != nil {
		return errors.Wrap(err, "Could not connect to RocketChat!")
	}
	log.Info.Println("Establish connection to " + rc.url)

	//first message have to be an "connection message"
	err = rc.con.WriteJSON(messages.NewConnect())
	if err != nil {
		return errors.Wrap(err, "Failed to send the Connect-Message!")
	}

	rc.startSender()
	rc.startReceiver()

	return nil
}

func (rc *rcClient) Stop() {
	//send close signal
	close(rc.chanIn)
	rc.chanReceiverClose <- 1

	//wait for go routines to finish
	rc.waitGroup.Wait()

	rc.con.Close()
}

func (rc *rcClient) startSender() {
	rc.chanIn = make(chan interface{}, 1)
	rc.waitGroup.Add(1)

	go func() {
		defer rc.waitGroup.Done()
		for {
			call, open := <-rc.chanIn
			if !open {
				//close signal received: that means we have to go
				return
			}
			log.Trace.Printf("[OUT] %+v\n", call)
			err := rc.con.WriteJSON(call)

			if err != nil {
				log.Error.Printf("Error sending Message: %v", err)
			}
		}
	}()
}

func (rc *rcClient) startReceiver() {
	rc.chanReceiverClose = make(chan interface{}, 1)
	rc.waitGroup.Add(1)

	go func() {
		defer rc.waitGroup.Done()

		type response struct {
			messageType int
			message     []byte
			err         error
		}
		internalChan := make(chan *response)

		for {

			//wrapper go func to use the channel functionality
			go func() {
				resp := &response{}
				resp.messageType, resp.message, resp.err = rc.con.ReadMessage()

				internalChan <- resp
			}()

			select {
			case <-rc.chanReceiverClose:
				//close signal received: that means we have to go
				return
			case resp := <-internalChan:
				if resp.err == nil {
					rc.handleMessage(string(resp.message))
				}
			}
		}
	}()
}

func (rc *rcClient) handleMessage(message string) {
	log.Trace.Printf("[IN] %s\n", message)

	//try to convert in GeneralMessage
	genResp := &messages.GeneralResponse{}
	err := json.Unmarshal([]byte(message), genResp)

	if err == nil {
		switch genResp.Msg {
		case "ping":
			//we have to pong :D
			rc.chanIn <- messages.NewPingResponse()
		default:
			rc.history.AddIncomingMessage(message)
		}
	}
}

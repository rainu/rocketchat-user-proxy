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

	//Send a message to the given recipients
	SendMessage(string, ...string)
}

type rcClient struct {
	url string
	con *websocket.Conn

	chanIn  chan interface{}
	chanOut chan []byte

	history History

	userDictionary map[string]string

	chanSenderClose   chan interface{}
	chanReceiverClose chan interface{}
	waitGroup         sync.WaitGroup
}

func NewRocketChat(url string) RocketChat {
	return &rcClient{
		url:            url,
		history:        NewHistory(),
		userDictionary: make(map[string]string),
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
	rc.chanSenderClose <- 1
	rc.chanReceiverClose <- 1
	rc.con.Close()

	//wait for go routines to finish
	rc.waitGroup.Wait()
}

func (rc *rcClient) startSender() {
	rc.chanIn = make(chan interface{}, 1)
	rc.chanSenderClose = make(chan interface{}, 1)
	rc.waitGroup.Add(1)

	go func() {
		defer rc.waitGroup.Done()
		for {
			select {
			case <-rc.chanSenderClose:
				//close signal received: that means we have to go
				return
			case call := <-rc.chanIn:
				log.Trace.Printf("[OUT] %+v\n", call)
				err := rc.con.WriteJSON(call)

				if err != nil {
					log.Error.Printf("Error sending Message: %v", err)
				}
			}
		}
	}()
}

func (rc *rcClient) startReceiver() {
	rc.chanOut = make(chan []byte, 100)
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
					providable := rc.handleMessageBeforeProvide(string(resp.message))

					if providable {
						rc.chanOut <- resp.message
					}
				}
			}
		}
	}()
}

func (rc *rcClient) handleMessageBeforeProvide(message string) bool {
	log.Trace.Printf("[IN] %s\n", message)

	//try to convert in GeneralMessage
	genResp := &messages.GeneralResponse{}
	err := json.Unmarshal([]byte(message), genResp)

	if err == nil {
		switch genResp.Msg {
		case "ping":
			//we have to pong :D
			rc.chanIn <- messages.NewPingResponse()
			return false
		default:
			rc.history.AddIncomingMessage(message)
		}
	}

	return true
}

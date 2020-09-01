package ws

import (
	"context"
	"errors"
	"github.com/alexdemen/ws_chat/domain"
	"net/http"
	"nhooyr.io/websocket"
)

type Client struct {
	msgs chan []byte
	name string
}

func NewWSClient() *Client {
	return &Client{msgs: make(chan []byte, 10), name: "One"}
}

func (W *Client) SendMessage(m domain.Message) error {
	W.msgs <- []byte(m.Text)
	return nil
}

func (W *Client) Close() {
	close(W.msgs)
}

type Handler struct {
	sender *domain.Sender
}

func NewHandler(sender *domain.Sender) *Handler {
	return &Handler{sender: sender}
}

func (ws Handler) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(wr, r, nil)
	if err != nil {
		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	err = ws.process(r.Context(), c)
	if errors.Is(err, context.Canceled) {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
		return
	}
	if err != nil {
		return
	}
}

func (ws *Handler) process(ctx context.Context, c *websocket.Conn) error {
	ctx = c.CloseRead(ctx)

	client := NewWSClient()
	err := ws.sender.AddClient(client)
	defer ws.sender.DeleteClient(client)
	err = c.Write(ctx, websocket.MessageText, []byte("hello"))
	if err != nil {
		return err
	}

	for {
		select {
		case msg := <-client.msgs:
			err = c.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

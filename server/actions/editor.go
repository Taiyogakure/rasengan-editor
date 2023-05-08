package actions

import (
	"log"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// EditorHandler: serves up the editor page.
func EditorHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("editor/index.plush.html"))
}

// WSHandler: handles websocket connection
func WSHandler(c buffalo.Context, h *Hub) error {
	log.Println("Handling Websocket connection from", c.Request().RemoteAddr)

	uid, name, err := h.Authorize(c.Request().URL.Query().Get("token"))

	if err != nil {
		c.Response().WriteHeader(403)
		log.Println(err)
		errors.WithStack(err)
	}

	ws, err := websocket.Upgrade(c.Response(), c.Request(), c.Response().Header(), 1024, 1024)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			errors.WithStack(err)
		}
	}

	client := &Client{hub: h, conn: ws, uid: uid, name: name, buffer: make(chan []byte, 256)}
	client.hub.register <- client

	go client.Reader()
	client.Writer()

	client.hub.unregister <- client
	client.conn.Close()

	return nil
}
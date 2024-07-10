package core

import "github.com/Allenxuxu/gev/plugins/websocket/ws"

type Handler func(c *Conn, data []byte) (messageType ws.MessageType, out []byte)

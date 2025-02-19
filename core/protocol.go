package core

import (
	"errors"
	"net/http"

	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/plugins/websocket/ws"
	"github.com/Allenxuxu/ringbuffer"
	"github.com/gobwas/pool/pbytes"
	"go.uber.org/zap"
)

const (
	upgradedKey     = "ctx:ws_upgraded"
	headerBufferKey = "ctx:header_buf"
)

// Protocol 应用层协议
type Protocol interface {
	Accept(c *Conn, uri string, headers http.Header) error
	Handler(c *Conn, data []byte) (messageType ws.MessageType, out []byte)
	Pack(c *Conn, data []byte) ([]byte, error)
}

type gevProtocol struct {
	upgrade *ws.Upgrader
}

func newGevProtocol(u *ws.Upgrader) gev.Protocol {
	return &gevProtocol{
		upgrade: u,
	}
}

// UnPacket 解析 websocket 协议，返回 header ，payload
func (p *gevProtocol) UnPacket(c *gev.Connection, buffer *ringbuffer.RingBuffer) (ctx interface{}, out []byte) {
	_, ok := c.Get(upgradedKey)
	if !ok {
		var err error
		out, _, err = p.upgrade.Upgrade(c, buffer)
		if err != nil {
			Log(zap.ErrorLevel, "websocket upgrade failed",
				Error(err),
				Addr(c.PeerAddr()),
			)
			return
		}
		c.Set(upgradedKey, true)
		c.Set(headerBufferKey, pbytes.Get(0, ws.MaxHeaderSize-2))
	} else {
		bts, _ := c.Get(headerBufferKey)
		header, err := ws.VirtualReadHeader(bts.([]byte), buffer)
		if err != nil {
			if !errors.Is(err, ws.ErrHeaderNotReady) {
				Log(zap.ErrorLevel, "decode websocket protocol failed",
					Error(err),
					Addr(c.PeerAddr()),
				)
			}
			return
		}
		if buffer.VirtualLength() >= int(header.Length) {
			buffer.VirtualFlush()

			payload := make([]byte, int(header.Length))
			_, _ = buffer.Read(payload)

			if header.Masked {
				ws.Cipher(payload, header.Mask, 0)
			}

			ctx = &header
			out = payload
		} else {
			buffer.VirtualRevert()
		}
	}
	return
}

// Packet 直接返回
func (p *gevProtocol) Packet(c *gev.Connection, data interface{}) []byte {
	return data.([]byte)
}

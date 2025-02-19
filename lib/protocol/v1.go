package protocol

import (
	"net/http"
	"net/url"

	"github.com/jhq0113/gochat/core"
	"github.com/jhq0113/gochat/lib/pogo"
	"github.com/jhq0113/gochat/lib/utils"

	"github.com/Allenxuxu/gev/plugins/websocket/ws"
	"github.com/Allenxuxu/toolkit/convert"
	"github.com/goccy/go-json"
)

const (
	Auth   = `Authorization`
	CtxKey = `ctx:key`
)

type V1 struct {
	handler Handler
	rsa     *utils.Rsa
}

func NewV1(handler Handler, rsa *utils.Rsa) core.Protocol {
	return &V1{
		handler: handler,
		rsa:     rsa,
	}
}

func (v1 *V1) Accept(c *core.Conn, uri string, headers http.Header) error {
	u, err := url.Parse(uri)
	if err != nil {
		//fmt.Printf("parse url failed with error: %v\n", err)
		return err
	}

	auth := u.Query().Get(Auth)
	if len(auth) == 0 {
		//fmt.Printf("authorization param is empty\n")
		return ErrForbidden
	}

	data, err := utils.Base64UrlDecode(convert.StringToBytes(auth))
	if err != nil {
		//fmt.Printf("base64 url decode failed with error: %v\n", err)
		return err
	}

	data, err = v1.rsa.Decrypt(data)
	if err != nil {
		//fmt.Printf("rsa decrypt failed with error: %v\n", err)
		return err
	}

	if len(data) != 32 {
		//fmt.Printf("key length is not 32: %s\n", data)
		return ErrForbidden
	}

	aes, err := utils.NewAesWithBytes(data[:16], data[16:])
	if err != nil {
		//fmt.Printf("create aes failed with error: %v\n", err)
		return err
	}

	c.Set(CtxKey, aes)

	return nil
}

func (v1 *V1) Handler(c *core.Conn, data []byte) (messageType ws.MessageType, out []byte) {
	value, _ := c.Get(CtxKey)
	aes, ok := value.(*utils.Aes)
	if !ok {
		c.Close()
		return
	}

	data, err := utils.Base64UrlDecode(data)
	if err != nil {
		c.Close()
		return
	}

	eventData, err := aes.CbcDecrypt(data)
	if err != nil {
		c.Close()
		return
	}

	event := pogo.AcqEvent()
	defer event.Close()

	if err = json.Unmarshal(eventData, event); err != nil {
		c.Close()
		return
	}

	v1.handler(c, event)
	return
}

func (v1 *V1) Pack(c *core.Conn, data []byte) ([]byte, error) {
	value, _ := c.Get(CtxKey)
	aes, ok := value.(*utils.Aes)
	if !ok {
		return nil, ErrProtocol
	}

	return utils.Base64UrlEncode(aes.CbcEncrypt(data)), nil
}

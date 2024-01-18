package gateway

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"io"
)

// wsCodec is the websocket context to be set
// when a gnet.conn arrived
type wsCodec struct {
	upgraded bool
	buf      bytes.Buffer
	wsMsgBuf wsMessageBuf
}

type wsMessageBuf struct {
	firstHeader *ws.Header
	curHeader   *ws.Header
	cachedBuf   bytes.Buffer
}

func (wc *wsCodec) upgrade(c gnet.Conn) (ok bool, action gnet.Action) {
	if wc.upgraded {
		ok = true
		action = 0
	}

	buf := &wc.buf
	reader := bytes.NewReader(buf.Bytes())
	oldlen := reader.Len()
	logging.Infof("upgrade\n")

	hs, err := ws.Upgrade(struct {
		io.Reader
		io.Writer
	}{
		reader,
		c,
	})
	skipN := oldlen - reader.Len()
	if err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return
		}
		buf.Next(skipN)
		logging.Infof("[ERROR]ws.Upgrade() Error: conn[%v] %v", c.RemoteAddr().String(), err.Error())
		action = gnet.Close
		return
	}
	buf.Next(skipN)
	logging.Infof("conn[%v] upgrade websocket complete! handshake: %v", c.RemoteAddr().String(), hs)
	ok = true
	wc.upgraded = true
	return
}
func (wc *wsCodec) readBuffBytes(c gnet.Conn) gnet.Action {
	size := c.InboundBuffered()
	buf := bufferPoolInstance.Get().(bytes.Buffer)
	read, err := c.Read(buf.Bytes())
	if err != nil {
		logging.Infof("read err! %w", err)
		return gnet.Close
	}
	if read < size {
		logging.Infof("read err! read size %v != expected size %v", read, size)
		return gnet.Close
	}
	wc.buf.Write(buf.Bytes())
	return gnet.None
}

func (wc *wsCodec) readWsMessage() (messages []wsutil.Message, err error) {
	msbuf := &wc.wsMsgBuf
	in := &wc.buf
	for {
		if msbuf.curHeader == nil {
			if in.Len() < ws.MinHeaderSize {
				return
			}
			var head ws.Header
			if in.Len() >= ws.MaxHeaderSize {
				head, err = ws.ReadHeader(in)
				if err != nil {
					return
				}
			} else {
				tmpReader := bytes.NewReader(in.Bytes())
				oldLen := tmpReader.Len()
				head, err = ws.ReadHeader(tmpReader)
				skipN := oldLen - tmpReader.Len()
				if err != nil {
					if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) { //数据不完整
						return messages, nil
					}
					in.Next(skipN)
					return nil, err
				}
				in.Next(skipN)
			}
			msbuf.curHeader = &head
			err := ws.WriteHeader(&msbuf.cachedBuf, head)
			if err != nil {
				return nil, err
			}
		}
		dataLen := (int)(msbuf.curHeader.Length)
		if dataLen > 0 {
			if in.Len() >= dataLen {
				_, err = io.CopyN(&msbuf.cachedBuf, in, int64(dataLen))
				if err != nil {
					return
				}
			} else { //数据不完整
				fmt.Println(in.Len(), dataLen)
				logging.Infof("incomplete data")
				return
			}
		}
		if msbuf.curHeader.Fin { //当前 header 已经是一个完整消息
			messages, err = wsutil.ReadClientMessage(&msbuf.cachedBuf, messages)
			if err != nil {
				return nil, err
			}
			msbuf.cachedBuf.Reset()
		} else {
			logging.Infof("The data is split into multiple frames")
		}
		msbuf.curHeader = nil
	}
}

func (wc *wsCodec) Decode(c gnet.Conn) (out []wsutil.Message, err error) {
	logging.Infof("decoding")
	messages, err := wc.readWsMessage()
	if err != nil {
		logging.Infof("Error reading message! %v", err)
		return nil, err
	}
	if messages == nil || len(messages) <= 0 {
		return
	}
	for _, message := range messages {
		if message.OpCode.IsControl() {
			err = wsutil.HandleClientControlMessage(c, message)
			if err != nil {
				return
			}
			continue
		}
		if message.OpCode == ws.OpText || message.OpCode == ws.OpBinary {
			out = append(out, message)
		}
	}
	return
}

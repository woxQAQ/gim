package gateway

import (
	"bytes"
	"errors"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"io"
)

// TODO should I really write codec like this? read more about gobwas/ws

// wsCodec is the websocket context to be set
// when a gnet.conn arrived
type wsCodec struct {
	upgraded bool
	buf      bytes.Buffer
	wsMsgBuf wsMessageBuf
}

// wsMessageBuf is used to
type wsMessageBuf struct {
	firstHeader *ws.Header
	curHeader   *ws.Header
	cachedBuf   bytes.Buffer
}

// upgrade is used to upgrade a gnet.conn to a websocket conn.
// to identify a conn which is upgrade or not, we use upgraded field
// and a websocket conn sees not so different from a gnet conn--just
// diff with a context
func (wc *wsCodec) upgrade(c gnet.Conn) (ok bool, action gnet.Action) {
	// if a wsCodec--a gnet conn's context, is upgrade, we return directly
	if wc.upgraded {
		ok = true
		action = 0
		return
	}

	// as a websocket server, it should receive msg from connection and
	// write something back to it.We use codec's buf as reader and Conn as writer
	buf := &wc.buf
	reader := bytes.NewReader(buf.Bytes())
	oldLen := reader.Len()
	logging.Infof("upgrade\n")
	hs, err := ws.Upgrade(struct {
		io.Reader
		io.Writer
	}{
		reader,
		c,
	})
	// after reader--the codec's buf, upgrade with conn to be a websocket conn
	skipN := oldLen - reader.Len()
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

// readBuffBytes read message from gnet's Conn to Codec's buf
func (wc *wsCodec) readBuffBytes(c gnet.Conn) gnet.Action {
	// read from c
	// get the size of message for validate the read
	size := c.InboundBuffered()
	buffer := bufferPoolInstance.Get().(bytes.Buffer)
	read, err := c.Read(buffer.Bytes())
	if err != nil {
		logging.Infof("read err! %w", err)
		return gnet.Close
	}
	if read < size {
		logging.Infof("read err! read size %v != expected size %v", read, size)
		return gnet.Close
	}

	// write buffer's data to Codec's buf
	wc.buf.Write(buffer.Bytes())
	return gnet.None
}

func (wc *wsCodec) readWsMessage() (messages []wsutil.Message, err error) {
	msbuf := &wc.wsMsgBuf
	in := &wc.buf
	for {
		// reading header from wsCodec's buf
		// after websocket connection established, we don't send the header any more.
		// so the first time we read ws message, we need to read header into codec
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

			// write header into cachedBuf
			err = ws.WriteHeader(&msbuf.cachedBuf, head)
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
				logging.Infof("incomplete data, expected length: %v and actual length: %v", dataLen, in.Len())
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

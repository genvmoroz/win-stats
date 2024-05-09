package http

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

func setDefaultCloseHandler(conn *websocket.Conn, cancel func(), logger logrus.FieldLogger) {
	conn.SetCloseHandler(func(code int, text string) error {
		cancel()
		logger.
			WithField("Code", code).
			WithField("Msg", text).
			Infof("ws connection is closed by the peer")
		return websocket.ErrCloseSent
	})
}

func runDefaultConnectionReader(ctx context.Context, conn *websocket.Conn, logger logrus.FieldLogger) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		code, msg, rErr := conn.ReadMessage()
		switch {
		case errors.Is(rErr, websocket.ErrCloseSent):
			return
		case rErr != nil:
			logger.Errorf("failed to read from the peer: %s", rErr.Error())
		default:
			logger.
				WithField("Code", code).
				WithField("Msg", string(msg)).
				Printf("received message from the peer")
		}
	}
}

func (s *Server) writeMessageWS(conn *websocket.Conn, msgType int, msg string) {
	err := conn.WriteMessage(msgType, []byte(msg))
	if err != nil && !errors.Is(err, websocket.ErrCloseSent) {
		s.logger.Errorf("failed to write the message (%s) into websocket: %s", msg, err.Error())
	}
}

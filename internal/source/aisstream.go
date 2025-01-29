package source

import (
	"encoding/json"
	"iter"
	"log"
	"os"

	aisstream "github.com/aisstream/ais-message-models/golang/aisStream"
	"github.com/gorilla/websocket"
)

type AisStreamSource struct {
	// TODO: add iterator
	ws *websocket.Conn
}

func NewAisStreamSource() (*AisStreamSource, error) {
	url := "wss://stream.aisstream.io/v0/stream"
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	subMsg := aisstream.SubscriptionMessage{
		APIKey:        os.Getenv("API_TOKEN"),
		BoundingBoxes: [][][]float64{{{-90.0, -180.0}, {90.0, 180.0}}}, // bounding box for the entire world
	}

	subMsgBytes, _ := json.Marshal(subMsg)
	if err := ws.WriteMessage(websocket.TextMessage, subMsgBytes); err != nil {
		return nil, err
	}

	return &AisStreamSource{
		ws: ws,
	}, nil
}

func (s *AisStreamSource) Close() {
	s.ws.Close()
}

func (s *AisStreamSource) Read() iter.Seq[aisstream.AisStreamMessage] {
	return func(yield func(aisstream.AisStreamMessage) bool) {
		// items := []model.AisMessage{}
		// for _, v := range items {
		// 	if !yield(v) {
		// 		return
		// 	}
		// }

		for i := 0; i < 1000; i++ {
			_, p, err := s.ws.ReadMessage()
			if err != nil {
				log.Fatalln(err)
			}
			var packet aisstream.AisStreamMessage

			err = json.Unmarshal(p, &packet)
			if err != nil {
				return
			}

			if !yield(packet) {
				return
			}
		}
	}
}

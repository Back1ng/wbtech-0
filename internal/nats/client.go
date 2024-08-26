package natsclient

import (
	"encoding/json"
	"github.com/Back1ng/wbtech-0/internal/entity"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/nats.go"
	"log"
)

type Nats struct {
	conn *nats.Conn
}

func NewClient() (Nats, error) {
	nc, err := nats.Connect(nats.DefaultURL, nats.Name("server"))
	if err != nil {
		log.Fatal("Error when nats.Connect:", err)
	}

	return Nats{
		conn: nc,
	}, err
}

func (n Nats) ListenOrdersSubject(subj string) chan entity.Order {
	ordersChan := make(chan entity.Order)

	_, err := n.conn.Subscribe(subj, func(msg *nats.Msg) {
		var order entity.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Println(err)
		}

		if err := validator.New().Struct(order); err != nil {
			log.Println(err)
			return
		}

		ordersChan <- order
	})
	if err != nil {
		log.Println("n.conn.Subscribe", err)
	}

	return ordersChan
}

func (n Nats) Close() {
	n.conn.Close()
}

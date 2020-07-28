package messages

import (
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
)

type Producer struct {
	client      *twitch.Client
	subscribers map[string]*Subscriber
	connected   bool
	mtx         sync.RWMutex
}

type Subscriber struct {
	id       string
	messages chan<- Message
}

type Message struct {
	User      string
	Text      string
	Timestamp time.Time
}

func NewProducer(twitchChannel string) *Producer {
	producer := &Producer{
		client:      twitch.NewAnonymousClient(),
		subscribers: map[string]*Subscriber{},
	}
	producer.client.OnPrivateMessage(func(msg twitch.PrivateMessage) { producer.privateMessageCb(msg) })
	producer.client.Join(twitchChannel)
	// fix error handling
	go func() {
		err := producer.client.Connect()
		if err != nil {
			panic(err)
		}
	}()
	return producer
}

func (p *Producer) privateMessageCb(msg twitch.PrivateMessage) {
	message := Message{
		User:      msg.User.DisplayName,
		Text:      msg.Message,
		Timestamp: msg.Time,
	}
	for _, subscriber := range p.subscribers {
		subscriber.messages <- message
	}
}

func (p *Producer) Subscribe(id string) (<-chan Message, error) {
	const msgChanCapacity = 10
	messages := make(chan Message, msgChanCapacity)
	p.mtx.Lock()
	defer p.mtx.Unlock()

	// note: should refactor this, don't want to override previously
	// defined subscriber.
	p.subscribers[id] = &Subscriber{id: id, messages: messages}
	return messages, nil
}

func (p *Producer) Unsubscribe(id string) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	delete(p.subscribers, id)
}

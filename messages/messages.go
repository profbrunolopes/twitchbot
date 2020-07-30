package messages

import (
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
)

type Notify func(Message)

type Producer struct {
	client      *twitch.Client
	subscribers map[string]Notify
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

type ProducerOptions struct {
	UserName   string `envconfig:"BOT_USERNAME"`
	OAuthToken string `envconfig:"OAUTH_TOKEN"`
	Channel    string `ignored:"true"`
}

func NewProducer(options *ProducerOptions) *Producer {
	producer := &Producer{
		client:      twitch.NewClient(options.UserName, options.OAuthToken),
		subscribers: map[string]Notify{},
	}
	producer.client.OnPrivateMessage(func(msg twitch.PrivateMessage) { producer.privateMessageCb(msg) })
	producer.client.Join(options.Channel)
	return producer
}

func (p *Producer) privateMessageCb(msg twitch.PrivateMessage) {
	message := Message{
		User:      msg.User.DisplayName,
		Text:      msg.Message,
		Timestamp: msg.Time,
	}
	for _, notify := range p.subscribers {
		go notify(message)
	}
}

func (p *Producer) Subscribe(notify Notify) string {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	id := "id aleatório"
	p.subscribers[id] = notify
	return id
}

func (p *Producer) Unsubscribe(id string) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	delete(p.subscribers, id)
}

func (p *Producer) Start() error {
	return p.client.Connect()
}

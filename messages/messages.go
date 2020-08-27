package messages

import (
	"crypto/rand"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
)

var modNoticeRegexp = regexp.MustCompile(`^The moderators of this channel are: (.+)$`)

type Notification struct {
	Message      Message
	channel      string
	twitchClient *twitch.Client
}

func (n *Notification) Reply(msg string) {
	n.twitchClient.Say(n.channel, msg)
}

type Notify func(Notification)

type Producer struct {
	channel     string
	client      *twitch.Client
	subscribers map[string]Notify
	mtx         sync.RWMutex
	notices     chan twitch.NoticeMessage
}

type User struct {
	DisplayName string
	name        string
	producer    *Producer
}

func (u *User) IsMod() (bool, error) {
	const timeout = 500 * time.Millisecond
	u.producer.say("/mods")
	select {
	case notice := <-u.producer.notices:
		m := modNoticeRegexp.FindStringSubmatch(notice.Message)
		mods := strings.Split(m[1], ",")
		for _, mod := range mods {
			if u.name == strings.TrimSpace(mod) {
				return true, nil
			}
		}
		return false, nil
	case <-time.After(timeout):
		return false, fmt.Errorf("couldn't determine if user is mod in %s", timeout)
	}
}

func (u *User) IsStreamer() bool {
	return u.name == u.producer.channel
}

type Message struct {
	User      User
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
		channel:     options.Channel,
		client:      twitch.NewClient(options.UserName, options.OAuthToken),
		subscribers: map[string]Notify{},
		notices:     make(chan twitch.NoticeMessage),
	}
	producer.client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		producer.privateMessageCb(msg)
	})
	producer.client.OnNoticeMessage(func(msg twitch.NoticeMessage) {
		select {
		case producer.notices <- msg:
		case <-time.After(time.Second):
		}
	})
	producer.client.Join(options.Channel)
	return producer
}

func (p *Producer) privateMessageCb(msg twitch.PrivateMessage) {
	notification := Notification{
		Message: Message{
			User: User{
				DisplayName: msg.User.DisplayName,
				name:        msg.User.Name,
				producer:    p,
			},
			Text:      msg.Message,
			Timestamp: msg.Time,
		},
		channel:      msg.Channel,
		twitchClient: p.client,
	}
	p.mtx.RLock()
	defer p.mtx.RUnlock()
	for _, notify := range p.subscribers {
		go notify(notification)
	}
}

func (p *Producer) Subscribe(notify Notify) (string, error) {
	var buf [8]byte
	n, err := rand.Read(buf[:])
	if err != nil {
		return "", err
	}
	if n != len(buf) {
		return "", io.ErrShortWrite
	}
	p.mtx.Lock()
	defer p.mtx.Unlock()
	id := fmt.Sprintf("%x", buf[:])
	p.subscribers[id] = notify
	return id, err
}

func (p *Producer) Unsubscribe(id string) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	delete(p.subscribers, id)
}

func (p *Producer) Start() error {
	return p.client.Connect()
}

func (p *Producer) say(msg string) {
	p.client.Say(p.channel, msg)
	fmt.Printf("said %q\n", msg)
}

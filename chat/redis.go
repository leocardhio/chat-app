package chat

import (
	"chat-app/config"
	"log"
	"strings"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

var (
	client *redis.Client
	sub *redis.PubSub
	cfg config.Config
)

const (
	channel = "chat"
	users = "chat-users"
)

func init() {
	cfg = config.Cfg

	log.Println("connecting to Redis..")
	redisOpt := redis.Options{
		Network: "tcp",
		Addr: cfg.RedisHost,
		Password: cfg.RedisPwd,
	}
	client = redis.NewClient(&redisOpt)

	_, err := client.Ping().Result()
	if err != nil {
		client.Close()
		log.Fatal("failed to connect to redis", err)
	}

	log.Println("connected to redis", cfg.RedisHost)
	startSubscriber()
}

func SendToChannel(msg string) {
	err := client.Publish(channel, msg).Err()
	if err != nil {
		log.Println("could not publish to channel", err)
	}
}

func startSubscriber() {
	go func() {
		sub = client.Subscribe(channel)
		messages := sub.Channel()
		for message := range messages {
			from := strings.Split(message.Payload, ":")[0]
			for user, peer := range Peers {
				if from != user { // prevent self-sent message
					peer.WriteMessage(websocket.TextMessage, []byte(message.Payload))
				}
			}
		}
	}()
}

func RemoveUser(user string) {
	err := client.SRem(users, user)
	if err != nil {
		log.Println("failed to remove user:", user)
		return
	}

	log.Println("removed user from redis:", user)
}

func Cleanup() {
	for user, peer := range Peers {
		client.SRem(users, user)
		peer.Close()
	}
	log.Println("cleaned up users and sessions...")
	err := sub.Unsubscribe(channel)
	if err != nil {
		log.Println("failed to unsubscribe redis channel subscription:", err)
	}
	err = sub.Close()
	if err != nil {
		log.Println("failed to close redis channel subscription:", err)
	}

	err = client.Close()
	if err != nil {
		log.Println("failed to close redis connection: ", err)
		return
	}
}
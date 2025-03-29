package service

import (
	"chat/internal/handler/ws"
	"chat/internal/model"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
)

type MessageService struct {
	redisClient *redis.Client
}

func NewMessageService(r *redis.Client) *MessageService {
	return &MessageService{r}
}

func (s *MessageService) PublishMessage(roomID uint, msg model.Message) error {
	data, _ := json.Marshal(msg)  //序列化消息
	return s.redisClient.Publish( //发布到redis频道
		context.Background(),
		"chat:room:"+string(roomID), //频道名格式
		data,
	).Err()
}

func (s *MessageService) SubscribeMessages(roomID uint) {
	pubsub := s.redisClient.Subscribe( //订阅指定房间频道
		context.Background(),
		"chat:room:"+string(roomID),
	)

	defer pubsub.Close()

	ch := pubsub.Channel() //获取消息通道
	for msg := range ch {
		var message model.Message
		json.Unmarshal([]byte(msg.Payload), &message) //反序列化
		ws.HubInstance.Broadcast <- message           //转发给hub
	}
}

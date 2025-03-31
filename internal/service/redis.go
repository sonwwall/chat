package service

import (
	"chat/internal/handler/ws"
	"chat/internal/model"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"log"
	"sync"
)

type MessageService struct {
	redisClient *redis.Client
	activeSubs  sync.Map
}

func NewMessageService(r *redis.Client) *MessageService {
	return &MessageService{r, sync.Map{}}
}

func (s *MessageService) PublishMessage(roomID uint, msg model.Message) error {
	data, _ := json.Marshal(msg)  //序列化消息
	return s.redisClient.Publish( //发布到redis频道
		context.Background(),
		fmt.Sprintf("chat:room:%d", roomID), //频道名格式
		data,
	).Err()
}

var roomLocks sync.Map

func (s *MessageService) SubscribeMessages(roomID uint) {
	//生成频道名
	channel := fmt.Sprintf("chat:room:%d", roomID)

	//获取房间专属锁
	lock, _ := roomLocks.LoadOrStore(roomID, &sync.Mutex{})
	mu := lock.(*sync.Mutex)

	mu.Lock()
	defer mu.Unlock()

	//检查重复订阅
	if _, loaded := s.activeSubs.Load(roomID); loaded {
		return
	}

	pubsub := s.redisClient.Subscribe( //订阅指定房间频道
		context.Background(),
		channel,
	)
	log.Printf("已订阅房间%v", roomID)

	s.activeSubs.Store(roomID, pubsub)

	go s.listenMessages(pubsub, roomID)

}

func (s *MessageService) listenMessages(pubsub *redis.PubSub, roomID uint) {
	defer func() {
		pubsub.Close()
		s.activeSubs.Delete(roomID)
		log.Printf("取消订阅房间%d", roomID)
	}()

	ch := pubsub.Channel() //获取消息通道
	for msg := range ch {
		var message model.Message
		err := json.Unmarshal([]byte(msg.Payload), &message) //反序列化
		if err != nil {
			log.Printf("消息解析失败:%v", err.Error())
			continue
		}
		ws.HubInstance.Broadcast <- message //转发给hub
	}
}

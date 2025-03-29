package controllers

import (
	"chat/internal/handler/ws"
	"chat/internal/model"
	"chat/internal/service"
	myerrors "chat/pkg/errors"
	"chat/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ChatController struct {
	db             *gorm.DB
	messageService *service.MessageService
}

func NewChatController(db *gorm.DB, ms *service.MessageService) *ChatController {
	return &ChatController{
		db:             db,
		messageService: ms,
	}
}

func (c *ChatController) ConnectWebSocket(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.Param("room_id"))
	userID := ctx.MustGet("userid").(uint)

	//协议升级
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Fail(myerrors.ErrCodeUpdateFail, "协议升级失败"))
		return
	}

	//客户端注册
	client := &ws.Client{
		Conn:   conn,
		UserID: userID,
		RoomID: uint(roomID),
	}
	ws.HubInstance.Register <- client

	go c.messageService.SubscribeMessages(uint(roomID))

	//连接维持协程
	go func() {
		defer func() {
			ws.HubInstance.Unregister <- client
			_ = conn.Close()
		}()
		for {
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				break
			}
			var msg model.Message
			json.Unmarshal(msgBytes, &msg)
			msg.UserID = userID
			msg.RoomID = uint(roomID)
			msg.CreatedAt = time.Now()

			//msg=model.Message{
			//	UserID: userID,
			//	RoomID: uint(roomID),
			//	Content: string(msgBytes),
			//	CreatedAt: time.Now(),
			//}
			if err := c.db.Create(&msg).Error; err != nil {
				log.Println(err.Error())
			}
			_ = c.messageService.PublishMessage(uint(roomID), msg)

		}
	}()

}

func (c *ChatController) SendMessage(ctx *gin.Context) {
	var req struct {
		RoomID  uint   `json:"room_id" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Fail(myerrors.ErrCodeParamInvalid, "参数错误"))
		return
	}
	message := model.Message{
		UserID:    ctx.MustGet("userid").(uint),
		RoomID:    req.RoomID,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}
	if err := c.db.Create(&message).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Fail(myerrors.ErrCodeDatabaseError, "消息存入数据库失败"))
		return
	}
	_ = c.messageService.PublishMessage(req.RoomID, message)
	ctx.JSON(http.StatusOK, response.Success(message))
}

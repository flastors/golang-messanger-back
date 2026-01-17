package app

import (
	psqlChatRepo "GolangMessanger/internal/adapters/ChatRepository/postgresql"
	psqlMsgRepo "GolangMessanger/internal/adapters/MessageRepository/postgresql"
	"GolangMessanger/internal/config"
	"GolangMessanger/internal/controllers/restapi"
	"GolangMessanger/internal/domain/chat"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Context struct {
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) Run() error {
	cfg := config.Get()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.ChatService.ChatDB.User,
		cfg.ChatService.ChatDB.Password,
		cfg.ChatService.ChatDB.Host,
		cfg.ChatService.ChatDB.Port,
		cfg.ChatService.ChatDB.DBName,
		cfg.ChatService.ChatDB.SSLMode,
	)
	migrationErr := RunPostgresMigrations(dsn)
	if migrationErr != nil {
		return migrationErr
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	chatRepo := psqlChatRepo.NewRepository(db)
	messageRepo := psqlMsgRepo.NewRepository(db)
	chatUsecase := chat.NewChatUsecase(chatRepo, messageRepo, cfg.ChatService)
	usecases := restapi.Usecases{
		ChatUsecase: chatUsecase,
	}
	restapiServer := restapi.NewServer(usecases, cfg.Server.RestAPI)
	return restapiServer.Start()
}

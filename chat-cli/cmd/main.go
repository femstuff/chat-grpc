package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"chat-grpc/pkg/config"
	"chat-grpc/proto_gen"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	authClient   proto_gen.AuthServiceClient
	chatClient   proto_gen.ChatServiceClient
	refreshToken string
	accessToken  string
	mu           sync.Mutex
	log          *zap.Logger
)

func main() {
	cfg := config.LoadConfig()
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	log = logger
	defer log.Sync()

	authConn, err := grpc.NewClient("localhost:"+cfg.ServerPortAuth, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to auth service", zap.Error(err))
	}
	defer authConn.Close()
	authClient = proto_gen.NewAuthServiceClient(authConn)

	chatConn, err := grpc.NewClient("localhost:"+cfg.ServerPortChat, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to chat service", zap.Error(err))
	}
	defer chatConn.Close()
	chatClient = proto_gen.NewChatServiceClient(chatConn)

	fmt.Println("Добро пожаловать в gRPC-чат!")
	fmt.Println("Команды:")
	fmt.Println("  register <name> <email> <password> <role> - Создать пользователя")
	fmt.Println("  login <username> <password> - Войти в систему")
	fmt.Println("  get_access - Получить access token")
	fmt.Println("  create_chat <user1,user2,...> - Создать чат")
	fmt.Println("  send_message <chat_id> <from> <text> - Отправить сообщение")
	fmt.Println("  connect <chat_id> - Подключиться к чату")
	fmt.Println("  exit - Выйти")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		args := strings.Fields(input)

		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "register":
			if len(args) < 5 {
				fmt.Println("Формат: register <name> <email> <password> <role>")
				continue
			}
			err := registerUser(args[1], args[2], args[3], args[4])
			if err != nil {
				log.Error("Registration failed", zap.Error(err))
			}
		case "login":
			if len(args) < 3 {
				fmt.Println("Формат: login <email> <password>")
				continue
			}
			err := login(args[1], args[2])
			if err != nil {
				log.Error("Login failed", zap.Error(err))
			}

		case "get_access":
			err := getAccess()
			if err != nil {
				log.Error("Failed to get access token", zap.Error(err))
			}

		case "create_chat":
			if len(args) < 2 {
				fmt.Println("Формат: create_chat <user1,user2,...>")
				continue
			}
			users := strings.Split(args[1], ",")
			err := createChat(users)
			if err != nil {
				log.Error("Failed to create chat", zap.Error(err))
			}

		case "send_message":
			if len(args) < 4 {
				fmt.Println("Формат: send_message <chat_id> <from> <text>")
				continue
			}
			chatID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				log.Warn("Invalid chat ID", zap.String("input", args[1]))
				continue
			}
			err = sendMessage(chatID, args[2], strings.Join(args[3:], " "))
			if err != nil {
				log.Error("Failed to send message", zap.Error(err))
			}

		case "connect":
			if len(args) < 2 {
				fmt.Println("Формат: connect <chat_id>")
				continue
			}
			chatID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				log.Warn("Invalid chat ID", zap.String("input", args[1]))
				continue
			}
			connectToChat(chatID)

		case "exit":
			log.Info("Exiting CLI")
			return

		default:
			log.Info("Unknown command")
		}
	}
}

func registerUser(name, email, password, roleStr string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	role, err := parseRole(roleStr)
	if err != nil {
		log.Warn("Invalid role", zap.String("input", roleStr), zap.Error(err))
		return fmt.Errorf("ошибка указания роли: %w", err)
	}

	resp, err := authClient.Create(ctx, &proto_gen.CreateUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     role,
	})
	if err != nil {
		return fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	log.Info("User registered", zap.Int64("user_id", resp.Id))
	fmt.Println("Пользователь создан, ID:", resp.Id)
	return nil
}

func login(email, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := authClient.Login(ctx, &proto_gen.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("ошибка авторизации: %w", err)
	}

	setRefreshToken(resp.RefreshToken)
	log.Info("Login successful, refresh token received")
	fmt.Println("Успешный вход! Используйте 'get_access' для получения access token.")
	return nil
}

func getAccess() error {
	rt := getRefreshToken()
	if rt == "" {
		return fmt.Errorf("необходимо авторизоваться")
	}

	log.Info("Requesting new access token")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := authClient.GetAccessToken(ctx, &proto_gen.AccessTokenRequest{
		RefreshToken: rt,
	})
	if err != nil {
		return fmt.Errorf("ошибка получения access token: %w", err)
	}

	setAccessToken(resp.AccessToken)
	log.Info("Access token updated")
	fmt.Println("Access token обновлен!")
	return nil
}

func createChat(users []string) error {
	ctx := authContext()
	if ctx == nil {
		return fmt.Errorf("необходимо получить access_token")
	}

	log.Info("Creating chat", zap.Strings("users", users))

	resp, err := chatClient.Create(ctx, &proto_gen.CreateRequest{Usernames: users})
	if err != nil {
		return fmt.Errorf("ошибка создания чата: %w", err)
	}

	log.Info("Chat created", zap.Int64("chat_id", resp.Id))
	fmt.Println("Чат создан, ID:", resp.Id)
	return nil
}

func sendMessage(chatID int64, from, text string) error {
	ctx := authContext()
	if ctx == nil {
		return fmt.Errorf("необходимо получить access_token")
	}

	log.Info("Sending message", zap.Int64("chat_id", chatID), zap.String("from", from))

	_, err := chatClient.SendMessage(ctx, &proto_gen.SendMessageRequest{
		ChatId: chatID,
		From:   from,
		Text:   text,
	})
	if err != nil {
		return fmt.Errorf("ошибка отправки сообщения: %w", err)
	}

	log.Info("Message sent")
	fmt.Println("Сообщение отправлено")
	return nil
}

func connectToChat(chatID int64) {
	ctx := authContext()
	if ctx == nil {
		log.Warn("Access token is missing")
		fmt.Println("Необходимо получить access_token")
		return
	}

	log.Info("Connecting to chat", zap.Int64("chat_id", chatID))

	historyResp, err := chatClient.GetMessages(ctx, &proto_gen.GetMessagesRequest{ChatId: chatID})
	if err != nil {
		log.Error("Failed to load chat history", zap.Error(err))
		fmt.Println("Ошибка загрузки истории:", err)

		return
	}

	fmt.Println("История чата:")
	for _, msg := range historyResp.Messages {
		fmt.Printf("[%s] %s: %s\n", msg.Timestamp.AsTime().Format("15:04"), msg.From, msg.Text)
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Error("Failed to connect to NATS", zap.Error(err))
		fmt.Println("Ошибка подключения к NATS:", err)
		return
	}

	subject := fmt.Sprintf("chat.%d", chatID)
	_, err = nc.Subscribe(subject, func(msg *nats.Msg) {
		var m proto_gen.Message
		if err := protojson.Unmarshal(msg.Data, &m); err != nil {
			log.Error("Failed to parse protojson message", zap.Error(err))
			return
		}

		log.Info("incoming message (NATS)",
			zap.String("from", m.From),
			zap.String("text", m.Text))
	})
	if err != nil {
		log.Error("Failed to subscribe to NATS", zap.Error(err))
		return
	}

	stream, err := chatClient.Connect(ctx, &proto_gen.ConnectRequest{ChatId: chatID})
	if err != nil {
		log.Error("Failed to connect to chat stream", zap.Error(err))
		fmt.Println("Ошибка подключения к чату:", err)
		return
	}

	log.Info("Connected to chat stream", zap.Int64("chat_id", chatID))
	fmt.Println("Подключен к чату. Ожидание сообщений...")

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Error("Error receiving message", zap.Error(err))
				fmt.Println("Ошибка при получении сообщения:", err)
				return
			}

			fmt.Printf("[%s] %s: %s\n", msg.Timestamp.AsTime().Format("15:04"), msg.From, msg.Text)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if scanner.Text() == "exit" {
				log.Info("User exited chat", zap.Int64("chat_id", chatID))
				return
			}
		}
	}()
}

func setRefreshToken(token string) {
	mu.Lock()
	defer mu.Unlock()
	refreshToken = token
}

func getRefreshToken() string {
	mu.Lock()
	defer mu.Unlock()
	return refreshToken
}

func setAccessToken(token string) {
	mu.Lock()
	defer mu.Unlock()
	accessToken = token
}

func getAccessToken() string {
	mu.Lock()
	defer mu.Unlock()
	return accessToken
}

func authContext() context.Context {
	token := getAccessToken()
	if token == "" {
		return nil
	}
	return metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))
}

func parseRole(roleStr string) (proto_gen.Role, error) {
	switch strings.ToLower(roleStr) {
	case "admin":
		return proto_gen.Role_AdminRole, nil
	case "user":
		return proto_gen.Role_UserRole, nil
	default:
		return proto_gen.Role_UserRole, fmt.Errorf("unknown role: %s", roleStr)
	}
}

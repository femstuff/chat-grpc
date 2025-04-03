package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"chat-grpc/proto_gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	authClient proto_gen.AuthServiceClient
	chatClient proto_gen.ChatServiceClient
	token      string
)

func main() {
	authConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Ошибка подключения к auth-сервису:", err)
		return
	}
	defer authConn.Close()
	authClient = proto_gen.NewAuthServiceClient(authConn)

	chatConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Ошибка подключения к chat-сервису:", err)
		return
	}
	defer chatConn.Close()
	chatClient = proto_gen.NewChatServiceClient(chatConn)

	fmt.Println("Добро пожаловать в gRPC-чат!")
	fmt.Println("Команды:")
	fmt.Println("  login <username> <password> - Войти в систему")
	fmt.Println("  create_chat <user1,user2,...> - Создать чат")
	fmt.Println("  send_message <chat_id> <from> <text> - Отправить сообщение")
	fmt.Println("  connect <chat_id> - Подключиться к чату")
	fmt.Println("  exit - Выйти")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		args := strings.Fields(input)

		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "login":
			if len(args) < 3 {
				fmt.Println("Формат: login <username> <password>")
				continue
			}
			err := login(args[1], args[2])
			if err != nil {
				fmt.Println("Ошибка:", err)
			}

		case "create_chat":
			if len(args) < 2 {
				fmt.Println("Формат: create_chat <user1,user2,...>")
				continue
			}
			users := strings.Split(args[1], ",")
			err := createChat(users)
			if err != nil {
				fmt.Println("Ошибка:", err)
			}

		case "send_message":
			if len(args) < 4 {
				fmt.Println("Формат: send_message <chat_id> <from> <text>")
				continue
			}
			chatID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				fmt.Println("Неверный ID чата")
				continue
			}
			err = sendMessage(chatID, args[2], strings.Join(args[3:], " "))
			if err != nil {
				fmt.Println("Ошибка:", err)
			}

		case "connect":
			if len(args) < 2 {
				fmt.Println("Формат: connect <chat_id>")
				continue
			}
			chatID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				fmt.Println("Неверный ID чата")
				continue
			}
			connectToChat(chatID)

		case "exit":
			fmt.Println("Выход из программы")
			return

		default:
			fmt.Println("Неизвестная команда")
		}
	}
}

func login(username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := authClient.Login(ctx, &proto_gen.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("ошибка авторизации: %w", err)
	}

	token = resp.RefreshToken
	fmt.Println("Успешный вход! Refresh Token:", token)
	return nil
}

func createChat(users []string) error {
	if token == "" {
		return fmt.Errorf("необходимо авторизоваться")
	}

	ctx := context.Background()
	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := chatClient.Create(ctx, &proto_gen.CreateRequest{Usernames: users})
	if err != nil {
		return fmt.Errorf("ошибка создания чата: %w", err)
	}

	fmt.Println("Чат создан, ID:", resp.Id)
	return nil
}

func sendMessage(chatID int64, from, text string) error {
	if token == "" {
		return fmt.Errorf("необходимо авторизоваться")
	}

	ctx := context.Background()
	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err := chatClient.SendMessage(ctx, &proto_gen.SendMessageRequest{
		ChatId: chatID,
		From:   from,
		Text:   text,
	})
	if err != nil {
		return fmt.Errorf("ошибка отправки сообщения: %w", err)
	}

	fmt.Println("Сообщение отправлено")
	return nil
}

func connectToChat(chatID int64) error {
	if token == "" {
		return fmt.Errorf("необходимо авторизоваться")
	}

	ctx := context.Background()
	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := chatClient.Connect(ctx, &proto_gen.ConnectRequest{ChatId: chatID})
	if err != nil {
		return fmt.Errorf("ошибка подключения к чату: %w", err)
	}

	fmt.Println("Подключен к чату. Ожидание сообщений...")

	for {
		msg, err := stream.Recv()
		if err != nil {
			fmt.Println("Ошибка при получении сообщения:", err)
			return err
		}

		fmt.Printf("[%s] %s: %s\n", msg.Timestamp.AsTime().Format(time.RFC822), msg.From, msg.Text)
	}
}

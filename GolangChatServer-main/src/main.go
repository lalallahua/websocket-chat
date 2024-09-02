package main

// 导入必要的包
import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// 初始化客户端和广播频道
var clients = make(map[*websocket.Conn]bool) // 存储所有连接的客户端
var broadcast = make(chan Message)           // 广播频道，用于向所有客户端发送消息

// 升级器，用于将普通HTTP连接升级为WebSocket连接
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源的连接（在生产环境中应进行更严格的检查）
	},
}

// 定义消息对象
type Message struct {
	Email    string `json:"email"`    // 发送者的电子邮件
	Username string `json:"username"` // 发送者的用户名
	Message  string `json:"message"`  // 发送的消息内容
}

// 主函数，Go程序的入口点
func main() {
	// 创建文件服务器，使用户可以访问前端文件（app.js和style.css）
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// 处理WebSocket连接请求
	http.HandleFunc("/ws", handleConnections)

	// 异步调用函数以处理所有消息
	go handleMessages()

	// 启动服务器并监听8000端口
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err) // 记录错误并退出程序
	}
}

// 处理WebSocket连接的函数
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// 将HTTP请求升级为WebSocket连接
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 函数返回时关闭连接
	defer ws.Close()

	// 将新的客户端添加到客户端字典中
	clients[ws] = true

	// 循环监听任何消息，并将其发送到广播频道
	for {
		var msg Message
		// 从WebSocket连接中读取JSON消息
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws) // 出错时将客户端从字典中删除
			break               // 退出循环
		}
		// 将消息发送到广播频道
		broadcast <- msg
	}
}

// 从广播频道中获取消息并将其转发给所有客户端
func handleMessages() {
	for {
		// 从广播频道中获取下一条消息
		msg := <-broadcast
		// 将消息发送给当前所有连接的客户端
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()          // 关闭出错的客户端连接
				delete(clients, client) // 将出错的客户端从字典中删除
			}
		}
	}
}

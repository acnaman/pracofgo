package main
import(
  "github.com/gorilla/websocket"
)
// clientはチャットを行っている1人のユーザ
type client struct{
  // socketはこのクライアントのためのwebsocketです
  socket *websocket.Conn
  // sendはメッセージが送られるチャネルです
  send chan []byte
  // roomはコオクライアントが参加しているチャットルーム
  room *room
}

package server

type Server struct {
	idGenerator      func() int
	connectedClients map[int]*ClientConnection
}

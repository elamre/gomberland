package packet

type ConnectionAction int

type ConnectionPacket struct {
	UserId  int
	Action  ConnectionAction
	Message string
}

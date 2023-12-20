package services

import "encoding/json"

const (
	EchoMessageType   = "echo"
	EchoOkMessageType = "echo_ok"
)

type EchoMessage struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body struct {
		Type  string `json:"type"`
		MsgId int32  `json:"msg_id"`
		Echo  string `json:"echo"`
	} `json:"body"`
}

type EchoOkMessage struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body struct {
		Type      string `json:"type"`
		MsgId     int32  `json:"msg_id"`
		InReplyTo int32  `json:"in_reply_to"`
		Echo      string `json:"echo"`
	} `json:"body"`
}

func (n *Node) HandleEcho(message []byte) ([]byte, error) {
	echo := EchoMessage{}
	err := json.Unmarshal(message, &echo)
	if err != nil {
		return nil, err
	}
	echoOk := EchoOkMessage{}
	echoOk.Src, echoOk.Dest = echo.Dest, echo.Src
	echoOk.Body.Type = EchoOkMessageType
	echoOk.Body.MsgId = n.GetNextMessageId()
	echoOk.Body.InReplyTo = echo.Body.MsgId
	echoOk.Body.Echo = echo.Body.Echo

	response, err := json.Marshal(echoOk)
	return response, err
}

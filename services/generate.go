package services

import (
	"encoding/json"

	"github.com/google/uuid"
)

const (
	GenerateMessageType   = "generate"
	GenerateOkMessageType = "generate_ok"
)

type GenerateMessage struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body struct {
		Type  string `json:"type"`
		MsgId int32  `json:"msg_id"`
	} `json:"body"`
}

type GenerateOkMessage struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body struct {
		Type      string `json:"type"`
		MsgId     int32  `json:"msg_id"`
		InReplyTo int32  `json:"in_reply_to"`
		Id        string `json:"id"`
	} `json:"body"`
}

func (n *Node) HandleGenerate(message []byte) ([]byte, error) {
	generate := GenerateMessage{}
	err := json.Unmarshal(message, &generate)
	if err != nil {
		return nil, err
	}

	generateOk := GenerateOkMessage{}
	generateOk.Dest, generateOk.Src = generate.Src, generate.Dest
	generateOk.Body.Type = GenerateOkMessageType
	generateOk.Body.MsgId = n.GetNextMessageId()
	generateOk.Body.InReplyTo = generate.Body.MsgId
	generateOk.Body.Id = uuid.New().String()

	response, err := json.Marshal(generateOk)
	return response, err
}

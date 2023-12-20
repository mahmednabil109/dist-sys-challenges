package services

import "encoding/json"

const (
	InitMessageType   = "init"
	InitOkMessageType = "init_ok"
)

type InitMessage struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body struct {
		Type    string   `json:"type"`
		MsgId   int32    `json:"msg_id"`
		NodeId  string   `json:"node_id"`
		NodeIds []string `json:"node_ids"`
	} `json:"body"`
}

type InitOkMessage struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body struct {
		Type      string `json:"type"`
		InReplyTo int32  `json:"in_reply_to"`
	} `json:"body"`
}

func (n *Node) HandleInit(message []byte) ([]byte, error) {
	init := InitMessage{}
	err := json.Unmarshal(message, &init)
	if err != nil {
		return nil, err
	}

	n.NodeId = init.Body.NodeId
	n.NextNodeIds = init.Body.NodeIds

	initOk := InitOkMessage{}
	initOk.Src, initOk.Dest = init.Dest, init.Src
	initOk.Body.Type = InitOkMessageType
	initOk.Body.InReplyTo = init.Body.MsgId

	response, err := json.Marshal(initOk)
	return response, err
}

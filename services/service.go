package services

import (
	"bytes"
	"errors"
)

type NodeHandleFunc func([]byte) ([]byte, error)

func Nop([]byte) ([]byte, error) {
	return nil, nil
}

type Node struct {
	NodeId      string
	NextNodeIds []string

	SendBuffer [][]byte

	nextMessageId int32
	msgTypeTabel  map[string]NodeHandleFunc
	store         map[int32]bool
}

func NewNode() (node *Node) {
	node = &Node{
		nextMessageId: 1,
		store:         make(map[int32]bool),
	}
	node.msgTypeTabel = map[string]NodeHandleFunc{
		InitMessageType:        node.HandleInit,
		EchoMessageType:        node.HandleEcho,
		GenerateMessageType:    node.HandleGenerate,
		BroadcastMessageType:   node.HandleBroadcast,
		BroadcastOkMessageType: Nop,
		ReadMessageType:        node.HandleRead,
		TopologyMessageType:    node.HandleTopology,
	}
	return
}

func (n *Node) GetNextMessageId() (id int32) {
	id = n.nextMessageId
	n.nextMessageId++
	return
}

func (n *Node) Send(message []byte) {
	n.SendBuffer = append(n.SendBuffer, message)
}

func (n *Node) HandleMessage(message []byte) ([]byte, error) {
	idx := bytes.Index(message, []byte("type"))
	if idx == -1 {
		return nil, errors.New("can't interpert message with no type")
	}
	eidx := bytes.IndexByte(message[idx:], ',')
	if eidx == -1 {
		eidx = len(message) - idx - 2
	}

	messageType := string(message[idx+7 : idx+eidx-1])
	nodeFunc, ok := n.msgTypeTabel[messageType]
	if !ok {
		return nil, errors.New("can't handle message with unkown type: " + messageType)
	}
	response, err := nodeFunc(message)
	return response, err
}

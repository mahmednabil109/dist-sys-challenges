package node

import (
	"encoding/json"
)

const (
	BroadcastMessageType   = "broadcast"
	BroadcastOkMessageType = "broadcast_ok"
	ReadMessageType        = "read"
	ReadOkMessageType      = "read_ok"
	TopologyMessageType    = "topology"
	TopologyOkMessageType  = "topology_ok"
)

type BroadcastMessage struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body struct {
		Type    string `json:"type"`
		Message int32  `json:"message"`
		MsgId   int32  `json:"msg_id"`
	} `json:"body"`
}

type BroadcastOkMessage struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body struct {
		Type      string `json:"type"`
		InReplyTo int32  `json:"in_reply_to"`
	} `json:"body"`
}

type ReadMessage struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body struct {
		Type  string `json:"type"`
		MsgId int32  `json:"msg_id"`
	} `json:"body"`
}

type ReadOkMessage struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body struct {
		Type      string  `json:"type"`
		Messages  []int32 `json:"messages"`
		InReplyTo int32   `json:"in_reply_to"`
	} `json:"body"`
}

type TopologyMessage struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body struct {
		Type     string              `json:"type"`
		Topology map[string][]string `json:"topology"`
		MsgId    int32               `json:"msg_id"`
	} `json:"body"`
}

type TopologyOkMessage struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body struct {
		Type      string `json:"type"`
		InReplyTo int32  `json:"in_reply_to"`
	} `json:"body"`
}

func (n *Node) HandleBroadcast(message []byte) ([]byte, error) {
	broadcast := BroadcastMessage{}
	if err := json.Unmarshal(message, &broadcast); err != nil {
		return nil, err
	}

	// TODO optimize the gossiping
	if _, ok := n.store[broadcast.Body.Message]; !ok {
		n.store[broadcast.Body.Message] = true
		for _, nextNodeId := range n.NextNodeIds {
			if nextNodeId == broadcast.Src || nextNodeId == n.NodeId {
				continue
			}
			newBroadcast := BroadcastMessage{
				Src:  n.NodeId,
				Dest: nextNodeId,
			}
			newBroadcast.Body.Type = BroadcastMessageType
			newBroadcast.Body.Message = broadcast.Body.Message
			newBroadcast.Body.MsgId = n.GetNextMessageId()
			msg, _ := json.Marshal(newBroadcast)
			n.Send(msg)
			n.awaitAck(newBroadcast.Body.MsgId, msg)
		}
	}

	broadcastOk := BroadcastOkMessage{}
	broadcastOk.Src, broadcastOk.Dest = broadcast.Dest, broadcast.Src
	broadcastOk.Body.Type = BroadcastOkMessageType
	broadcastOk.Body.InReplyTo = broadcast.Body.MsgId

	response, err := json.Marshal(broadcastOk)
	return response, err
}

func (n *Node) HandleRead(message []byte) ([]byte, error) {
	read := ReadMessage{}
	if err := json.Unmarshal(message, &read); err != nil {
		return nil, err
	}

	readOk := ReadOkMessage{}
	readOk.Src, readOk.Dest = read.Dest, read.Src
	readOk.Body.Type = ReadOkMessageType
	messages := make([]int32, 0, len(n.store))
	for msg := range n.store {
		messages = append(messages, msg)
	}
	readOk.Body.Messages = messages
	readOk.Body.InReplyTo = read.Body.MsgId

	response, err := json.Marshal(readOk)
	return response, err
}

func (n *Node) HandleTopology(message []byte) ([]byte, error) {
	topology := TopologyMessage{}
	if err := json.Unmarshal(message, &topology); err != nil {
		return nil, err
	}

	// TODO save the topology some where to perform smart gossiping

	topologyOk := TopologyOkMessage{}
	topologyOk.Src, topologyOk.Dest = topology.Dest, topology.Src
	topologyOk.Body.Type = TopologyOkMessageType
	topologyOk.Body.InReplyTo = topology.Body.MsgId

	response, err := json.Marshal(topologyOk)
	return response, err
}

func (n *Node) HandleBroadcastOk(message []byte) ([]byte, error) {
	broadcastOk := BroadcastOkMessage{}
	if err := json.Unmarshal(message, &broadcastOk); err != nil {
		return nil, err
	}
	n.ackMsg(broadcastOk.Body.InReplyTo)
	return nil, nil
}

package node

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type NodeHandleFunc func([]byte) ([]byte, error)

func Nop([]byte) ([]byte, error) {
	return nil, nil
}

type sentMessage struct {
	time  time.Time
	msgId int32
	msg   []byte

	isAcked bool
}

type Node struct {
	NodeId      string
	NextNodeIds []string

	nextMessageId int32
	msgTypeTabel  map[string]NodeHandleFunc
	store         map[int32]bool

	stream         io.Writer
	send           chan []byte
	awaitAckMsg    []sentMessage
	awaitAckMsgMux sync.Mutex
}

func New(outputStream io.Writer) (node *Node) {
	node = &Node{
		nextMessageId: 1,
		store:         make(map[int32]bool),
		send:          make(chan []byte),
		stream:        outputStream,
	}
	node.msgTypeTabel = map[string]NodeHandleFunc{
		InitMessageType:        node.HandleInit,
		EchoMessageType:        node.HandleEcho,
		GenerateMessageType:    node.HandleGenerate,
		BroadcastMessageType:   node.HandleBroadcast,
		BroadcastOkMessageType: node.HandleBroadcastOk,
		ReadMessageType:        node.HandleRead,
		TopologyMessageType:    node.HandleTopology,
	}

	go node.writer()
	go node.retry()

	return
}

func (n *Node) GetNextMessageId() (id int32) {
	id = n.nextMessageId
	n.nextMessageId++
	return
}

func (n *Node) Send(message []byte) {
	go func() {
		n.send <- message
	}()
}

func (n *Node) writer() {
	for msg := range n.send {
		if msg == nil {
			continue
		}
		fmt.Fprintln(os.Stderr, string(msg))
		_, err := fmt.Fprintln(n.stream, string(msg))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Err occured when sending msg: %v", err)
		}
	}
}

func (n *Node) ackMsg(msgId int32) {
	n.awaitAckMsgMux.Lock()
	defer n.awaitAckMsgMux.Unlock()
	for _, m := range n.awaitAckMsg {
		if m.msgId == msgId {
			m.isAcked = true
			break
		}
	}
}

func (n *Node) awaitAck(msgId int32, msg []byte) {
	n.awaitAckMsgMux.Lock()
	defer n.awaitAckMsgMux.Unlock()
	n.awaitAckMsg = append(n.awaitAckMsg, sentMessage{
		time:  time.Now(),
		msgId: msgId,
		msg:   msg,
	})
}

func (n *Node) retry() {
	timer := time.NewTicker(5 * time.Second)
	for range timer.C {
		n.awaitAckMsgMux.Lock()
		newList := make([]sentMessage, len(n.awaitAckMsg))
		for _, m := range n.awaitAckMsg {
			if time.Since(m.time) > time.Microsecond*10 {
				n.Send(m.msg)
				m.time = time.Now()
			}
			if !m.isAcked {
				newList = append(newList, m)
			}
		}
		n.awaitAckMsg = newList
		n.awaitAckMsgMux.Unlock()
	}
}

func (n *Node) HandleMessage(message []byte) {
	// get message type from the byte stream
	idx := bytes.Index(message, []byte("type"))
	if idx == -1 {
		fmt.Fprintln(os.Stderr, "can't interpert message with no type")
	}
	eidx := bytes.IndexByte(message[idx:], ',')
	if eidx == -1 {
		eidx = len(message) - idx - 2
	}
	messageType := string(message[idx+7 : idx+eidx-1])

	// dispatch corresponding function
	nodeFunc, ok := n.msgTypeTabel[messageType]
	if !ok {
		fmt.Fprintln(os.Stderr, "can't handle message with unkown type: "+messageType)
	}
	response, err := nodeFunc(message)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	n.Send(response)
}

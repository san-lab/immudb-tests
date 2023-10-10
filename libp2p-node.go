package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

var NET string

// MPCNode represents a subscription to a single PubSub topic. Messages
// can be published to the topic with MPCNode.Publish, and received
// messages are pushed to the Messages channel.
type Node struct {
	// Messages is a channel of messages received from other peers in the chat room
	ctx   context.Context
	ps    *pubsub.PubSub
	topic *pubsub.Topic
	sub   *pubsub.Subscription
	self  peer.ID
}

// readLoop pulls messages from the pubsub topic and pushes them onto the Messages channel.
func (node *Node) readLoop() {
	for {
		msg, err := node.sub.Next(node.ctx)
		if err != nil {
			log.Fatalln(err)
			return
		}
		// only forward messages delivered by others
		if msg.ReceivedFrom == node.self {
			continue
		}

		node.ProcessMessage(msg)
	}
}

func (node *Node) ProcessMessage(msg *pubsub.Message) {
	txmsg := new(MT103Message)
	err := json.Unmarshal(msg.Data, txmsg)
	if err != nil {
		fmt.Println("bad frame:", err)
		return
	}
	ProcessInterBankTx(txmsg)
}

func (node *Node) SendMsg(txmsg *MT103Message) {
	b, _ := json.Marshal(txmsg)
	node.topic.Publish(context.Background(), b)
}

func (node *Node) PeerCount() int {
	mutex.Lock()
	defer mutex.Unlock()
	return len(node.topic.ListPeers())
}

// tries to subscribe to the PubSub topic for the room name, returning
// an MPCNode on success.
func JoinNet(ctx context.Context, ps *pubsub.PubSub, selfID peer.ID, roomName string) (*Node, error) {
	// join the pubsub topic
	mutex.Lock()
	defer mutex.Unlock()
	topic, err := ps.Join(NET)
	if err != nil {
		return nil, err
	}

	// and subscribe to it
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	node := &Node{
		ctx:   ctx,
		ps:    ps,
		topic: topic,
		sub:   sub,
		self:  selfID,
	}

	// start reading messages from the subscription in a loop
	go node.readLoop()
	return node, nil
}

func (node *Node) GetNodeID() string {
	return string(node.self)
}

package transactions

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	libp2pconf "github.com/libp2p/go-libp2p/config"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// DiscoveryInterval is how often we re-publish our mDNS records.
const DiscoveryInterval = time.Hour

// DiscoveryServiceTag is used in our mDNS advertisements to discover other chat peers.
const DiscoveryServiceTag = "ImmuDBNet"

// LibP2P Host
var h host.Host
var mutex = &sync.Mutex{}

func GetHost() (host.Host, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if h == nil {
		kopt := func(c *libp2pconf.Config) error {
			c.PeerKey = nil
			return nil
		}
		nh, err := libp2p.New(kopt, libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))

		if err != nil {
			return nil, err
		}
		// setup local mDNS discovery
		if err = setupDiscovery(nh); err != nil {
			return nil, err
		}
		h = nh
	}
	return h, nil
}

// MPC Node
var node *Node

func GetNode() (*Node, error) {
	if node == nil {
		h, err := GetHost()
		if err != nil {
			return nil, err
		}
		// create a new PubSub service using the GossipSub router
		ctx := context.Background()
		ps, err := pubsub.NewGossipSub(ctx, h)
		if err != nil {
			return nil, err
		}

		node, err = JoinNet(ctx, ps, h.ID(), NET)
		if err != nil {
			return nil, err
		}
	}

	return node, nil
}

//***************

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
	h host.Host
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	//fmt.Printf("discovered new peer %s\n", pi.ID.Pretty())
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", pi.ID, err)
	}
}

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func setupDiscovery(h host.Host) error {
	// setup mDNS discovery to find local peers
	s := mdns.NewMdnsService(h, DiscoveryServiceTag, &discoveryNotifee{h: h})
	return s.Start()
}

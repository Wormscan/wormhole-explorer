package main

import (
	"context"
	"fmt"
	"os"

	"github.com/certusone/wormhole/node/pkg/common"
	"github.com/certusone/wormhole/node/pkg/p2p"
	gossipv1 "github.com/certusone/wormhole/node/pkg/proto/gossip/v1"
	"github.com/certusone/wormhole/node/pkg/supervisor"
	eth_common "github.com/ethereum/go-ethereum/common"
	ipfslog "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/crypto"
	"go.uber.org/zap"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	rootCtx       context.Context
	rootCtxCancel context.CancelFunc
)

var (
	p2pNetworkID string
	p2pPort      uint
	p2pBootstrap string
	nodeKeyPath string
	logLevel string
)

func main() {
	p2pNetworkID = "/wormhole/mainnet/2"
	p2pPort = 8999
	p2pBootstrap = "/dns4/wormhole-mainnet-v2-bootstrap.certus.one/udp/8999/quic/p2p/12D3KooWQp644DK27fd3d4Km3jr7gHiuJJ5ZGmy8hH4py7fP4FP7"
	nodeKeyPath = "/tmp/node.key"
	logLevel = "info"
	common.SetRestrictiveUmask()

	lvl, err := ipfslog.LevelFromString(logLevel)
	if err != nil {
		fmt.Println("Invalid log level")
		os.Exit(1)
	}

	logger := ipfslog.Logger("wormhole-spy").Desugar()

	ipfslog.SetAllLoggers(lvl)

	// Verify flags
	if nodeKeyPath == "" {
		logger.Fatal("Please specify --nodeKey")
	}
	if p2pBootstrap == "" {
		logger.Fatal("Please specify --bootstrap")
	}

	// Setup DB
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		logger.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	hbColl := client.Database("wormhole").Collection("heartbeats")
	obsColl := client.Database("wormhole").Collection("observations")
	vaaColl := client.Database("wormhole").Collection("vaas")

	// Node's main lifecycle context.
	rootCtx, rootCtxCancel = context.WithCancel(context.Background())
	defer rootCtxCancel()

	// Outbound gossip message queue
	sendC := make(chan []byte)

	// Inbound observations
	obsvC := make(chan *gossipv1.SignedObservation, 50)

	// Inbound signed VAAs
	signedInC := make(chan *gossipv1.SignedVAAWithQuorum, 50)

	// Heartbeat updates
	heartbeatC := make(chan *gossipv1.Heartbeat, 50)

	// Guardian set state managed by processor
	gst := common.NewGuardianSetState(heartbeatC)

	// Bootstrap guardian set, otherwise heartbeats would be skipped
	gst.Set(&common.GuardianSet{
		Index: 2,
		Keys: []eth_common.Address{
			eth_common.HexToAddress("0x58CC3AE5C097b213cE3c81979e1B9f9570746AA5"),
			eth_common.HexToAddress("0xfF6CB952589BDE862c25Ef4392132fb9D4A42157"),
			eth_common.HexToAddress("0x114De8460193bdf3A2fCf81f86a09765F4762fD1"),
			eth_common.HexToAddress("0x107A0086b32d7A0977926A205131d8731D39cbEB"),
			eth_common.HexToAddress("0x8C82B2fd82FaeD2711d59AF0F2499D16e726f6b2"),
			eth_common.HexToAddress("0x11b39756C042441BE6D8650b69b54EbE715E2343"),
			eth_common.HexToAddress("0x54Ce5B4D348fb74B958e8966e2ec3dBd4958a7cd"),
			eth_common.HexToAddress("0x66B9590e1c41e0B226937bf9217D1d67Fd4E91F5"),
			eth_common.HexToAddress("0x74a3bf913953D695260D88BC1aA25A4eeE363ef0"),
			eth_common.HexToAddress("0x000aC0076727b35FBea2dAc28fEE5cCB0fEA768e"),
			eth_common.HexToAddress("0xAF45Ced136b9D9e24903464AE889F5C8a723FC14"),
			eth_common.HexToAddress("0xf93124b7c738843CBB89E864c862c38cddCccF95"),
			eth_common.HexToAddress("0xD2CC37A4dc036a8D232b48f62cDD4731412f4890"),
			eth_common.HexToAddress("0xDA798F6896A3331F64b48c12D1D57Fd9cbe70811"),
			eth_common.HexToAddress("0x71AA1BE1D36CaFE3867910F99C09e347899C19C3"),
			eth_common.HexToAddress("0x8192b6E7387CCd768277c17DAb1b7a5027c0b3Cf"),
			eth_common.HexToAddress("0x178e21ad2E77AE06711549CFBB1f9c7a9d8096e8"),
			eth_common.HexToAddress("0x5E1487F35515d02A92753504a8D75471b9f49EdB"),
			eth_common.HexToAddress("0x6FbEBc898F403E4773E95feB15E80C9A99c8348d"),
		},
	})

	// Ignore observations
	go func() {
		for {
			select {
			case <-rootCtx.Done():
				return
			case o := <- obsvC:
				logger.Info("Received observation", zap.Any("observation", o))
				result, err := obsColl.InsertOne(context.TODO(), o)
				if err != nil {
					logger.Error("Error inserting observation", zap.Error(err))
				}
				logger.Info("Inserted document", zap.Any("id", result.InsertedID))
			}
		}
	}()

	// Log signed VAAs
	go func() {
		for {
			select {
			case <-rootCtx.Done():
				return
			case v := <-signedInC:
				logger.Info("Received signed VAA",
					zap.Any("vaa", v.Vaa))
				result, err := vaaColl.InsertOne(context.TODO(), v)
				if err != nil {
					logger.Error("Error inserting vaa", zap.Error(err))
				}
				logger.Info("Inserted document", zap.Any("id", result.InsertedID))
			}
		}
	}()

	// Ignore heartbeats
	go func() {
		for {
			select {
			case <-rootCtx.Done():
				return
			case hb := <- heartbeatC:
				logger.Info("Received heartbeat", zap.Any("heartbeat", hb))
				result, err := hbColl.InsertOne(context.TODO(), hb)
				if err != nil {
					logger.Error("Error inserting heartbeat", zap.Error(err))
				}
				logger.Info("Inserted document", zap.Any("id", result.InsertedID))
			}
		}
	}()

	// Load p2p private key
	var priv crypto.PrivKey
	priv, err = common.GetOrCreateNodeKey(logger, nodeKeyPath)
	if err != nil {
		logger.Fatal("Failed to load node key", zap.Error(err))
	}

	// Run supervisor.
	supervisor.New(rootCtx, logger, func(ctx context.Context) error {
		if err := supervisor.Run(ctx, "p2p", p2p.Run(obsvC, nil, nil, sendC, signedInC, priv, nil, gst, p2pPort, p2pNetworkID, p2pBootstrap, "", false, rootCtxCancel, nil)); err != nil {
			return err
		}

		logger.Info("Started internal services")

		<-ctx.Done()
		return nil
	},
		// It's safer to crash and restart the process in case we encounter a panic,
		// rather than attempting to reschedule the runnable.
		supervisor.WithPropagatePanic)

	<-rootCtx.Done()
	logger.Info("root context cancelled, exiting...")
	// TODO: wait for things to shut down gracefully
}

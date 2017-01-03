package wormhole

import (
	"encoding/hex"
	"io"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/superfly/smux"
	kcp "github.com/xtaci/kcp-go"
)

// Constants, useful to both server and client
const (
	NoDelay      = 0
	Interval     = 30
	Resend       = 2
	NoCongestion = 1
	MaxBuffer    = 4194304
	KeepAlive    = 10

	// KCP
	KCPShards = 10
	KCPParity = 3
	DSCP      = 0

	SecretLength = 32
)

var (
	passphrase string
	publicKey  string
	version    string

	logLevel   = os.Getenv("LOG_LEVEL")
	smuxConfig *smux.Config
)

func init() {
	if logLevel == "" {
		log.SetLevel(log.InfoLevel)
	} else if logLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	}
	// logging
	textFormatter := &log.TextFormatter{FullTimestamp: true}
	log.SetFormatter(textFormatter)

	smuxConfig = smux.DefaultConfig()
	smuxConfig.MaxReceiveBuffer = MaxBuffer
	smuxConfig.KeepAliveInterval = KeepAlive * time.Second
}

func ensureEnvironment() {
	if passphrase == "" {
		passphrase = os.Getenv("PASSPHRASE")
		if passphrase == "" {
			log.Fatalln("PASSPHRASE needs to be set")
		} else if len([]byte(passphrase)) < SecretLength {
			log.Fatalf("PASSPHRASE needs to be at least %d bytes long\n", SecretLength)
		}
	}
	if publicKey == "" {
		publicKey = os.Getenv("PUBLIC_KEY")
		if publicKey == "" {
			log.Fatalln("PUBLIC_KEY needs to be set")
		}
	}
	publicKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		log.Fatalf("PUBLIC_KEY needs to be in hex format. Details: %s", err.Error())
	}
	if len(publicKeyBytes) != SecretLength {
		log.Fatalf("PUBLIC_KEY needs to be %d bytes long\n", SecretLength)
	}
	copy(smuxConfig.ServerPublicKey[:], publicKeyBytes)
	if version == "" {
		version = "latest"
	}
}

// OutputSNMP ...
func OutputSNMP() {
	log.Printf("KCP SNMP:%+v", kcp.DefaultSnmp.Copy())
}

// DebugSNMP ...
func DebugSNMP() {
	for _ = range time.Tick(120 * time.Second) {
		OutputSNMP()
	}
}

// CopyCloseIO ...
func CopyCloseIO(c1, c2 io.ReadWriteCloser) (err error) {
	defer c1.Close()
	defer c2.Close()

	// start tunnel
	c1die := make(chan struct{})
	go func() {
		_, err = io.Copy(c1, c2)
		close(c1die)
	}()

	c2die := make(chan struct{})
	go func() {
		_, err = io.Copy(c2, c1)
		close(c2die)
	}()

	// wait for tunnel termination
	select {
	case <-c1die:
	case <-c2die:
	}
	return err
}

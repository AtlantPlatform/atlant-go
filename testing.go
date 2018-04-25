//+build testing

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
	"sync"

	"github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
	"github.com/xlab/closer"

	"github.com/AtlantPlatform/atlant-go/eth"
	"github.com/AtlantPlatform/atlant-go/eth/contracts"
	"github.com/AtlantPlatform/atlant-go/fs"
	"github.com/AtlantPlatform/atlant-go/logging"
)

func init() {
	// this init will be called only within builds that have `-tags testing`
	defaultLogLevel = "4"

	testingCommands = []testingCmd{{
		Name: "test-ipfs-put",
		Desc: "Test for IPFS API: Put object method",
		Init: testIpfsPut,
	}, {
		Name: "test-ipfs-get",
		Desc: "Test for IPFS API: Get object method",
		Init: testIpfsGet,
	}, {
		Name: "test-ipfs-list",
		Desc: "Test for IPFS API: List objects method",
		Init: testIpfsList,
	}, {
		Name: "test-ipfs-pubsub",
		Desc: "Test for IPFS pub sub",
		Init: testIpfsPubSub,
	}, {
		Name: "test-ipfs-listen",
		Desc: "Test for IPFS pub sub",
		Init: testIpfsListen,
	}, {
		Name: "test-ipfs-dial",
		Desc: "Test for IPFS pub sub",
		Init: testIpfsDial,
	}, {
		Name: "test-eth-balance",
		Desc: "Test for eth balance",
		Init: testEthBalance,
	}, {
		Name: "test-eth-kyc",
		Desc: "Test for kyc contract",
		Init: testEthKYC,
	}, {
		Name: "test-authcenter",
		Desc: "Test for authcenter",
		Init: testAuthCenter,
	}}
}

func testEthKYC(c *cli.Cmd) {
	ACCOUNT := "0xc669b4edba4491c4f03060577d732bb228b9b7b8"
	ACCOUNT2 := "0xc669b4edba4491c4f03060577d732bb228b9b7b9"
	KYCABI := []byte("[{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"approveAddr\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"registerProvider\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getStatus\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"kycStatus\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"removeProvider\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"suspendAddr\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"ProviderAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"ProviderRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AddrApproved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"by\",\"type\":\"address\"}],\"name\":\"AddrSuspended\",\"type\":\"event\"}]")
	//KYCBin := `0x6060604052341561000f57600080fd5b60028054600160a060020a03191633600160a060020a0316179055610532806100396000396000f3006060604052600436106100825763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663055273c981146100875780630e260016146100a857806330ccebb5146100c757806346cc599e146100fc5780638a355a571461011b5780638da5cb5b1461013a578063996a6f8214610169575b600080fd5b341561009257600080fd5b6100a6600160a060020a0360043516610188565b005b34156100b357600080fd5b6100a6600160a060020a0360043516610254565b34156100d257600080fd5b6100e6600160a060020a03600435166102c5565b60405160ff909116815260200160405180910390f35b341561010757600080fd5b6100e6600160a060020a03600435166102e3565b341561012657600080fd5b6100a6600160a060020a03600435166102f8565b341561014557600080fd5b61014d610369565b604051600160a060020a03909116815260200160405180910390f35b341561017457600080fd5b6100a6600160a060020a0360043516610378565b60025460009033600160a060020a03908116911614806101ae57506101ae600033610440565b15156101b957600080fd5b50600160a060020a03811660009081526001602081905260409091205460ff16908114156101e657600080fd5b600160a060020a03828116600090815260016020819052604091829020805460ff1916909117905533909116907fa3673b71ec0beba775defcf8c7ad027536fdbac996023d594b5efe0c4181acd090849051600160a060020a03909116815260200160405180910390a25050565b60025433600160a060020a0390811691161461026f57600080fd5b61027a600082610463565b151561028557600080fd5b7fae9c2c6481964847714ce58f65a7f6dcc41d0d8394449bacdf161b5920c4744a81604051600160a060020a03909116815260200160405180910390a150565b600160a060020a031660009081526001602052604090205460ff1690565b60016020526000908152604090205460ff1681565b60025433600160a060020a0390811691161461031357600080fd5b61031e6000826104b6565b151561032957600080fd5b7f1589f8555933761a3cff8aa925061be3b46e2dd43f621322ab611d300f62b1d981604051600160a060020a03909116815260200160405180910390a150565b600254600160a060020a031681565b60025460009033600160a060020a039081169116148061039e575061039e600033610440565b15156103a957600080fd5b50600160a060020a03811660009081526001602052604090205460ff1660028114156103d457600080fd5b600160a060020a0382811660009081526001602052604090819020805460ff1916600217905533909116907fc17dc8dc1039ea489e8aa5d0bd8bfc32c1afb87b98959710df9646cd80c331cb90849051600160a060020a03909116815260200160405180910390a25050565b600160a060020a03811660009081526020839052604090205460ff165b92915050565b600160a060020a03811660009081526020839052604081205460ff161561048c5750600061045d565b50600160a060020a0316600090815260209190915260409020805460ff1916600190811790915590565b600160a060020a03811660009081526020839052604081205460ff1615156104e05750600061045d565b50600160a060020a0316600090815260209190915260409020805460ff191690556001905600a165627a7a72305820415f5abe7dce20370befdff065f7754fa832d832658b540ded160c3e2b3085d30029`
	PATHONE := "/Projects/go/src/github.com/AtlantPlatform/ethfw/var/chain/geth.ipc"
	PATHTWO := "/go/src/github.com/AtlantPlatform/ethfw/var/chain/geth.ipc"

	curUsr, err := user.Current()
	if err != nil {
		log.Warningf("can't obtain <home> path: %v", err)
	}
	sock := curUsr.HomeDir + PATHONE
	cli, err := eth.NewClient(sock)
	if err != nil {
		cli, _ = eth.NewClient(curUsr.HomeDir + PATHTWO)
	}
	//KYC
	kyc, err := contracts.BindKYC(cli, KYCABI)
	if err != nil {
		log.Warningf("can't bind kyc: %v", err)
	}
	tier, err := kyc.GetInfo(ACCOUNT, ACCOUNT2)
	fmt.Println("KYC INFO: ", tier, err)

	cli.Close()
}

func testEthBalance(c *cli.Cmd) {
	curUsr, err := user.Current()
	if err != nil {
		log.Warningf("can't obtain <home> path: %v", err)
	}
	path := "/.ethereum/geth.ipc"
	sock := curUsr.HomeDir + path
	cli, _ := eth.NewClient(sock)
	//ETH
	eth := contracts.BindETH(cli)
	bal, err := eth.BalanceOf("0xEcE807666997519F2C2b5199EB4261b03392756a")
	fmt.Println(bal, err, "ETH")
	//PTO
	pto, err := contracts.BindPTO(cli, nil)
	if err != nil {
		log.Errorln(err)
	}
	result, err := pto.BalanceOf("0x6ed945867ac5dc998e8c6aae4e437225345a6302")
	fmt.Println(result, err, "PTO")
	//ATL
	atl, err := contracts.BindATL(cli, nil)
	if err != nil {
		log.Errorln(err)
	}
	result, err = atl.BalanceOf("0xa12431d0b9db640034b0cdfceef9cce161e62be4")
	fmt.Println(result, err, "ATL")
	return
}

func testIpfsPut(c *cli.Cmd) {
	meta := logging.WithFn()
	id := c.StringOpt("I id", "", "Specify object ID.")
	path := c.StringOpt("P path", "hello.txt", "Specify object path.")
	version := c.StringOpt("V version", "", "Specify object version (IPFS CID).")
	prevVersion := c.StringOpt("P prev-version", "", "Specify object previous version (IPFS CID).")
	content := c.StringArg("CONTENT", "hello world!", "Content to store.")
	c.Action = func() {
		printInfo(meta)
		runWithPlanetaryContext(func(ctx PlanetaryContext) {
			log.Println("IPFS identity:", ctx.FileStore().NodeID())
			ref := fs.ObjectRef{
				ID:              *id,
				Path:            *path,
				Version:         *version,
				VersionPrevious: *prevVersion,
			}
			body := ioutil.NopCloser(bytes.NewReader([]byte(*content)))
			newRef, err := ctx.FileStore().PutObject(context.Background(), ref, body)
			if err != nil {
				log.Errorln(err)
				return
			}
			log.Printf("New object ID=%s, Version=%s, VersionPrevious=%s", newRef.ID, newRef.Version, newRef.VersionPrevious)
		})
	}
}

func testIpfsGet(c *cli.Cmd) {
	meta := logging.WithFn()
	id := c.StringOpt("I id", "", "Specify object ID.")
	path := c.StringOpt("P path", "hello.txt", "Specify object path.")
	version := c.StringOpt("V version", "", "Specify object version (IPFS CID).")
	prevVersion := c.StringOpt("P prev-version", "", "Specify object previous version (IPFS CID).")
	versionOffset := c.IntOpt("O offset", 0, "Specify version offset.")
	c.Action = func() {
		printInfo(meta)
		runWithPlanetaryContext(func(ctx PlanetaryContext) {
			log.Println("IPFS identity:", ctx.FileStore().NodeID())
			ref := fs.ObjectRef{
				ID:              *id,
				Path:            *path,
				Version:         *version,
				VersionOffset:   *versionOffset,
				VersionPrevious: *prevVersion,
			}
			newRef, err := ctx.FileStore().GetObject(context.Background(), ref)
			if err != nil {
				log.Errorln(err)
				return
			}
			log.Printf("Loaded object ID=%s, Version=%s VersionOffset=%d", newRef.ID, newRef.Version, newRef.VersionOffset)
			v, _ := ioutil.ReadAll(newRef.Body)
			if len(v) > 0 {
				fmt.Println("Body:", string(v))
			} else {
				fmt.Println("No body.")
			}
		})
	}
}

func testIpfsList(c *cli.Cmd) {
	meta := logging.WithFn()
	id := c.StringOpt("I id", "", "Specify object ID.")
	path := c.StringOpt("P path", "hello.txt", "Specify object path.")
	version := c.StringOpt("V version", "", "Specify object version (IPFS CID).")
	prevVersion := c.StringOpt("P prev-version", "", "Specify object previous version (IPFS CID).")
	versionOffset := c.IntOpt("O offset", 0, "Specify version offset.")
	c.Action = func() {
		printInfo(meta)
		runWithPlanetaryContext(func(ctx PlanetaryContext) {
			log.Println("IPFS identity:", ctx.FileStore().NodeID())
			ref := fs.ObjectRef{
				ID:              *id,
				Path:            *path,
				Version:         *version,
				VersionOffset:   *versionOffset,
				VersionPrevious: *prevVersion,
			}
			list, err := ctx.FileStore().ListObjects(context.Background(), ref)
			if err != nil {
				log.Errorln(err)
				return
			}
			v, _ := json.MarshalIndent(list, "", "\t")
			fmt.Printf("List: %v\n", string(v))
		})
	}
}

func testIpfsPubSub(c *cli.Cmd) {
	meta := logging.WithFn()
	c.Action = func() {
		printInfo(meta)
		runWithPlanetaryContext(func(ctx PlanetaryContext) {
			log.Println("IPFS identity:", ctx.FileStore().NodeID())
			ps, err := ctx.FileStore().PubSub()
			if err != nil {
				log.Errorln(err)
				return
			}
			listenOn := func(topic string, wg *sync.WaitGroup) {
				if err := ps.Subscribe(topic, func(m *fs.Message) error {
					log.Printf("message from %s on %s: %s", m.From, m.TopicIDs[0], m.Data)
					wg.Done()
					return nil
				}); err != nil {
					log.Errorln(err)
					return
				}
			}
			wg := new(sync.WaitGroup)
			wg.Add(4)
			{
				listenOn("test01", wg)
				listenOn("test02", wg)
				ps.Publish("test01", []byte("hello01"))
				ps.Publish("test02", []byte("hello01"))
				ps.Publish("test01", []byte("hello02"))
				ps.Publish("test01", []byte("hello03"))
			}
			wg.Wait()
		})
	}
}

func testAuthority(c *cli.Cmd) {
	meta := logging.WithFn()
	c.Action = func() {
		printInfo(meta)
		runWithPlanetaryContext(func(ctx PlanetaryContext) {
			dnsAuth := authority.GetDnsAuthorityInstance()
			auth := authority.NewAuthority(dnsAuth)
			nodeID := ctx.FileStore().NodeID()
			isWhiteNode := auth.Check(nodeID)
			log.Info("Node id:", nodeID)
			log.Info("Node in white list:", isWhiteNode)

			printMap := func(pki authority.PKI) {
				for k, v := range auth.List() {
					log.Info(k)
					for _, j := range v {
						log.Info("\t", j)
					}
				}
			}

			printMap(auth)
		})
	}
}

func printInfo(meta log.Fields) {
	log.WithFields(meta).Println("---TEST START---")
	closer.Bind(func() {
		log.WithFields(meta).Println("---TEST ENDED---")
	})
}

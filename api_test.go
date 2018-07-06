package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"flag"
	"github.com/AtlantPlatform/atlant-go/api"
	"github.com/AtlantPlatform/atlant-go/contracts"
	"github.com/AtlantPlatform/atlant-go/proto"
	"github.com/AtlantPlatform/atlant-go/rs"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/xlab/closer"
	"strings"
	"time"
	"github.com/ethereum/go-ethereum/core/types"
)

func init() {
	defaultLogLevel = "5"

	ethAddress = flag.String("E", "0x0", "eth wallet addr")
	fsDir = flag.String("F", "var/fs", "file storage path")
	stateDir = flag.String("S", "var/state", "db storage addr")
	fsListenAddr = flag.String("L", ":33770", "file storage port")
	webListenAddr = flag.String("W", "0.0.0.0:33780", "public socket")
	flag.Parse()
}

func TestMain(m *testing.M) {
	if !isNodeRunning() {
		go main()
	}
	//go setUpPrivateAPI() //use instead of main()
	pingForSync()
	code := m.Run()
	closer.Close()
	os.Exit(code)
}

func isNodeRunning() bool {
	_, err := http.Get("http://" + *webListenAddr + "/checkForRunning")
	if err == nil {
		return true
	}
	return false
}

func pingForSync() {
	for i := 0; i < 200; i++ {
		_, err := http.Get("http://" + *webListenAddr + "/checkForSync")
		if i%10 == 0 {
			fmt.Printf("waits times #%d: %#v\n", i, err)
		}
		if err == nil {
			return
		}
		time.Sleep(5 * time.Second)
	}
	log.Println("waiting for node syncing too long")
	os.Exit(1)
}

func setUpPrivateAPI() {
	runWithPlanetaryContext(func(ctx PlanetaryContext) {
		store, err := rs.NewPlanetaryRecordStore(ctx.NodeID(), ctx.FileStore(), ctx.StateStore())
		if err != nil {
			log.Fatalln(err)
		}
		closer.Bind(func() {
			if err := store.Close(); err != nil {
				log.Warningln(err)
			}
		})
		apiCtx := api.NewContext(ctx, store, nil, *ethAddress, "")

		privateServer := api.NewPrivateServer()
		privateServer.RouteAPI(apiCtx)
		_, err = privateServer.Listen("127.0.0.1:33999")
		if err != nil {
			log.Fatalln(err)
		}
	})
}

//positive
func TestPing(t *testing.T) {
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/ping", *webListenAddr))
	r.NoError(err)
	r.Equal(200, resp.StatusCode)
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	r.Equal(len(bytes), 49)
	fmt.Println(string(bytes))
}

func TestTokenDistributionInfo(t *testing.T) {
	r := require.New(t)
	var disInfo api.DistributionInfo
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/tokenDistributionInfo", *webListenAddr))
	r.NoError(err)
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	err = json.Unmarshal(bytes, &disInfo)
	r.NoError(err)
	fmt.Printf("%#v\n", disInfo)
	r.Equal(200, resp.StatusCode)
}

func TestKYCStatus(t *testing.T) {
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/kycStatus", *webListenAddr))
	r.NoError(err)
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	switch string(bytes) {
	case string(contracts.StatusUnknown), string(contracts.StatusApproved), string(contracts.StatusSuspended):
		fmt.Println(string(bytes))
	default:
		t.Error("error while getting kyc status:", string(bytes))
		return
	}
	r.Equal(200, resp.StatusCode)
}

func TestKYCApprove(t *testing.T) {
	var tx types.Transaction
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/kycApprove", *webListenAddr))
	r.NoError(err)
	r.Equal(200, resp.StatusCode)
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	err = json.Unmarshal(bytes, &tx)
	r.NoError(err)
	r.NotEqual(len(tx.String()), 0)
}

func TestKYCSuspend(t *testing.T) {
	var tx types.Transaction
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/kycSuspend", *webListenAddr))
	r.NoError(err)
	r.Equal(200, resp.StatusCode)
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	err = json.Unmarshal(bytes, &tx)
	r.NoError(err)
	r.NotEqual(len(tx.String()), 0)
}

func TestEthBalance(t *testing.T) {
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/ethBalance", *webListenAddr))
	r.NoError(err)
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	fmt.Println(string(bytes))
	r.Equal(200, resp.StatusCode)
}

func TestATLBalance(t *testing.T) {
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/atlBalance", *webListenAddr))
	r.NoError(err)
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	fmt.Println(string(bytes))
	r.Equal(200, resp.StatusCode)
}

func TestPTOBalance(t *testing.T) {
	r := require.New(t)
	resp, err := http.Get(
		fmt.Sprintf("http://%s/api/v1/ptoBalance/0xa936055b4c9b4a1213e64b7fc8c7ff295939ce71", *webListenAddr),
	)
	r.NoError(err)
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	fmt.Println(string(bytes))
	r.Equal(200, resp.StatusCode)
}

func TestPut(t *testing.T) {
	r := require.New(t)
	resp, err := http.Post(fmt.Sprintf("http://%s/api/v1/put/content/test.txt", *webListenAddr), "", nil)
	r.NoError(err)

	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	fmt.Println(string(bytes))
	r.Equal(200, resp.StatusCode)
}

func TestNewID(t *testing.T) {
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/newID", *webListenAddr))
	r.NoError(err)

	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	r.Equal(len(bytes), 26)
	fmt.Println(string(bytes))
	r.Equal(200, resp.StatusCode)
}

func TestEnv(t *testing.T) {
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/env", *webListenAddr))
	r.NoError(err)
	r.Equal(200, resp.StatusCode)
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	switch string(bytes) {
	case "main", "test":
		fmt.Println(string(bytes))
	default:
		t.Error("unknown env:", string(bytes))
		return
	}
}

func TestSession(t *testing.T) {
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/session", *webListenAddr))
	r.NoError(err)
	r.Equal(200, resp.StatusCode)
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	r.Equal(len(bytes), 26)
	fmt.Println(string(bytes))
}

func TestVersion(t *testing.T) {
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/version", *webListenAddr))
	r.NoError(err)
	r.Equal(200, resp.StatusCode)
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	r.NotEqual(len(bytes), 0)
	fmt.Println(string(bytes))
}

func TestStats(t *testing.T) {
	r := require.New(t)
	var stats api.Stats
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/stats", *webListenAddr))
	r.NoError(err)
	r.Equal(200, resp.StatusCode)
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	err = json.Unmarshal(bytes, &stats)
	r.NoError(err)
	r.NotEqual(len(bytes), 0)
	fmt.Println(stats.Uptime)
}

func TestLogs(t *testing.T) {
	r := require.New(t)
	var logs []string
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/logs", *webListenAddr))
	r.NoError(err)
	r.Equal(200, resp.StatusCode)
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	if len(bytes) != 0 {
		err := json.Unmarshal(bytes, &logs)
		r.NoError(err)
		fmt.Println(logs)
		return
	}
	fmt.Println("empty logs")
}

func TestLog(t *testing.T) {
	r := require.New(t)
	dateStr := time.Now().Local().Format("2006-01-02")
	queryDate := strings.Replace(dateStr, "-", "/", -1)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/log/%s", *webListenAddr, queryDate))
	r.NoError(err)
	r.Equal(200, resp.StatusCode)
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	fmt.Println(string(bytes))
}

func TestAssets(t *testing.T) {
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/assets", *webListenAddr))
	r.NoError(err)
	r.Equal(200, resp.StatusCode)
}

func TestIndex(t *testing.T) {
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/index/configs/", *webListenAddr))
	r.NoError(err)
	r.Equal(200, resp.StatusCode)
}

func TestContent(t *testing.T) {
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/content/configs/pto/atl123.json", *webListenAddr))
	r.NoError(err)
	r.Equal(200, resp.StatusCode)
	//r.NotEmpty(resp.Header["X-Meta-ID"])
	r.NotEmpty(resp.Header["X-Meta-Version"])
	//r.NotEmpty(resp.Header["X-Meta-Previous"])
	r.NotEmpty(resp.Header["X-Meta-Path"])
	//r.NotEmpty(resp.Header["X-Meta-UserMeta"])
	_, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
}

func TestListAll(t *testing.T) {
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/listAll/configs/", *webListenAddr))
	r.NoError(err)
	r.Equal(200, resp.StatusCode)
}

func TestListVersions(t *testing.T) {
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/listVersions/configs/pto/atl123.json", *webListenAddr))
	r.NoError(err)
	r.Equal(200, resp.StatusCode)
}

func TestMeta(t *testing.T) {
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/meta/configs/pto/atl123.json", *webListenAddr))
	r.NoError(err)
	r.Equal(200, resp.StatusCode)
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	var meta proto.ObjectMeta
	err = json.Unmarshal(bytes, &meta)
	r.NoError(err)
	fmt.Printf("%#v\n", meta)
}

func TestPrivatePing(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	r := require.New(t)
	resp, err := http.Get("http://127.0.0.1:33999/private/v1/ping")
	r.NoError(err)
	r.Equal(200, resp.StatusCode)

}

func TestPrivateRecords(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	r := require.New(t)
	resp, err := http.Get("http://127.0.0.1:33999/private/v1/records")
	r.NoError(err)
	r.Equal(200, resp.StatusCode)
}

func TestPrivateAnnounce(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	r := require.New(t)
	resp, err := http.Post("http://127.0.0.1:33999/private/v1/announce", "", nil)
	r.NoError(err)
	r.Equal(200, resp.StatusCode)
}

//negative
func TestFailEmptyPTOBalance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	fmt.Println("======================NEGATIVE testing===================")
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/ptoBalance/", *webListenAddr))
	r.NoError(err)
	r.Equal(404, resp.StatusCode)
}

func TestFailImproperFaPTOBalance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/ptoBalance/****", *webListenAddr))
	r.NoError(err)
	r.Equal(500, resp.StatusCode)
}

func TestFailEmptyLog(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/log/", *webListenAddr))
	r.NoError(err)
	r.Equal(404, resp.StatusCode)
}

func TestFailImproperLog(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/log/9999/99/99", *webListenAddr))
	r.NoError(err)
	r.Equal(404, resp.StatusCode)
}

func TestFailNotAuthorizedPut(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	r := require.New(t)
	resp, err := http.Post(fmt.Sprintf("http://%s/api/v1/put/nothing", *webListenAddr), "", nil)
	r.NoError(err)
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	r.Contains(string(bytes), "not authorized")
	r.Equal(500, resp.StatusCode)
}

func TestFailEmptyPut(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	r := require.New(t)
	resp, err := http.Post(fmt.Sprintf("http://%s/api/v1/put", *webListenAddr), "", nil)
	r.NoError(err)
	r.Equal(400, resp.StatusCode)
}

func TestFailDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	r := require.New(t)
	resp, err := http.Post(
		fmt.Sprintf("http://%s/api/v1/delete/please", *webListenAddr),
		"",
		nil,
	)
	r.NoError(err)
	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r.NoError(err)
	r.Contains(string(bytes), "not authorized")
	r.Equal(500, resp.StatusCode)
}

func TestFailEmptyDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	r := require.New(t)
	resp, err := http.Post(fmt.Sprintf("http://%s/api/v1/delete/", *webListenAddr), "", nil)
	r.NoError(err)
	r.Equal(404, resp.StatusCode)
}

func TestFailContent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/content/```", *webListenAddr))
	r.NoError(err)
	r.Equal(404, resp.StatusCode)
}

func TestFailMeta(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/meta/somemetaplease", *webListenAddr))
	r.NoError(err)
	r.Equal(404, resp.StatusCode)
}

func TestFailListVersions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/aaaaaaaaaahelpme", *webListenAddr))
	r.NoError(err)
	r.Equal(404, resp.StatusCode)
}

func TestFailListAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/listAll/abbracadabbra/", *webListenAddr))
	r.NoError(err)
	r.Equal(404, resp.StatusCode)
}

func TestFailPTO(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/ptoBalance/fuckyou", *webListenAddr))
	r.NoError(err)
	r.Equal(500, resp.StatusCode)
}

func TestFailEmptyPTO(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/ptoBalance", *webListenAddr))
	r.NoError(err)
	r.Equal(404, resp.StatusCode)
}

func TestFailIndex(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	r := require.New(t)
	resp, err := http.Get(fmt.Sprintf("http://%s/index/canttouchthis", *webListenAddr))
	r.NoError(err)
	r.Equal(404, resp.StatusCode)
}

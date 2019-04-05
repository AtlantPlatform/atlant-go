// Copyright 2017-2019 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package fs

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sync"

	files "github.com/ipfs/go-ipfs-files"
	coreapi "github.com/ipfs/go-ipfs/core/coreapi"
	options "github.com/ipfs/interface-go-ipfs-core/options"
	log "github.com/sirupsen/logrus"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/corerepo"

	bitswap "github.com/ipfs/go-bitswap"
	cid "github.com/ipfs/go-cid"
	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/namesys"
	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	ipld "github.com/ipfs/go-ipld-format"
	ipath "github.com/ipfs/go-path"
	"github.com/ipfs/go-path/resolver"
	uio "github.com/ipfs/go-unixfs/io"
	ipnet "github.com/libp2p/go-libp2p-interface-pnet"
	peer "github.com/libp2p/go-libp2p-peer"
	pnet "github.com/libp2p/go-libp2p-pnet"

	"github.com/AtlantPlatform/atlant-go/logging"
	"github.com/AtlantPlatform/atlant-go/proto"
)

func init() {
	ipnet.ForcePrivateNetwork = true
}

// ipfsStore implements PlanetaryFileStore.
type ipfsStore struct {
	prefix string
	opts   *ipfsOptions
	node   *core.IpfsNode
	repo   repo.Repo
	resolv *resolver.Resolver

	pubsub     *ipfsPubSub
	pubsubOnce sync.Once

	listener     *p2pListener
	listenerOnce sync.Once
	client       *p2pClient
	clientOnce   sync.Once
}

func (s *ipfsStore) NodeID() string {
	return s.node.Identity.Pretty()
}

// ErrNotFound is thrown if boject is not found in IPFS
var ErrNotFound = errors.New("not found")

func (s *ipfsStore) PutObject(ctx context.Context, ref ObjectRef,
	userMeta []byte, body io.ReadCloser) (*ObjectRef, error) {
	return s.putObject(ctx, ref, userMeta, body, false)
}

func (s *ipfsStore) DeleteObject(ctx context.Context, ref ObjectRef) (*ObjectRef, error) {
	// also unpin previous versions
	return s.putObject(ctx, ref, nil, nil, true)
}

func (s *ipfsStore) putObject(ctx context.Context, ref ObjectRef,
	userMeta []byte, body io.ReadCloser, isDelete bool) (*ObjectRef, error) {
	// fileAdder, err := coreunix.NewAdder(ctx, s.node.Pinning, s.node.Blockstore, s.node.DAG)
	// if err != nil {
	// err = fmt.Errorf("failed to init IPFS file adder: %v", err)
	// return nil, err
	// }
	if len(ref.ID) == 0 {
		ref.ID = proto.NewID()
	}
	meta, err := ref.ToProto()
	if err != nil {
		err = fmt.Errorf("failed to create object meta: %v", err)
		return nil, err
	}
	if isDelete {
		meta.SetIsDeleted(true)
	}
	meta.SetUserMeta(string(userMeta))
	dir, err := NewObjectDir(meta, body)
	if err != nil {
		err = fmt.Errorf("failed to create object directory: %v", err)
		return nil, err
	}
	api, err := coreapi.NewCoreAPI(s.node)
	if err != nil {
		err = fmt.Errorf("failed to init core API: %v", err)
		return nil, err
	}
	path, err := api.Unixfs().Add(context.Background(), dir.(files.Directory), options.Unixfs.Wrap(true))
	if err != nil {
		err = fmt.Errorf("failed to Add API: %v", err)
		return nil, err
	}
	log.Infoln("[ipfsStore.putObject] node.path=", path.String())
	node, err := api.Object().Get(context.Background(), path)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		err = fmt.Errorf("failed to add and pin object file to DAG: %v", err)
		return nil, err
	}
	ref.Version = node.Cid().String()
	if err := s.PinNewest(ref, 3); err != nil {
		err = fmt.Errorf("failed to pin object file, it will be soon collected by GC: %v", err)
		return nil, err
	}
	meta.SetVersion(ref.Version)
	ref.SetMeta(&meta)
	return &ref, nil
}

func (s *ipfsStore) HeadObject(ctx context.Context, ref ObjectRef) (*ObjectRef, error) {
	normRef := s.resolveObjectVersion(ctx, ref)
	if normRef == nil || normRef.Meta() == nil {
		normRef = s.cidToObjectRef(ctx, normRef.Version)
		if normRef == nil || normRef.Meta() == nil {
			return nil, ErrNotFound
		}
	}
	if _, err := ipath.ParseCidToPath(normRef.Version); err != nil {
		err = fmt.Errorf("failed to parse object version CID: %v", err)
		return nil, err
	}
	return normRef, nil
}

func (s *ipfsStore) GetObject(ctx context.Context, ref ObjectRef) (*Object, error) {
	normRef := s.resolveObjectVersion(ctx, ref)
	if normRef == nil || normRef.Meta() == nil {
		normRef = s.cidToObjectRef(ctx, normRef.Version)
		if normRef == nil || normRef.Meta() == nil {
			return nil, ErrNotFound
		}
	}
	p, err := ipath.ParseCidToPath(normRef.Version)
	if err != nil {
		err = fmt.Errorf("failed to parse object version CID: %v", err)
		return nil, err
	}
	obj := &Object{
		ObjectRef: *normRef,
		Meta:      normRef.Meta(),
	}
	if normRef.Meta().IsDeleted() {
		return obj, nil
	}
	dagNode, err := core.Resolve(ctx, s.node.Namesys, s.resolv, p)
	if err != nil {
		return obj, ErrNotFound
	}

	var contentNode ipld.Node
	for _, link := range dagNode.Links() {
		if link.Name == "content" {
			n, err := link.GetNode(ctx, s.node.DAG)
			if err != nil {
				err = fmt.Errorf("failed to get object content node: %v", err)
				return nil, err
			}
			contentNode = n
			break
		}
	}
	if contentNode == nil {
		return obj, ErrNotFound
	}
	reader, err := uio.NewDagReader(ctx, contentNode, s.node.DAG)
	if err != nil {
		err = fmt.Errorf("failed to read node content: %v", err)
		return obj, err
	}
	obj.Body = reader
	return obj, nil
}

func (s *ipfsStore) resolveObjectVersion(ctx context.Context, ref ObjectRef) *ObjectRef {
	if ref.VersionOffset == 0 {
		return &ref
	}
	if ref.VersionOffset > 0 {
		// TODO(max): cannot handle such cases yet.
		// Needs caching implemetation in place for ipfsStore.
		return &ref
	}
	if len(ref.VersionPrevious) > 0 {
		ref.Version = ref.VersionPrevious
		ref.VersionPrevious = ""
		ref.VersionOffset++
	}
	obj := s.cidToObjectRef(ctx, ref.Version)
	for obj != nil && ref.VersionOffset < 0 {
		if len(obj.VersionPrevious) > 0 {
			ref.Version = obj.VersionPrevious
			ref.VersionOffset++
		} else {
			return &ref
		}
	}
	return &ref
}

func (s *ipfsStore) ListObjects(ctx context.Context, ref ObjectRef) ([]ObjectRef, error) {
	var list []ObjectRef
	if len(ref.Version) > 0 {
		if ref.VersionOffset == 0 {
			// specific version
			if obj := s.cidToObjectRef(ctx, ref.Version); obj != nil {
				list = append(list, *obj)
				return list, nil
			}
			return nil, ErrNotFound
		}
		return nil, errors.New("ListObjects: version offsets are not supported yet (VersionOffset == 0)")
	} else if len(ref.Version) > 0 && ref.VersionOffset != 0 {
		return nil, errors.New("ListObjects: version offsets are not supported yet (VersionOffset != 0)")
	}
	cidC, err := s.node.Blockstore.AllKeysChan(ctx)
	if err != nil {
		return nil, err
	}
	for c := range cidC {
		cStr := c.String()
		var match bool
		if len(ref.Version) > 0 {
			match = (cStr == ref.Version)
		} else if len(ref.VersionPrevious) > 0 {
			match = (cStr == ref.VersionPrevious)
		} else {
			match = true
		}
		if !match {
			// filtered out
			continue
		}
		if obj := s.cidToObjectRef(ctx, cStr); obj != nil {
			if len(ref.ID) > 0 && obj.ID != ref.ID {
				// filtered out
				continue
			} else if len(ref.Path) > 0 && obj.Path != ref.Path {
				// filtered out
				continue
			}
			list = append(list, *obj)
		}
	}
	return list, nil
}

func (s *ipfsStore) PinObject(ref ObjectRef) error {
	p, err := ipath.ParseCidToPath(ref.Version)
	if err != nil {
		log.WithFields(logging.WithFn()).Errorln("failed to parse object CID:", err)
		return err
	}
	dagNode, err := core.Resolve(s.node.Context(), s.node.Namesys, s.resolv, p)
	if err != nil {
		return err
	}
	if err := s.node.Pinning.Pin(s.node.Context(), dagNode, true); err != nil {
		return err
	}
	return s.node.Pinning.Flush()
}

func (s *ipfsStore) UnpinObject(ref ObjectRef) error {
	id, err := cid.Parse(ref.Version)
	if err != nil {
		log.WithFields(logging.WithFn()).Errorln("failed to parse object CID:", err)
		return err
	}
	if err := s.node.Pinning.Unpin(s.node.Context(), id, true); err != nil {
		return err
	}
	return s.node.Pinning.Flush()
}

func (s *ipfsStore) PinNewest(ref ObjectRef, depth int) error {
	if err := s.PinObject(ref); err != nil {
		return err
	}
	if ref.VersionPrevious == "" || depth < 0 {
		return nil
	}

	prevVer := ref.VersionPrevious
	for prevVer != "" {
		if depth--; depth < 0 {
			objRef := s.cidToObjectRef(s.node.Context(), prevVer)
			prevVer = objRef.VersionPrevious
			id, err := cid.Parse(objRef.Version)
			if err != nil {
				return fmt.Errorf("failed to parse ID to CID: %v", err)
			}
			if _, ok, _ := s.node.Pinning.IsPinned(id); !ok {
				continue
			}
			if err := s.node.Pinning.Unpin(s.node.Context(), id, true); err != nil {
				return err
			}
			if _, ok, _ := s.node.Pinning.IsPinned(id); ok {
				return fmt.Errorf("failed object unpinning: cid %s", id.String())
			}
			log.Debugf("successful unpinning cid: %s", id.String())
		}
	}
	return s.node.Pinning.Flush()
}

func (s *ipfsStore) cidToObjectRef(ctx context.Context, cid string) *ObjectRef {
	p, err := ipath.ParseCidToPath(cid)
	if err != nil {
		log.WithFields(logging.WithFn()).Errorln("failed to parse object CID:", err)
		return nil
	}
	dagNode, err := core.Resolve(ctx, s.node.Namesys, s.resolv, p)
	if err != nil {
		return nil
	}
	var metaNode ipld.Node
	for _, link := range dagNode.Links() {
		if link.Name == "meta" {
			m, err := link.GetNode(ctx, s.node.DAG)
			if err != nil {
				log.WithFields(logging.WithFn()).Warningln("failed to get link node:", err)
				return nil
			}
			metaNode = m
			break
		}
	}
	if metaNode == nil {
		return nil
	}
	reader, err := uio.NewDagReader(ctx, metaNode, s.node.DAG)
	if err != nil {
		log.WithFields(logging.WithFn()).Warningf("no reader for meta node %s: %v", cid, err)
		return nil
	}
	var meta proto.ObjectMeta
	func() {
		defer reader.Close()
		// TODO(max): potential buffer reuse for multiple cidToObjectRef calls.
		if m, err := readObjectFileMeta(reader); err != nil {
			log.WithFields(logging.WithFn()).Warningf("failed to read object file meta: %v", err)
		} else {
			meta = m
		}
	}()
	if len(meta.IdBytes()) == 0 {
		log.WithFields(logging.WithFn()).Warningln("empty meta for", cid)
		return nil
	}
	meta.SetVersion(cid)
	ref := &ObjectRef{
		ID:   meta.Id(),
		Path: meta.Path(),
		Size: meta.Size(),

		Version:         cid,
		VersionPrevious: meta.VersionPrevious(),

		meta: &meta,
	}
	return ref
}

func (s *ipfsStore) PubSub() (PlanetaryPubSub, error) {
	if !s.opts.PubSubEnabled {
		return nil, ErrNoPubSub
	}
	s.pubsubOnce.Do(func() {
		s.pubsub = newIpfsPubSub(s.node)
	})
	if s.pubsub == nil {
		return nil, ErrNoPubSub
	}
	return s.pubsub, nil
}

func (s *ipfsStore) Listener() PlanetaryListener {
	s.listenerOnce.Do(func() {
		s.listener = newListener(s.node)
	})
	return s.listener
}

func (s *ipfsStore) Client() PlanetaryClient {
	s.clientOnce.Do(func() {
		s.client = newClient(s.node)
	})
	return s.client
}

func newIpfsStore(prefix string, needInit bool, opts ...IpfsOpt) (*ipfsStore, error) {
	s := &ipfsStore{
		prefix: prefix,
		opts:   defaultIpfsOptions(),
	}
	for _, o := range opts {
		if o != nil {
			o(s.opts)
		}
	}
	cfg := &core.BuildCfg{
		Online: true,
		ExtraOpts: map[string]bool{
			"pubsub": s.opts.PubSubEnabled,
			"ipnsps": false,
			"mplex":  false,
		},
	}

	if s.opts.StoreEnabled {
		locked, err := fsrepo.LockedByOtherProcess(prefix)
		if err != nil {
			return nil, err
		} else if locked {
			err := fmt.Errorf("specified fs store prefix is locked by another process (prefix=%s)", prefix)
			return nil, err
		}
		if needInit {
			if err := checkWriteable(prefix); err != nil {
				return nil, err
			}
			conf, err := config.Init(ioutil.Discard, 2048)
			if err != nil {
				return nil, err
			}

			// force signature verification on initialization
			conf.Pubsub.DisableSigning = false
			conf.Pubsub.StrictSignatureVerification = true

			// force use of BadgerDB upon the init
			if err := config.Profiles["badgerds"].Transform(conf); err != nil {
				log.Warningf("failed to apply badgerds profile: %v", err)
				return nil, err
			}
			if err := fsrepo.Init(prefix, conf); err != nil {
				return nil, err
			}
			if err := initializeIpnsKeyspace(prefix); err != nil {
				return nil, err
			}
			cfg.Online = false
		}
		r, err := s.openRepo(prefix)
		if err != nil {
			return nil, err
		}
		cfg.Permanent = true
		cfg.Repo = r
		s.repo = r
	} else {
		cfg.Permanent = false
		cfg.NilRepo = true
	}
	n, err := core.NewNode(context.Background(), cfg)
	if err != nil {
		return nil, err
	}
	s.node = n
	s.resolv = &resolver.Resolver{
		DAG:         n.DAG,
		ResolveOnce: uio.ResolveUnixfsOnce,
	}
	return s, nil
}

func (s *ipfsStore) applyConfig(cfg *config.Config) error {
	// network presets should be applied first
	switch s.opts.NetworkProfile {
	case NetworkDefault:
		if err := config.Profiles["default-networking"].Transform(cfg); err != nil {
			log.Warningf("failed to apply default-networking profile: %v", err)
		}
		if err := config.Profiles["local-discovery"].Transform(cfg); err != nil {
			log.Warningf("failed to apply local-discovery profile: %v", err)
		}
	case NetworkServer:
		if err := config.Profiles["default-networking"].Transform(cfg); err != nil {
			log.Warningf("failed to apply default-networking profile: %v", err)
		}
		if err := config.Profiles["server"].Transform(cfg); err != nil {
			log.Warningf("failed to apply server profile: %v", err)
		}
	case NetworkTest:
		return errors.New("network test profile is not supported yet")
		// if err := config.Profiles["test"].Transform(cfg); err != nil {
		// 	log.Warningf("failed to apply test profile: %v", err)
		// }
	}
	cfg.Ipns = config.Ipns{
		ResolveCacheSize: 1024,
	}
	if s.opts.RelayEnabled {
		cfg.Swarm.DisableRelay = false
		cfg.Swarm.EnableRelayHop = true
	} else {
		cfg.Swarm.DisableRelay = true
		cfg.Swarm.EnableRelayHop = false
	}
	cfg.Experimental.Libp2pStreamMounting = true
	cfg.Swarm.DisableBandwidthMetrics = false
	cfg.SetBootstrapPeers(s.opts.BootstrapPeers)
	cfg.Addresses.Swarm = []string{
		fmt.Sprintf("/ip4/%s/tcp/%d", s.opts.ListenHost, s.opts.ListenPort),
	}
	// disable extra IPFS networking
	cfg.Addresses.API = []string{""}
	cfg.Addresses.Gateway = []string{""}
	cfg.API = config.API{}
	cfg.Gateway = config.Gateway{}
	return nil
}

func (s *ipfsStore) openRepo(prefix string) (repo.Repo, error) {
	r, err := fsrepo.Open(prefix)
	if err != nil {
		return nil, err
	}
	cfg, err := r.Config()
	if err != nil {
		r.Close()
		return nil, err
	}
	if err := s.applyConfig(cfg); err != nil {
		r.Close()
		return nil, err
	}
	if err := r.SetConfig(cfg); err != nil {
		log.Warningf("failed to apply current options to IPFS config: %v", err)
	}
	return r, nil
}

func (s *ipfsStore) Close() error {
	if err := s.node.Close(); err != nil {
		log.Errorf("IPFS node shutdown failed: %v", err)
	}
	if err := s.pubsub.Close(); err != nil {
		log.Errorf("IPFS PubSub shutdown failed: %v", err)
	}
	return s.repo.Close()
}

func newIpfsPrivateKey() (io.Reader, error) {
	return pnet.GenerateV1PSK()
}

// checkWriteable from ipfs-go/init.go.
func checkWriteable(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		// dir exists, make sure we can write to it
		testfile := path.Join(dir, "test")
		fi, err := os.Create(testfile)
		if err != nil {
			if os.IsPermission(err) {
				return fmt.Errorf("%s is not writeable by the current user", dir)
			}
			return fmt.Errorf("unexpected error while checking writeablility of repo root: %s", err)
		}
		fi.Close()
		return os.Remove(testfile)
	}
	if os.IsNotExist(err) {
		// dir doesnt exist, check that we can create it
		return os.Mkdir(dir, 0775)
	}
	if os.IsPermission(err) {
		return fmt.Errorf("cannot write to %s, incorrect permissions", err)
	}
	return err
}

// initializeIpnsKeyspace from ipfs-go/init.go
func initializeIpnsKeyspace(prefix string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := fsrepo.Open(prefix)
	if err != nil { // NB: repo is owned by the node
		return err
	}
	nd, err := core.NewNode(ctx, &core.BuildCfg{Repo: r})
	if err != nil {
		return err
	}
	defer nd.Close()
	// if err := nd.SetupOfflineRouting(); err != nil {
	//	return err
	// }
	return namesys.InitializeKeyspace(ctx, nd.Namesys, nd.Pinning, nd.PrivateKey)
}

func (s *ipfsStore) nodesForPaths(ctx context.Context, paths []string) ([]ipld.Node, error) {
	nodes := make([]ipld.Node, len(paths))
	for i, sp := range paths {
		p, err := ipath.ParsePath(sp)
		if err != nil {
			return nil, err
		}
		node, err := core.Resolve(ctx, s.node.Namesys, s.node.Resolver, p)
		if err != nil {
			return nil, err
		}
		nodes[i] = node
	}
	return nodes, nil
}

func (s *ipfsStore) DiskStats() (*DiskStats, error) {
	// var fs syscall.Statfs_t
	// if err := syscall.Statfs(s.prefix, &fs); err != nil {
	// 	return nil, err
	// }
	ds := &DiskStats{
		// BytesAll:  fs.Blocks * uint64(fs.Bsize),
		// BytesFree: fs.Bfree * uint64(fs.Bsize),
	}
	ds.BytesUsed = ds.BytesAll - ds.BytesFree
	return ds, nil
}

func (s *ipfsStore) BandwidthStats() *BandwidthStats {
	if s.node.Reporter == nil {
		return nil
	}
	totals := s.node.Reporter.GetBandwidthTotals()
	return &BandwidthStats{
		TotalIn:  totals.TotalIn,
		TotalOut: totals.TotalOut,
		RateIn:   totals.RateIn,
		RateOut:  totals.RateOut,
	}
}

func (s *ipfsStore) RepoStats() *RepoStats {
	stats, err := corerepo.RepoStat(s.node.Context(), s.node)
	if err != nil {
		log.Warningf("failed to read core stats: %v", err)
		return nil
	}
	return &RepoStats{
		NumObjects: stats.NumObjects,
		RepoSize:   stats.RepoSize,
		RepoPath:   stats.RepoPath,
		Version:    stats.Version,
		StorageMax: stats.StorageMax,
	}
}

func (s *ipfsStore) BitswapStats() *BitswapStats {
	b, ok := s.node.Exchange.(*bitswap.Bitswap)
	if !ok {
		return nil
	}
	stats, err := b.Stat()
	if err != nil {
		log.Warningf("failed to read bitswap stats: %v", err)
		return nil
	}
	return &BitswapStats{
		ProvideBufLen:   stats.ProvideBufLen,
		WantlistLen:     len(stats.Wantlist),
		Peers:           stats.Peers,
		BlocksReceived:  stats.BlocksReceived,
		DataReceived:    stats.DataReceived,
		BlocksSent:      stats.BlocksSent,
		DataSent:        stats.DataSent,
		DupBlksReceived: stats.DupBlksReceived,
		DupDataReceived: stats.DupDataReceived,
	}
}

func (s *ipfsStore) SignData(nodeID string, data []byte) ([]byte, error) {
	pk := s.node.PrivateKey.GetPublic()
	id, err := peer.IDFromEd25519PublicKey(pk)
	if err != nil {
		return nil, err
	}
	if nodeID != id.Pretty() {
		return nil, errors.New("sign nodeID mismatch")
	}
	return s.node.PrivateKey.Sign(data)
}

func (s *ipfsStore) VerifyNode(code string) (string, error) {
	id, err := peer.IDFromPrivateKey(s.node.PrivateKey)
	if err != nil {
		return "", err
	}
	signed, err := s.node.PrivateKey.Sign([]byte(code))
	if err != nil {
		return "", err
	}
	jsoned := fmt.Sprintf("[%q, %q]", hex.EncodeToString(signed), id.Pretty())
	return base64.StdEncoding.EncodeToString([]byte(jsoned)), nil
}

// VerifyDataSignature - checks the signature for the node
func VerifyDataSignature(nodeID string, sig string, data []byte) (bool, error) {
	id, err := peer.IDB58Decode(nodeID)
	if err != nil {
		return false, err
	}
	pk, err := id.ExtractEd25519PublicKey()
	if err != nil {
		return false, err
	}
	return pk.Verify(data, []byte(sig))
}

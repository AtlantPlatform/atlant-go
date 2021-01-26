// Copyright 2017-2021 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD-3-Clause "New" or "Revised"
// License (BSD-3-Clause) that can be found in the LICENSE file.

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"math"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/AtlantPlatform/atlant-go/contracts"
	"github.com/AtlantPlatform/atlant-go/fs"
	"github.com/AtlantPlatform/atlant-go/proto"
	"github.com/AtlantPlatform/atlant-go/rs"
)

// PublicServer to contain server and time of the start
type PublicServer struct {
	mux       *gin.Engine
	startedAt time.Time
}

// NewPublicServer is a constructor of the PublicServer
func NewPublicServer() *PublicServer {
	return &PublicServer{
		startedAt: time.Now(),
	}
}

// ListenAndServe starts server binded to address i.e. "0.0.0.0:33780"
func (p *PublicServer) ListenAndServe(addr string) error {
	return p.mux.Run(addr)
}

// RouteAPI initializes GIN routes
func (p *PublicServer) RouteAPI(ctx APIContext) {
	r := gin.Default()
	r.POST("/api/v1/put/*path", p.PutHandler(ctx))
	r.POST("/api/v1/delete/:id", p.DeleteHandler(ctx))
	r.GET("/api/v1/content/*path", p.ContentHandler(ctx))
	r.GET("/api/v1/meta/*path", p.MetaHandler(ctx))
	r.GET("/api/v1/listVersions/*path", p.ListVersionsHandler(ctx))
	r.GET("/api/v1/listAll/*prefix", p.ListAllHandler(ctx))

	r.GET("/api/v1/tokenDistributionInfo", p.TokenDistributionInfo(ctx))
	r.GET("/api/v1/kycStatus", p.KYCStatus(ctx))
	r.GET("/api/v1/ethBalance", p.TokenBalance(ctx, contracts.TokenETH))
	r.GET("/api/v1/atlBalance", p.TokenBalance(ctx, contracts.TokenATL))
	r.GET("/api/v1/ptoBalance/:token", p.PropertyTokenBalance(ctx))

	r.GET("/api/v1/newID", p.IDHandler(ctx))
	r.GET("/api/v1/ping", p.PingHandler(ctx))
	r.GET("/api/v1/env", p.EnvHandler(ctx))
	r.GET("/api/v1/session", p.SessionHandler(ctx))
	r.GET("/api/v1/version", p.VersionHandler(ctx))
	r.GET("/api/v1/stats", p.StatsHandler(ctx))
	r.GET("/api/v1/logs", p.LogListHandler(ctx))
	r.GET("/api/v1/log/:year/:month/:day", p.LogGetHandler(ctx))

	r.GET("/index/*prefix", p.IndexHandler(ctx))
	r.StaticFS("/assets", assetFS())

	p.mux = r
}

// PingHandler returns HTTP response with Node ID
func (p *PublicServer) PingHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, ctx.NodeID())
	}
}

// IDHandler returns HTTP response with New ID
func (p *PublicServer) IDHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, proto.NewID())
	}
}

// EnvHandler returns HTTP response with current environment
func (p *PublicServer) EnvHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, ctx.Env())
	}
}

// SessionHandler returns HTTP response with current sesion Id
func (p *PublicServer) SessionHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, ctx.SessionID())
	}
}

// VersionHandler returns HTTP response with build version
func (p *PublicServer) VersionHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, ctx.Version())
	}
}

const (
	// KB - kilobytes
	KB = 1024
	// MB - megabytes
	MB = 1024 * KB
	// GB - gigabytes
	GB = 1024 * MB
)

// DiskStats contains stats of disks usage
type DiskStats struct {
	*fs.DiskStats

	KBytesAll  float64 `json:"kb_all"`
	KBytesUsed float64 `json:"kb_used"`
	KBytesFree float64 `json:"kb_free"`

	MBytesAll  float64 `json:"mb_all"`
	MBytesUsed float64 `json:"mb_used"`
	MBytesFree float64 `json:"mb_free"`

	GBytesAll  float64 `json:"gb_all"`
	GBytesUsed float64 `json:"gb_used"`
	GBytesFree float64 `json:"gb_free"`
}

// Stats is container for
type Stats struct {
	Uptime         string             `json:"uptime,omitempty"`
	DiskStats      *DiskStats         `json:"disk_stats,omitempty"`
	BandwidthStats *fs.BandwidthStats `json:"bandwidth_stats,omitempty"`
	RepoStats      *fs.RepoStats      `json:"repo_stats,omitempty"`
	BitswapStats   *fs.BitswapStats   `json:"bitswap_stats,omitempty"`
	BadgerStats    *rs.BadgerStats    `json:"badger_stats,omitempty"`
}

// StatsHandler endpoint returning JSON with all collects stats
func (p *PublicServer) StatsHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		stats := &Stats{
			Uptime:         fmt.Sprintf("%s", time.Since(p.startedAt)),
			BandwidthStats: ctx.FileStore().BandwidthStats(),
			RepoStats:      ctx.FileStore().RepoStats(),
			BadgerStats:    ctx.RecordStore().BadgerStats(),
		}
		if useBitswap := c.Query("bitswap"); useBitswap == "1" || useBitswap == "true" {
			stats.BitswapStats = ctx.FileStore().BitswapStats()
		}
		if ds, err := ctx.FileStore().DiskStats(); err == nil {
			stats.DiskStats = &DiskStats{
				DiskStats: ds,
			}
			stats.DiskStats.KBytesAll = float64(ds.BytesAll) / KB
			stats.DiskStats.KBytesUsed = float64(ds.BytesUsed) / KB
			stats.DiskStats.KBytesFree = float64(ds.BytesFree) / KB
			stats.DiskStats.MBytesAll = float64(ds.BytesAll) / MB
			stats.DiskStats.MBytesUsed = float64(ds.BytesUsed) / MB
			stats.DiskStats.MBytesFree = float64(ds.BytesFree) / MB
			stats.DiskStats.GBytesAll = float64(ds.BytesAll) / GB
			stats.DiskStats.GBytesUsed = float64(ds.BytesUsed) / GB
			stats.DiskStats.GBytesFree = float64(ds.BytesFree) / GB
		}
		c.JSON(200, stats)
	}
}

// ContentHandler is endpoint to return content
func (p *PublicServer) ContentHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, err := ctx.RecordStore().ReadRecord(ctx, c.Param("path"), rs.ReadOptions{
			Version: c.Query("ver"),
		})
		if err == rs.ErrRecordNotFound {
			if r != nil {
				if meta := r.Object.Meta(); meta != nil {
					serveMeta(c, meta)
					c.Status(404)
					return
				}
			}
			c.AbortWithStatus(404)
			return
		} else if err != nil {
			c.String(500, "error: %v", err)
			return
		}
		serveObject(c, r.Body, r.Object.Meta())
	}
}

// MetaHandler is endpoint to get meta information on the record
func (p *PublicServer) MetaHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, err := ctx.RecordStore().ReadRecord(ctx, c.Param("path"), rs.ReadOptions{
			Version:   c.Query("ver"),
			NoContent: true,
		})
		if err == rs.ErrRecordNotFound {
			if r != nil {
				c.JSON(200, r.Object.Meta())
				return
			}
			c.AbortWithStatus(404)
			return
		} else if err != nil {
			c.String(500, "error: %v", err)
			return
		}
		c.JSON(200, r.Object.Meta())
	}
}

// PutHandler is endpoint to set Meta information
func (p *PublicServer) PutHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		size, _ := strconv.ParseInt(c.Request.Header.Get("Content-Length"), 10, 64)
		userMeta := c.Request.Header.Get("X-Meta-UserMeta")
		if len(userMeta) > 0 {
			if !json.Valid([]byte(userMeta)) {
				c.String(400, "error: user meta json is not valid: %s", userMeta)
				return
			}
		}
		path := c.Param("path")
		if len(path) == 0 || path == "/" || len(filepath.Base(path)) == 0 {
			c.AbortWithStatus(400)
			return
		}
		r, err := ctx.RecordStore().CreateRecord(ctx, path, c.Request.Body, rs.CreateOptions{
			Size:     size,
			UserMeta: []byte(userMeta),
		})
		if err == rs.ErrRecordExists {
			log.Debugln("record exists, updating:", path)
			r, err = ctx.RecordStore().UpdateRecord(ctx, path, c.Request.Body, rs.UpdateOptions{
				Size:     size,
				UserMeta: []byte(userMeta),
			})
		} else if err == nil {
			log.Debugln("record not exists, created:", path, r.Id())
		}
		if err != nil {
			log.WithFields(log.Fields{
				"path": path,
				"id":   r.Id(),
			}).Errorf("error: %v", err)
			c.String(500, "error: %v", err)
			return
		}
		c.JSON(200, r.Object.Meta())
	}
}

// DeleteHandler endpoint to delete record from the store
func (p *PublicServer) DeleteHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, err := ctx.RecordStore().DeleteRecord(ctx, c.Param("id"))
		if err == rs.ErrRecordNotFound {
			if r != nil {
				if meta := r.Object.Meta(); meta != nil {
					serveMeta(c, meta)
					c.Status(200)
					return
				}
			}
			c.Status(404)
			return
		} else if err != nil {
			c.String(500, "error: %v", err)
			return
		}
		if meta := r.Object.Meta(); meta != nil {
			serveMeta(c, meta)
		}
		c.Status(200)
	}
}

// LogListHandler endpoint to return list of available logs
func (p *PublicServer) LogListHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		dir := ctx.LogDir()
		if len(dir) == 0 {
			return
		}
		var list []string
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			} else if info.IsDir() && path != dir {
				return filepath.SkipDir
			} else if filepath.Ext(path) != ".log" {
				return nil
			}
			list = append(list, filepath.Base(path))
			return nil
		})
		c.JSON(200, list)
	}
}

// LogGetHandler endpoint
func (p *PublicServer) LogGetHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		dir := ctx.LogDir()
		if len(dir) == 0 {
			c.Status(404)
			return
		}
		year := numeric(c.Param("year"))
		month := numeric(c.Param("month"))
		day := numeric(c.Param("day"))
		logFile := filepath.ToSlash(filepath.Join(dir, fmt.Sprintf("%s-%s-%s.log", year, month, day)))
		c.File(logFile)
	}
}

// DistributionInfo to store beat report
type DistributionInfo struct {
	Report     *rs.BeatReport `json:"report"`
	HoursTotal uint64         `json:"hours_total"`
}

// TokenDistributionInfo endpoint to show token distribution
func (p *PublicServer) TokenDistributionInfo(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountAddr := strings.ToLower(c.Query("account"))
		if len(accountAddr) == 0 {
			accountAddr = ctx.ETHAddr()
			if len(accountAddr) == 0 {
				c.String(400, "error: no ETH account specified")
			}
		}
		var report *rs.BeatReport
		r, err := ctx.RecordStore().ReadRecord(ctx, fmt.Sprintf("/beat_reports/%s.json", accountAddr))
		if err == rs.ErrRecordNotFound {
			c.JSON(200, &DistributionInfo{})
			return
		} else if err != nil {
			c.String(500, "error: %v", err)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
			c.String(500, "error: %v", err)
			return
		}
		r.Body.Close()
		var hours uint64
		for _, sess := range report.Sessions {
			hours += uint64(sess.Uptime)
		}
		c.JSON(200, &DistributionInfo{
			Report:     report,
			HoursTotal: hours,
		})
	}
}

// KYCStatus is endpoint to return KYC account status
func (p *PublicServer) KYCStatus(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountAddr := strings.ToLower(c.Query("account"))
		if len(accountAddr) == 0 {
			accountAddr = ctx.ETHAddr()
			if len(accountAddr) == 0 {
				c.String(400, "error: no ETH account specified")
			}
		}
		mgr, err := ctx.ContractsManager().KYCManager()
		if err != nil {
			c.String(500, "error: %v", err)
			return
		}
		status, err := mgr.AccountStatus(accountAddr)
		if err != nil {
			c.String(500, "error: %v", err)
			return
		}
		c.String(200, "%s", status)
	}
}

// TokenBalance - endpoint to return Token Balance for the wallet (plain text)
func (p *PublicServer) TokenBalance(ctx APIContext, token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountAddr := strings.ToLower(c.Query("account"))
		if len(accountAddr) == 0 {
			accountAddr = ctx.ETHAddr()
			if len(accountAddr) == 0 {
				c.String(400, "error: no ETH account specified")
			}
		}
		mgr, err := ctx.ContractsManager().TokenManager(token, "")
		if err != nil {
			c.String(500, "error: %v", err)
			return
		}
		balance, err := mgr.AccountBalance(accountAddr)
		if err != nil {
			c.String(500, "error: %v", err)
			return
		}
		c.String(200, "%f", balance)
	}
}

// PropertyTokenBalance is endpoint to respond with property token balance
func (p *PublicServer) PropertyTokenBalance(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountAddr := strings.ToLower(c.Query("account"))
		if len(accountAddr) == 0 {
			accountAddr = ctx.ETHAddr()
			if len(accountAddr) == 0 {
				c.String(400, "error: no ETH account specified")
			}
		}
		token := strings.ToLower(c.Param("token"))
		mgr, err := ctx.ContractsManager().TokenManager(contracts.TokenPTO, token)
		if err != nil {
			c.String(500, "error: %v", err)
			return
		}
		balance, err := mgr.AccountBalance(accountAddr)
		if err != nil {
			c.String(500, "error: %v", err)
			return
		}
		c.String(200, "%f", balance)
	}
}

func numeric(str string) string {
	var safe []rune
	for _, v := range str {
		if '0' <= v && v <= '9' {
			safe = append(safe, v)
		}
	}
	return string(safe)
}

func serveMeta(c *gin.Context, meta *proto.ObjectMeta) {
	c.Header("X-Meta-ID", meta.Id())
	c.Header("X-Meta-Version", meta.Version())
	if ver := meta.VersionPrevious(); len(ver) > 0 {
		c.Header("X-Meta-Previous", ver)
	}
	if p := meta.Path(); len(p) > 0 {
		c.Header("X-Meta-Path", p)
	}
	if m := meta.UserMeta(); len(m) > 0 {
		c.Header("X-Meta-UserMeta", m)
	}
	if meta.IsDeleted() {
		c.Header("X-Meta-Deleted", "true")
	}
}

func serveObject(c *gin.Context, r io.ReadCloser, meta *proto.ObjectMeta) {
	serveMeta(c, meta)
	ts := time.Unix(0, meta.CreatedAt())
	if seekable, ok := r.(io.ReadSeeker); ok {
		http.ServeContent(c.Writer, c.Request, meta.Path(), ts, seekable)
		return
	}
	// actually do all the work http.ServeContent does, but without support
	// of ranges and partial reads due to lack of io.Seeker interface.
	if !ts.IsZero() {
		c.Header("Last-Modified", ts.UTC().Format(http.TimeFormat))
	}
	ctype := mime.TypeByExtension(filepath.Ext(meta.Path()))
	c.Header("Content-Type", ctype)
	if meta.Size() > 0 {
		c.Header("Content-Length", strconv.FormatInt(meta.Size(), 10))
		io.CopyN(c.Writer, r, meta.Size())
		return
	}
	io.Copy(c.Writer, r)
	return
}

// IndexTemplate is a template to generate assets as Golang code
//go:generate go-bindata-assetfs -pkg api assets/templates assets/icons
var IndexTemplate = template.Must(
	template.New("Index").Parse(
		string(MustAsset("assets/templates/index.html.tpl")),
	),
)

// Index stores the files and folder descriptor
type Index struct {
	Prefix       string
	ParentPrefix string
	Files        []*IndexFile
}

// Compile apache-like folder
func (i *Index) Compile() ([]byte, error) {
	var buf bytes.Buffer
	if err := IndexTemplate.Execute(&buf, i); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// IndexFile is directory index descriptor
type IndexFile struct {
	Dir          bool
	Name         string
	Path         string
	LastModified string
	Size         string
	UserMeta     string
	Icon         string
	IconAlt      string
}

// IndexFilesByName is array of directories index descriptor
type IndexFilesByName []*IndexFile

func (s IndexFilesByName) Len() int           { return len(s) }
func (s IndexFilesByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s IndexFilesByName) Less(i, j int) bool { return s[i].Name < s[j].Name }

var indexIcons = map[string][]string{
	".pdf":  []string{"layout.png", "[DOC]"},
	".doc":  []string{"layout.png", "[DOC]"},
	".docx": []string{"layout.png", "[DOC]"},

	".png":  []string{"image2.png", "[IMG]"},
	".gif":  []string{"image2.png", "[IMG]"},
	".jpg":  []string{"image2.png", "[IMG]"},
	".jpeg": []string{"image2.png", "[IMG]"},
	".tiff": []string{"image2.png", "[IMG]"},
	".bmp":  []string{"image2.png", "[IMG]"},

	".txt":  []string{"text.png", "[TXT]"},
	".json": []string{"text.png", "[TXT]"},
	".yml":  []string{"text.png", "[TXT]"},
	".yaml": []string{"text.png", "[TXT]"},
	".conf": []string{"text.png", "[TXT]"},
	".cfg":  []string{"text.png", "[TXT]"},
	".ini":  []string{"text.png", "[TXT]"},

	"dir":     []string{"dir.png", "[DIR]"},
	"default": []string{"generic.png", "[OBJ]"},
}

// ListVersionsResponse contains list of versions of the object
type ListVersionsResponse struct {
	ID       string              `json:"id"`
	Versions []*proto.ObjectMeta `json:"versions"`
}

// ListVersionsHandler - endpoint to return all versions for the object under given path
func (p *PublicServer) ListVersionsHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var versions []*proto.ObjectMeta
		r, err := ctx.RecordStore().ReadRecord(ctx, c.Param("path"), rs.ReadOptions{
			NoContent: true,
		})
		if err == rs.ErrRecordNotFound {
			if r == nil {
				c.AbortWithStatus(404)
				return
			}
		} else if err != nil {
			c.String(500, "error: %v", err)
			return
		}
		versions = append(versions, r.Object.Meta())
		limit := r.Previous().Len()
		for i := 0; i < limit; i++ {
			r, err := ctx.RecordStore().ReadRecord(ctx, "", rs.ReadOptions{
				Version:   r.Previous().At(i).Version(),
				NoContent: true,
			})
			if err == rs.ErrRecordNotFound {
				if r == nil {
					continue
				}
			} else if err != nil {
				log.Warningf("failed to read record from store: %v", err)
				continue
			}
			versions = append(versions, r.Object.Meta())
		}
		c.JSON(200, &ListVersionsResponse{
			ID:       r.Id(),
			Versions: versions,
		})
	}
}

// ListResponse structure of response for the client = list of Dirs and Files
type ListResponse struct {
	Dirs  []string
	Files []*proto.ObjectMeta
}

// ObjectMetas is array of ObjectMeta
type ObjectMetas []*proto.ObjectMeta

func (s ObjectMetas) Len() int           { return len(s) }
func (s ObjectMetas) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ObjectMetas) Less(i, j int) bool { return s[i].Path() < s[j].Path() }

// ListAllHandler - endpoint to response with all current versions under provided path
func (p *PublicServer) ListAllHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		prefix := c.Param("prefix")
		if !strings.HasPrefix(prefix, "/") {
			prefix = "/" + prefix
		}
		if !strings.HasSuffix(prefix, "/") {
			prefix = prefix + "/"
		}
		resp := &ListResponse{}
		seenDirs := make(map[string]struct{})
		err := ctx.RecordStore().WalkRecords(ctx, "", func(path string, r *rs.Record) error {
			if len(path) == 0 {
				return nil
			} else if !strings.HasPrefix(path, prefix) {
				return nil
			}
			path = strings.TrimPrefix(path, prefix)
			parts := strings.Split(path, "/")
			if len(parts) > 1 {
				dir := parts[0]
				if _, ok := seenDirs[dir]; ok {
					return nil
				}
				seenDirs[dir] = struct{}{}
				resp.Dirs = append(resp.Dirs, filepath.ToSlash(filepath.Join(prefix, dir)+"/"))
				return nil
			}
			var meta *proto.ObjectMeta
			if metaRecord, err := ctx.RecordStore().ReadRecord(ctx, r.Path(), rs.ReadOptions{
				Version:   r.Current().Version(),
				NoContent: true,
			}); err == rs.ErrRecordNotFound {
				return nil
			} else if err != nil {
				log.Warningf("failed to fetch record: %v", err)
				return nil
			} else {
				meta = metaRecord.Object.Meta()
			}
			resp.Files = append(resp.Files, meta)

			return nil
		})
		if err == rs.ErrRecordNotFound || len(resp.Files)+len(resp.Dirs) == 0 {
			c.Status(404)
			return
		} else if err != nil {
			c.String(500, "error: %v", err)
			return
		}

		sort.Sort(sort.StringSlice(resp.Dirs))
		sort.Sort(ObjectMetas(resp.Files))
		c.JSON(200, resp)
	}
}

// IndexHandler endpoint to response with all contents
func (p *PublicServer) IndexHandler(ctx APIContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		prefix := c.Param("prefix")
		if !strings.HasPrefix(prefix, "/") {
			prefix = "/" + prefix
		}
		if !strings.HasSuffix(prefix, "/") {
			prefix = prefix + "/"
		}
		index := Index{
			Prefix: prefix,
		}
		if prefix != "/" {
			index.ParentPrefix = filepath.ToSlash(filepath.Dir(filepath.Dir(prefix)))
		}
		log.WithFields(log.Fields{
			"prefix": prefix,
		}).Debug("Index: walking record")

		seenDirs := make(map[string]struct{})
		err := ctx.RecordStore().WalkRecords(ctx, "", func(path string, r *rs.Record) error {
			if len(path) == 0 {
				return nil
			} else if !strings.HasPrefix(path, prefix) {
				return nil
			}
			log.WithField("path", path).Debug("Index: walking record")
			path = strings.TrimPrefix(path, prefix)
			parts := strings.Split(path, "/")
			if len(parts) > 1 {
				dir := parts[0]
				if _, ok := seenDirs[dir]; ok {
					return nil
				}
				seenDirs[dir] = struct{}{}
				index.Files = append(index.Files, &IndexFile{
					Dir:     true,
					Name:    dir,
					Path:    filepath.ToSlash(filepath.Join(prefix, dir)) + "/",
					Icon:    indexIcons["dir"][0],
					IconAlt: indexIcons["dir"][1],
				})
				return nil
			}
			var meta *proto.ObjectMeta
			if metaRecord, err := ctx.RecordStore().ReadRecord(ctx, r.Path(), rs.ReadOptions{
				Version:   r.Current().Version(),
				NoContent: true,
			}); err == rs.ErrRecordNotFound {
				return nil
			} else if err != nil {
				log.WithFields(log.Fields{
					"prefix":  prefix,
					"path":    path,
					"version": r.Current().Version(),
				}).Warningf("Index failed to fetch record: %v", err)
				return nil
			} else {
				meta = metaRecord.Object.Meta()
			}
			f := &IndexFile{
				Name:         parts[0],
				Path:         r.Path(),
				LastModified: time.Unix(0, meta.CreatedAt()).Format(time.RFC1123),
				Size:         humanBytes(meta.Size(), 1024),
				UserMeta:     meta.UserMeta(),
			}
			ext := strings.ToLower(filepath.Ext(path))
			if icon, ok := indexIcons[ext]; ok {
				f.Icon = icon[0]
				f.IconAlt = icon[1]
			} else {
				f.Icon = indexIcons["default"][0]
				f.IconAlt = indexIcons["default"][1]
			}
			index.Files = append(index.Files, f)
			return nil
		})
		if err == rs.ErrRecordNotFound || len(index.Files) == 0 {

			// for index page, prefix is "/"
			log.WithFields(log.Fields{
				"prefix": prefix,
			}).Debug("Index: Nothing was found")
			c.Status(404)
			return
		} else if err != nil {
			log.WithFields(log.Fields{
				"prefix": prefix,
			}).Errorf("Index error: %v", err)
			c.String(500, "error: %v", err)
			return
		}
		data, err := index.Compile()
		if err != nil {
			log.WithFields(log.Fields{
				"prefix":       prefix,
				"parentPrefix": index.ParentPrefix,
			}).Errorf("Index compilation error: %v", err)
			c.String(500, "error: %v", err)
			return
		}
		c.Data(200, "text/html", data)
	}
}

var sizes = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanBytes(s int64, base float64) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f %s"
	if val < 10 {
		f = "%.1f %s"
	}

	return fmt.Sprintf(f, val, suffix)
}

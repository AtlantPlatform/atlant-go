// Copyright 2017-2019 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package fs

import (
	"context"
	"strings"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
)

// RefWriter implementation from go-ipfs/commands/refs.go
type RefWriter struct {
	out chan string
	DAG ipld.DAGService
	Ctx context.Context

	Unique    bool
	Recursive bool
	PrintFmt  string

	seen *cid.Set
}

// WriteRefs writes refs of the given object to the underlying writer.
func (rw *RefWriter) WriteRefs(n ipld.Node) (int, error) {
	if rw.Recursive {
		return rw.writeRefsRecursive(n)
	}
	return rw.writeRefsSingle(n)
}

func (rw *RefWriter) writeRefsRecursive(n ipld.Node) (int, error) {
	nc := n.Cid()

	var count int
	for i, ng := range ipld.GetDAG(rw.Ctx, rw.DAG, n) {
		lc := n.Links()[i].Cid
		if rw.skip(&lc) {
			continue
		}

		if err := rw.WriteEdge(&nc, &lc, n.Links()[i].Name); err != nil {
			return count, err
		}

		nd, err := ng.Get(rw.Ctx)
		if err != nil {
			return count, err
		}

		c, err := rw.writeRefsRecursive(nd)
		count += c
		if err != nil {
			return count, err
		}
	}
	return count, nil
}

func (rw *RefWriter) writeRefsSingle(n ipld.Node) (int, error) {
	c := n.Cid()

	if rw.skip(&c) {
		return 0, nil
	}

	count := 0
	for _, l := range n.Links() {
		lc := l.Cid
		if rw.skip(&lc) {
			continue
		}
		if err := rw.WriteEdge(&c, &lc, l.Name); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

// skip returns whether to skip a cid
func (rw *RefWriter) skip(c *cid.Cid) bool {
	if !rw.Unique {
		return false
	}

	if rw.seen == nil {
		rw.seen = cid.NewSet()
	}

	has := rw.seen.Has(*c)
	if !has {
		rw.seen.Add(*c)
	}
	return has
}

// Write one edge
func (rw *RefWriter) WriteEdge(from, to *cid.Cid, linkname string) error {
	if rw.Ctx != nil {
		select {
		case <-rw.Ctx.Done(): // just in case.
			return rw.Ctx.Err()
		default:
		}
	}

	var s string
	switch {
	case rw.PrintFmt != "":
		s = rw.PrintFmt
		s = strings.Replace(s, "<src>", from.String(), -1)
		s = strings.Replace(s, "<dst>", to.String(), -1)
		s = strings.Replace(s, "<linkname>", linkname, -1)
	default:
		s += to.String()
	}

	rw.out <- s
	return nil
}

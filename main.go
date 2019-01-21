package main

import (
	"flag"
	"log"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/jstaf/onedriver/graph"
)

var auth graph.Auth

type fuseFs struct {
	pathfs.FileSystem
}

// these files will never exist, and we should ignore them
func ignore(path string) bool {
	ignoredFiles := []string{
		"/BDMV",
		"/.Trash",
		"/.Trash-1000",
		"/.xdg-volume-info",
		"/autorun.inf",
	}
	for _, ignore := range ignoredFiles {
		if path == ignore {
			return true
		}
	}
	return false
}

func (fs *fuseFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	name = "/" + name
	if ignore(name) {
		return nil, fuse.ENOENT
	}
	log.Printf("GetAttr(\"%s\")", name)

	return nil, fuse.ENOENT
}

func (fs *fuseFs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	name = "/" + name
	log.Printf("OpenDir(\"%s\")", name)
	return nil, fuse.ENOENT
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("Usage:\n  onedriver <mountPoint>")
	}

	auth = graph.Authenticate()

	fs := pathfs.NewPathNodeFs(&fuseFs{FileSystem: pathfs.NewDefaultFileSystem()}, nil)
	server, _, err := nodefs.MountRoot(flag.Arg(0), fs.Root(), nil)
	if err != nil {
		log.Fatalf("Mount failed: %v\n", err)
	}
	server.Serve()
}
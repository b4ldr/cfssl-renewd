package request

import (
	"log"
	"os"
	"os/user"
	"strconv"
)

type Perms struct {
	Owner string `yaml:"owner"`
	Group string `yaml:"group"`
	Mode  int    `yaml:"mode"`
}

func (p *Perms) gid() int {
	if p.Group != "" {
		group, err := user.LookupGroup(p.Group)
		if err != nil {
			log.Fatal(err)
		}
		gid, err := strconv.Atoi(group.Gid)
		if err != nil {
			log.Fatal(err)
		}
		return gid
	}
	return os.Getegid()
}
func (p *Perms) uid() int {
	if p.Owner != "" {
		owner, err := user.Lookup(p.Owner)
		if err != nil {
			log.Fatal(err)
		}
		uid, err := strconv.Atoi(owner.Uid)
		if err != nil {
			log.Fatal(err)
		}
		return uid
	}
	return os.Geteuid()
}

func (p *Perms) mode() os.FileMode {
	if p.Mode != 0 {
		return os.FileMode(p.Mode)
	}
	return 0400
}

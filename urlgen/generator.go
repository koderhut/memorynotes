package urlgen

import (
	"fmt"

	"github.com/koderhut/memorynotes/config"
)

type Url struct {
	Scheme string
	Host   string
	Prefix string
}

func (u *Url) Generate(path string) string {
	return fmt.Sprintf("%s%s%s%s", u.Scheme, u.Host, u.Prefix, path)
}

func FromConfig(c config.Context) *Url {
	return &Url{
		Scheme: "http://",
		Host:   c.Domain,
		Prefix: c.PathPrefix,
	}
}

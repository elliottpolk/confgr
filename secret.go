package peppermintsparkles

import (
	"encoding/json"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type AppInfo struct {
	Name string `json:"app_name"`
	Env  string `json:"env"`
}

type Meta struct {
	Team    string `json:"team,omitempty"`
	User    string `json:"username"`
	Expires int64  `json:"expires,omitempty"`
}

type Secret struct {
	Id       string      `json:"id,omitempty"`
	Content  interface{} `json:"content"`
	App      *AppInfo    `json:"app_info"`
	Metadata *Meta       `json:"metadata"`
}

func (s *Secret) String() string {
	out, err := json.Marshal(s)
	if err != nil {
		log.Error(errors.Wrap(err, "unable to marshal secret to string"))
	}

	return string(out)
}

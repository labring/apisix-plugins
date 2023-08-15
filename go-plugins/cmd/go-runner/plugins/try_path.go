package plugins

import (
	"encoding/json"
	"errors"
	pkgHTTP "github.com/apache/apisix-go-plugin-runner/pkg/http"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/apache/apisix-go-plugin-runner/pkg/plugin"
	"net/http"
	"strings"
	"time"
)

func init() {
	err := plugin.RegisterPlugin(&TryPath{})
	if err != nil {
		log.Fatalf("failed to register plugin try_files: %s", err)
	}
}

type TryPath struct {
	plugin.DefaultPlugin
}

func (t *TryPath) Name() string {
	return "try-path"
}

func (t *TryPath) ParseConf(in []byte) (interface{}, error) {
	cfg := TryPathConf{}
	err := json.Unmarshal(in, &cfg)
	if err != nil {
		return nil, err
	}
	if len(cfg.Paths) < 2 && len(cfg.Paths) <= 4 {
		return nil, errors.New("uris must be more than 1 but less than or equal to 4")
	}
	return cfg, err
}

type TryPathConf struct {
	Paths []string `json:"paths"`
	Host  string   `json:"host"`
}

func (t *TryPath) RequestFilter(conf interface{}, w http.ResponseWriter, r pkgHTTP.Request) {
	cfg := conf.(TryPathConf)

	// get path length
	pathLen := len(cfg.Paths)

	// loop path
	for index := 0; index < pathLen; index++ {
		newPath := strings.ReplaceAll(cfg.Paths[index], "$uri", string(r.Path()))

		// if the last path, directly set path
		if index == pathLen-1 {
			log.Warnf("redirect path: %s", newPath)
			r.SetPath([]byte(newPath))
			return
		}

		// try path exists
		st := time.Now()
		ok, err := tryPath(cfg.Host, newPath)
		log.Warnf("try path cost: %d", time.Now().Sub(st).Milliseconds())
		if err != nil {
			log.Errorf("try path failed: %s", err)
			return
		}
		if ok {
			log.Warnf("redirect path: %s", newPath)
			r.SetPath([]byte(newPath))
			return
		}
	}
}

func tryPath(host, path string) (bool, error) {
	log.Warnf("try path: %s", host+path)
	resp, err := http.Head(host + path)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 || resp.StatusCode == 403 || resp.StatusCode == 400 || resp.Header.Get("Content-Type") == "folder" {
		return false, nil
	}
	return true, nil
}

package trader

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/robertkrimen/otto/registry"
)

var (
	scripts = []string{}
	entyr   = registry.Register(func() string {
		return strings.Join(scripts, "")
	})
)

func init() {
	filepath.Walk("plugin", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".js") {
			return err
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		data, _ := ioutil.ReadAll(file)
		scripts = append(scripts, string(data))
		return nil
	})
}

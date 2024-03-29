package yaml

import (
	"fmt"
	"github.com/croissong/releasechecker/pkg/provider"
	"github.com/croissong/releasechecker/pkg/util/cmd"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"net/http"
)

type Yaml struct {
	config *config
}

type config struct {
	Path string
	Url  string
}

func (_ Yaml) NewProvider(conf map[string]interface{}) (provider.Provider, error) {
	config, err := validateConfig(conf)
	if err != nil {
		return nil, err
	}
	yaml := Yaml{config: config}
	return &yaml, nil
}

func (yaml Yaml) GetVersion() (string, error) {
	resp, err := http.Get(yaml.config.Url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(body)
	version, err := execYq(yaml.config.Path, bodyString)
	if err != nil {
		return "", err
	}
	return version, nil
}

func (yaml Yaml) GetVersions() ([]string, error) {
	return nil, nil
}

func execYq(path string, input string) (string, error) {
	command := fmt.Sprintf("yq r - %s", path)
	return cmdutil.RunCmd(command, cmdutil.CmdOptions{Input: input})
}

func validateConfig(conf map[string]interface{}) (*config, error) {
	var config config
	if err := mapstructure.Decode(conf, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

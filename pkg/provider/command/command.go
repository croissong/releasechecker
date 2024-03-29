package command

import (
	"errors"
	"fmt"
	"github.com/croissong/releasechecker/pkg/config"
	"github.com/croissong/releasechecker/pkg/log"
	"github.com/croissong/releasechecker/pkg/provider"
	"github.com/croissong/releasechecker/pkg/util/cmd"
	"github.com/mitchellh/mapstructure"
	"strings"
)

type Command struct {
	Command string
}

func (_ Command) NewProvider(config map[string]interface{}) (provider.Provider, error) {
	var command Command
	if err := mapstructure.Decode(config, &command); err != nil {
		return nil, err
	}
	if command.Command == "" {
		return nil, errors.New(fmt.Sprintf("Missing field 'command' in config"))
	}
	log.Logger.Debugf("%#v", command)
	return &command, nil
}

func (cmd Command) GetVersion() (string, error) {
	version, err := cmd.getOutput()
	if err != nil {
		if config.Config.InitDownstreams {
			log.Logger.Infof("Ignoring cmd err due to 'initSouces' set. (%s)", err)
			return "", nil
		} else {
			return "", err
		}
	}
	log.Logger.Debugf("Got source version: %s", version)
	return version, nil
}

func (cmd Command) GetVersions() ([]string, error) {
	output, err := cmd.getOutput()
	if err != nil {
		return nil, err
	}
	versions := strings.Split(output, "\n")
	log.Logger.Debug(versions)
	return versions, nil
}

func (cmd Command) getOutput() (string, error) {
	output, err := cmdutil.RunCmd(cmd.Command, cmdutil.CmdOptions{})
	if err != nil {
		errMessage := fmt.Sprintf("Command err: %s - %s", err, output)
		return "", errors.New(errMessage)
	}
	return output, nil
}

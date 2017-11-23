package common

import (
	log "github.com/Sirupsen/logrus"
	"github.com/yjiong/go_tg120/config"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	VERSION    = "1.0"
	MODEL      = "TG150"
	INTERFACES = "/etc/network/interfaces"
)

var CONFILEPATH string = "./config.ini"
var DEVFILEPATH = "./devlist.ini"
var TEMPLATE string = "./templates"
var MESSAGE string = "./message"
var Mqtt_connected bool = false

func init() {
	var pathfs string
	if runtime.GOOS == "Linux" {
		pathfs = "\\"
	} else {
		pathfs = "/"
	}
	if execfile, err := exec.LookPath(os.Args[0]); err == nil {
		//		fmt.Printf("%s\n", execfile)
		if path, err := filepath.Abs(execfile); err == nil {
			//			fmt.Printf("%s\n", path)
			i := strings.LastIndex(path, pathfs)
			basepath := string(path[0 : i+1])
			//			fmt.Printf("%s\n", path[0:i+1])
			CONFILEPATH = basepath + CONFILEPATH[2:]
			DEVFILEPATH = basepath + DEVFILEPATH[2:]
			TEMPLATE = basepath + TEMPLATE[2:]
			MESSAGE = basepath + MESSAGE[2:]
		}
	}
}

func NewConMap(confile string) (map[string]string, error) {
	_, err := os.Stat(confile)
	if os.IsNotExist(err) {
		return nil, nil
	}
	con, err := config.LoadConfigFile(confile)
	if err != nil {
		log.WithFields(log.Fields{
			"config": con,
		}).Errorf("load config file failed: %s", err)
		return nil, err
	}
	retm := make(map[string]string)
	for _, sec := range con.GetSectionList() {
		if m, err := con.GetSection(sec); err == nil {
			retm = Mergemap(retm, m)
		} else {
			log.Errorf("get config element failed: %s", err)
			return nil, err
		}
	}
	return retm, nil
}

func Mergemap(lm ...map[string]string) map[string]string {
	retmap := make(map[string]string)
	for _, m := range lm {
		for k, v := range m {
			retmap[k] = v
		}
	}
	return retmap
}

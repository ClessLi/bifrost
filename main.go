package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var confPath = flag.String("f", "", "jenkins `conf`.y(a)ml path.")

const (
	ERROR      = "ERROR"
	WARN       = "WARN"
	INFO       = "INFO"
	DEBUG      = "DEBUG"
	timeFormat = "2006-01-02 15:04:05"
)

type ParserConfig struct {
	//Name         string `json:"name"`
	//RelativePath string `json:"relative_path"`
	//Port         int    `json:"port"`
	//ConfPath     string `json:"conf_path"`
	Name         string `yaml:"name"`
	RelativePath string `yaml:"relativePath"`
	Port         int    `yaml:"port"`
	ConfPath     string `yaml:"confPath"`
	NginxBin     string `yaml:"nginxBin"`
}

type ParserConfigs struct {
	//Configs []ParserConfig `json:"configs"`
	Configs []ParserConfig `yaml:"configs"`
}

func main() {
	//flag.Parse()
	//if *confPath == "" {
	//	*confPath = "./configs/ng-conf-info.yml"
	//}

	//confPath := "./configs/ng-conf-info.json"
	confPath := "./configs/ng-conf-info.yml"
	//isExist, pathErr := PathExists(*confPath)
	isExist, pathErr := PathExists(confPath)
	if !isExist {
		if pathErr != nil {
			fmt.Println("The config file", confPath, "is not found.")
		} else {
			fmt.Println("Unkown error of the config file.")
		}
		//flag.Usage()
		os.Exit(1)
	}
	//confData,readErr := readFile(*confPath)
	confData, readErr := readFile(confPath)
	if readErr != nil {
		fmt.Println(readErr)
		os.Exit(1)
	}

	configs := &ParserConfigs{}
	//jsonErr := json.Unmarshal(confData, configs)
	jsonErr := yaml.Unmarshal(confData, configs)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		os.Exit(1)
	}
	for _, config := range configs.Configs {
		conf, err := resolv.Load(config.ConfPath)

		if err != nil {
			fmt.Println(err)
			continue
		}

		errChan := make(chan error)

		go run(conf, config.RelativePath, config.Port, config.NginxBin, errChan)

		err = <-errChan
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	//nginxConfPath := "test/tmp/nginx.conf"
	//nginxConfPath := "test/config_test/nginx.conf"
	//conf := NewConf(nil, "test.conf")
	//h := NewHttp()
	//s1 := NewServer()
	//l := NewLocation("/")
	//k := NewKey("$test","$testV")
	//k2 := NewKey("$remote_addr - $remote_user [$time_local] \"$request\" ","")
	//i := NewInclude("../../test/mime.types","../../test/mime.types")
	//l.Add(k)
	//l.Add(k2)
	//s1.Add(l)
	//h.Add(i)
	//h.Add(s1)
	//conf.Add(h)
}

func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fd, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

func run(conf *resolv.Config, relativePath string, port int, ngBin string, errChan chan error) {
	_, jerr := json.Marshal(conf)
	//confBytes, jerr := json.Marshal(conf)
	//confBytes, jerr := json.MarshalIndent(conf, "", "    ")
	if jerr != nil {
		errChan <- jerr
	}

	ngBin, absErr := filepath.Abs(ngBin)
	if absErr != nil {
		errChan <- absErr
	}

	router := gin.Default()
	router.GET(relativePath, func(c *gin.Context) {
		h := GET(conf, c)
		c.JSON(200, &h)
	})

	router.POST(relativePath, func(c *gin.Context) {
		//var confBrif string
		//confBytes := make([]byte, 1024)
		//n, _ := c.Request.Body.Read(confBytes)
		//if n > 200 {
		//	confBrif = fmt.Sprintf("%s...%s", string(confBytes[0:50]), string(confBytes[n-50:n]))
		//} else {
		//	confBrif = string(confBytes[0:n])
		//}
		//confStr := string(confBytes[0:n])
		h := POST(ngBin, conf, c)
		c.JSON(200, &h)

	})

	rErr := router.Run(fmt.Sprintf(":%d", port))
	if rErr != nil {
		errChan <- rErr
	}

	errChan <- nil

}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, err
	} else {
		return false, nil
	}
}

func log(level, message string) {
	fmt.Printf("[%s] [%s] %s\n", level, time.Now().Format(timeFormat), message)

}

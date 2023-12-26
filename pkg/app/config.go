package app

import (
	"errors"
	"github.com/BurntSushi/toml"
	jsoniter "github.com/json-iterator/go"
	"imsdk/pkg/funcs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Configuration 应用配置
type Configuration struct {
	configs map[string]map[string]interface{}
	once    sync.Once
}

var (
	config *Configuration
	json   = jsoniter.Config{EscapeHTML: true, TagKey: "toml"}.Froze()
	// ErrNodeNotExists 配置节点不存在
	ErrNodeNotExists = errors.New("node not exists")
)

func LoadConfig() {
	config = new(Configuration).singleLoad()

}

// Config 得到config对象
func Config() *Configuration {
	return config
}

func (conf *Configuration) copy(node string, value map[string]interface{}) {
	for key, val := range value {
		if conf.configs[node] == nil {
			conf.configs[node] = make(map[string]interface{})
		}
		conf.configs[node][key] = val
	}
}

func (conf *Configuration) walk(path string, info os.FileInfo, err error) error {
	if err == nil {
		if !info.IsDir() {
			if !strings.HasSuffix(path, ".toml") {
				return nil
			}
			var err error
			var config map[string]interface{}
			_, err = toml.DecodeFile(path, &config)
			if err != nil {
				// 配置读失败了
				log.Fatal(err)
			}
			conf.copy(strings.TrimSuffix(info.Name(), ".toml"), config)
		} else {
			return filepath.Walk(info.Name(), conf.walk)
		}
	}
	return nil
}

func (conf *Configuration) GetPublicConfigDir() string {
	rootDir := funcs.GetRoot()
	pubConfDir := rootDir + "/config/"
	return pubConfDir
}

func (conf *Configuration) GetConfigDirs() []string {
	pubConfDir := conf.GetPublicConfigDir()
	rootDir := funcs.GetRoot()
	confDirArr := []string{pubConfDir, rootDir + "/configs/" + "global/"}
	if runModule := os.Getenv("TMM_RUN_MODULE"); runModule != "" {
		confDirArr = append(confDirArr, rootDir+"/configs/"+runModule+"/")
	}
	return confDirArr
}

func (conf *Configuration) singleLoad() *Configuration {
	conf.once.Do(func() {
		conf.configs = make(map[string]map[string]interface{})
		confDirArr := conf.GetConfigDirs()
		for _, dir := range confDirArr {
			rd, err := ioutil.ReadDir(dir)
			if err != nil {
				continue
			}
			for _, fi := range rd {
				if !fi.IsDir() {
					if strings.HasSuffix(dir+fi.Name(), ".toml") {
						var config map[string]interface{}
						if _, err = toml.DecodeFile(dir+fi.Name(), &config); err != nil {
							// 配置读失败了
							log.Fatal(err)
						}
						conf.copy(strings.TrimSuffix(fi.Name(), ".toml"), config)
					}
				}
			}
		}
	})
	return conf
}

// Bind 将配置绑定到传入对象
//  node 其实是配置文件的文件名
//  key 是配置文件中的顶层key
//  具体可查看该方法的其他包的使用
func (conf *Configuration) Bind(node, key string, obj interface{}) error {
	nodeVal, ok := conf.configs[node]
	if !ok {
		return nil
	}

	var objVal interface{}

	if key != "" {
		objVal, ok = nodeVal[key]
		if !ok {
			return ErrNodeNotExists
		}
	} else {
		objVal = nodeVal
	}

	return conf.assignment(objVal, obj)
}

func (conf *Configuration) GetChildConf(node, key, subKey string) (interface{}, error) {
	var confInfo map[string]interface{}
	if err := conf.Bind(node, key, &confInfo); err != nil {
		return nil, err
	}
	val, ok := confInfo[subKey]
	if !ok {
		return nil, ErrNodeNotExists
	}
	return val, nil
}

func (conf *Configuration) assignment(val, obj interface{}) error {
	data, _ := json.Marshal(val)
	return json.Unmarshal(data, obj)
}

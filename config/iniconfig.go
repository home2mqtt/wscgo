package config

import (
	"log"
	"strings"

	"gopkg.in/ini.v1"
)

type iniConfigSection struct {
	*ini.Section
	id string
}

func (is *iniConfigSection) GetID() string {
	return is.id
}

func (is *iniConfigSection) FillData(d interface{}) error {
	return is.MapTo(d)
}

func (conf *WscgoConfiguration) processConfig(category string, id string, section *ini.Section) {
	is := &iniConfigSection{
		id:      id,
		Section: section,
	}
	switch category {
	case ini.DEFAULT_SECTION:
		section.MapTo(&conf.Node)
	case "mqtt":
		section.MapTo(&conf.MqttConfig)
	default:
		parser, err := GetConfigurationPartParser(category)
		if err != nil {
			log.Print(err.Error())
		} else {
			err = parser.ParseConfiguration(is, conf)
			if err != nil {
				log.Print(err.Error())
			}
		}
	}
}

// LoadConfig loads a configuration ini file
func LoadConfig(filename string) *WscgoConfiguration {
	cfg, err := ini.Load(filename)
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}

	conf := defaultConfiguration()

	for _, s := range cfg.Sections() {
		name := s.Name()
		cat := strings.Split(name, ":")
		l := len(cat)
		var category string
		var id string
		category = cat[0]
		id = ""
		if l > 1 {
			id = cat[1]
		}
		conf.processConfig(category, id, s)
	}

	return conf
}

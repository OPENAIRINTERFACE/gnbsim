// SPDX-FileCopyrightText: 2022 Great Software Laboratory Pvt. Ltd
// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0

/*
 * gNBSim Configuration Factory
 */

package factory

import (
	"fmt"
	"os"
	"sort"

	gnbctx "github.com/omec-project/gnbsim/gnodeb/context"
	"github.com/omec-project/openapi/models"
)

const (
	GNBSIM_EXPECTED_CONFIG_VERSION string = "1.0.0"
	GNBSIM_DEFAULT_CONFIG_PATH            = "./testscenario/config-default.yaml"
)

type Config struct {
	Info          *Info          `yaml:"info"`
	Configuration *Configuration `yaml:"configuration"`
	Logger        *Logger        `yaml:"logger"`
}

type Info struct {
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

type SecurityCapabilities struct {
	FgIa0 bool `yaml:"fgIa0" json:"fgIa0"`
	FgIa1 bool `yaml:"fgIa1" json:"fgIa1"`
	FgIa2 bool `yaml:"fgIa2" json:"fgIa2"`
	FgIa3 bool `yaml:"fgIa3" json:"fgIa3"`
	FgIa4 bool `yaml:"fgIa4" json:"fgIa4"`
	FgIa5 bool `yaml:"fgIa5" json:"fgIa5"`
	FgIa6 bool `yaml:"fgIa6" json:"fgIa6"`
	FgIa7 bool `yaml:"fgIa7" json:"fgIa7"`
	FgEa0 bool `yaml:"fgEa0" json:"fgEa0"`
	FgEa1 bool `yaml:"fgEa1" json:"fgEa1"`
	FgEa2 bool `yaml:"fgEa2" json:"fgEa2"`
	FgEa3 bool `yaml:"fgEa3" json:"fgEa3"`
	FgEa4 bool `yaml:"fgEa4" json:"fgEa4"`
	FgEa5 bool `yaml:"fgEa5" json:"fgEa5"`
	FgEa6 bool `yaml:"fgEa6" json:"fgEa6"`
	FgEa7 bool `yaml:"fgEa7" json:"fgEa7"`
	Eia0  bool `yaml:"eia0" json:"Eia0"`
	Eia1  bool `yaml:"eia1" json:"Eia1"`
	Eia2  bool `yaml:"eia2" json:"Eia2"`
	Eia3  bool `yaml:"eia3" json:"Eia3"`
	Eia4  bool `yaml:"eia4" json:"Eia4"`
	Eia5  bool `yaml:"eia5" json:"Eia5"`
	Eia6  bool `yaml:"eia6" json:"Eia6"`
	Eia7  bool `yaml:"eia7" json:"Eia7"`
	Eea0  bool `yaml:"eea0" json:"Eea0"`
	Eea1  bool `yaml:"eea1" json:"Eea1"`
	Eea2  bool `yaml:"eea2" json:"Eea2"`
	Eea3  bool `yaml:"eea3" json:"Eea3"`
	Eea4  bool `yaml:"eea4" json:"Eea4"`
	Eea5  bool `yaml:"eea5" json:"Eea5"`
	Eea6  bool `yaml:"eea6" json:"Eea6"`
	Eea7  bool `yaml:"eea7" json:"Eea7"`
}

type Nas struct {
	SecurityCapabilities SecurityCapabilities `yaml:"securityCapabilities" json:"securityCapabilities"`
	SeqNum               string               `yaml:"sequenceNumber" json:"sequenceNumber"`
	Dnn                  string               `yaml:"dnn" json:"dnn"`
	SNssai               *models.Snssai       `yaml:"sNssai" json:"sNssai"`
}

type Provision struct {
	CreateSubscriber  bool   `yaml:"createSubscriber" json:"createSubscriber"`
	CreateRestUrl     string `yaml:"createRestUrl" json:"createRestUrl"`
	DeleteRestUrl     string `yaml:"deleteRestUrl" json:"deleteRestUrl"`
	CreateJsonContent string `yaml:"createJsonContent" json:"createJsonContent"`
}

type UeProfile struct {
	Model     string         `yaml:"model" json:"model"`
	StartImsi string         `yaml:"startImsi" json:"startImsi"`
	NumUes    int            `yaml:"numUes"`
	Opc       string         `yaml:"opc" json:"opc"`
	Key       string         `yaml:"key" json:"key"`
	Nas       Nas            `yaml:"nas" json:"nas"`
	Plmn      *models.PlmnId `yaml:"plmnId" json:"plmnId"`
	Provision Provision      `yaml:"provision" json:"provision"`
}

type Amf struct {
	HostName string `yaml:"hostName"`
	IpAddr   string `yaml:"ipAddr"`
	Port     string `yaml:"port"`
}

type Configuration struct {
	Amfs                    map[string]*gnbctx.GnbAmf `yaml:"amfs"`
	Gnbs                    map[string]*gnbctx.GNodeB `yaml:"gnbs"`
	UeProfiles              map[string]*UeProfile     `yaml:"ueProfiles"`
	SingleInterface         bool                      `yaml:"singleInterface"`
	ExecScenariosInParallel bool                      `yaml:"execScenariosInParallel"`
	ExecUesInParallel       bool                      `yaml:"execUesInParallel"`
	Server                  HttpServer                `yaml:"httpServer"`
	GoProfile               ProfileServer             `yaml:"goProfile"`
}

type ProfileServer struct {
	Enable bool `yaml:"enable"`
	Port   int  `yaml:"port"`
}

type HttpServer struct {
	Enable bool   `yaml:"enable"`
	IpAddr string `yaml:"ipAddr"`
	Port   string `yaml:"port"`
}

type Logger struct {
	LogLevel string `yaml:"logLevel"`
}

func (c *Config) GetVersion() string {
	if c.Info != nil && c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}

func (c *Config) Validate() (err error) {

	if c.Info == nil {
		return fmt.Errorf("Info field missing")
	}

	if c.Configuration == nil {
		return fmt.Errorf("Configuration field missing")
	}

	if len(c.Configuration.Gnbs) == 0 {
		return fmt.Errorf("no gNB(s) configured")
	}

	if len(c.Configuration.Amfs) == 0 {
		return fmt.Errorf("no AMF(s) configured")
	}

	if len(c.Configuration.UeProfiles) == 0 {
		return fmt.Errorf("no UE profile(s) configured")
	}

	if c.Configuration.GoProfile.Enable == true {
		if c.Configuration.GoProfile.Port == 0 {
			c.Configuration.GoProfile.Port = 5000
		}
	}

	if c.Configuration.Server.IpAddr == "POD_IP" {
		c.Configuration.Server.IpAddr = os.Getenv("POD_IP")
	}

	if c.Configuration.SingleInterface == true {
		for _, gnb := range c.Configuration.Gnbs {
			if gnb.GnbN3Ip == "POD_IP" {
				gnb.GnbN3Ip = os.Getenv("POD_IP")
			}
		}
	}

	return nil
}

func (c *Configuration) GetAmf(name string) (*gnbctx.GnbAmf, error) {
	var err error
	amf, ok := c.Amfs[name]
	if !ok {
		err = fmt.Errorf("no corresponding Amf found for:%v", name)
	}
	return amf, err
}

func (c *Configuration) GetGNodeB(name string) (*gnbctx.GNodeB, error) {
	var err error
	gnb, ok := c.Gnbs[name]
	if !ok {
		err = fmt.Errorf("no corresponding gNodeB found for:%v", name)
	}
	return gnb, err
}

func (c *Configuration) GetGNodeBAt(pos int) (*gnbctx.GNodeB, error) {
	var err error
	if pos >= len(c.Gnbs) {
		err = fmt.Errorf("no corresponding gNodeB found at pos:%v", pos)
		return nil, err
	}
	keys := make([]string, 0, len(c.Gnbs))
	for k := range c.Gnbs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return c.Gnbs[keys[pos]], nil
}

func (c *Configuration) GetUeProfile(model string) (*UeProfile, error) {
	var err error
	ue, ok := c.UeProfiles[model]
	if !ok {
		err = fmt.Errorf("no corresponding Ue model found for:%v", model)
	}
	return ue, err
}

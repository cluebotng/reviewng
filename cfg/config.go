package cfg

// MIT License
//
// Copyright (c) 2021 Damian Zaremba
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var ReleaseTag = "development"

type Config struct {
	Runtime struct {
		Release string
	}
	Session struct {
		SecretKey string `yaml:"key"`
	}
	Db struct {
		Host string
		Port int
		User string
		Pass string
		Name string
	}
	OAuth struct {
		Token  string
		Secret string
	}
	App struct {
		UpdateStats bool `yaml:"update_stats"`
		AdminOnly   bool `yaml:"admin_only"`
	}
	Wikipedia struct {
		Username string `yaml:"username"`
		Password   string `yaml:"password"`
	}
}

func LoadConfigFromDisk(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := Config{}
	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		return nil, err
	}

	config.Runtime.Release = ReleaseTag
	return &config, nil
}

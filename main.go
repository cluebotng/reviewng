package main

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
	"embed"
	"fmt"
	"github.com/cluebotng/reviewng/cfg"
	"github.com/cluebotng/reviewng/controllers"
	"os"
)

//go:embed templates/*
var fsTemplates embed.FS

//go:embed static/*
var fsStatic embed.FS

func main() {
	configPath := "config.yaml"
	if val, ok := os.LookupEnv("REVIEW_CFG"); ok {
		configPath = val
	}

	config, err := cfg.LoadConfigFromDisk(configPath)
	if err != nil {
		panic(err)
	}

	app := controllers.NewApp(config, &fsTemplates, &fsStatic)

	listenAddr := "0.0.0.0:8080"
	if val, ok := os.LookupEnv("PORT"); ok {
		listenAddr = fmt.Sprintf("0.0.0.0:%s", val)
	}

	fmt.Printf("Listening on %+v\n", listenAddr)
	app.RunForever(listenAddr)
}

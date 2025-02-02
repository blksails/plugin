package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-plugin"
	bplugin "pkg.blksails.net/plugin" // 替换为实际的包路径
)

type ExamplePlugin struct {
	logger *log.Logger
}

func newLog(filename string) (*log.Logger, error) {
	// 设置日志文件
	logPath := filepath.Dir(filename)
	// 确保日志目录存在
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return log.New(logFile, "[ExamplePlugin] ", log.LstdFlags), nil
}

func (p *ExamplePlugin) Start(config map[string]string) (bplugin.PluginInfo, error) {
	p.logger.Println("Example plugin started")

	return bplugin.PluginInfo{
		Name:        "example",
		Description: "example plugin",
		Version:     "1.0.0",
		Author:      "blksails",
		Email:       "blksails@gmail.com",
		Url:         "https://blksails.net",
	}, nil
}

func (p *ExamplePlugin) Stop() error {
	if p.logger != nil {
		p.logger.Println("Example plugin stopped")
	}
	os.Exit(-1)
	return nil
}

func (p *ExamplePlugin) Install(ablink *bplugin.ABLink) ([]string, error) {
	if p.logger != nil {
		p.logger.Println("Installing example plugin")
	}
	return []string{"example.js"}, nil
}

func (p *ExamplePlugin) Enter(sessionId string, ablink *bplugin.ABLink) (bool, error) {
	if p.logger != nil {
		p.logger.Printf("Enter called with session ID: %s", sessionId)
	}
	return true, nil
}

func main() {
	var test *plugin.ServeTestConfig

	if os.Getenv("TEST") == "true" {
		test = &plugin.ServeTestConfig{
			CloseCh: make(chan struct{}),
		}
	}

	logger, err := newLog("logs/example.log")
	if err != nil {
		log.Fatal(err)
	}

	expPlugin := &ExamplePlugin{
		logger: logger,
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: bplugin.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"ablink": &bplugin.ABLinkPluginImpl{
				Impl: expPlugin,
			},
		},
		// GRPCServer: plugin.DefaultGRPCServer,
		Test: test,
	})

	logger.Println("Example plugin quit")
}

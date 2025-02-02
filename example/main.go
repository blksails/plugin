package main

import (
	"log"
	"os"

	"github.com/hashicorp/go-hclog"
	bplugin "pkg.blksails.net/plugin"
	"pkg.blksails.net/plugin/plugin" // 替换为实际的包路径
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:            "ablink",
		Level:           hclog.Info,
		Output:          os.Stdout,
		IncludeLocation: true,
	})

	manager := plugin.NewManager(logger)

	// 加载插件目录中的所有插件
	err := manager.LoadPlugins("./bin/plugins")
	if err != nil {
		log.Fatal(err)
	}
	defer manager.CloseAll()

	// 获取并使用特定插件
	p, err := manager.GetPlugin("example")
	if err != nil {
		log.Fatal(err)
	}

	// 使用插件
	info, err := p.Start(nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("plugin info: %+v", info)

	// call Install
	scripts, err := p.Install(&bplugin.ABLink{
		Id: "123",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("scripts: %+v", scripts)
}

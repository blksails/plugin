package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	bplugin "pkg.blksails.net/plugin"
)

type Manager struct {
	sync.RWMutex
	plugins map[string]*plugin.Client
	logger  hclog.Logger
}

func NewManager(logger hclog.Logger) *Manager {
	if logger == nil {
		logger = hclog.New(&hclog.LoggerOptions{
			Name:   "plugin-manager",
			Level:  hclog.Debug,
			Output: os.Stdout,
		})
	}

	return &Manager{
		plugins: make(map[string]*plugin.Client),
		logger:  logger,
	}
}

func (m *Manager) LoadPlugins(pluginPath string) error {
	files, err := os.ReadDir(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		pluginPath := filepath.Join(pluginPath, file.Name())
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: bplugin.HandshakeConfig,
			Plugins:         bplugin.PluginMap,
			Cmd:             exec.Command(pluginPath),
			Logger:          m.logger,
			AllowedProtocols: []plugin.Protocol{
				plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		})

		m.Lock()
		m.plugins[file.Name()] = client
		m.Unlock()

		// 连接到插件
		_, err := client.Client()
		if err != nil {
			m.logger.Error("failed to connect to plugin", "name", file.Name(), "error", err)
			continue
		}
	}

	return nil
}

func (m *Manager) GetPlugin(name string) (bplugin.ABLinkPlugin, error) {
	m.RLock()
	client, exists := m.plugins[name]
	m.RUnlock()

	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	rpcClient, err := client.Client()
	if err != nil {
		return nil, fmt.Errorf("failed to get plugin client: %v", err)
	}

	raw, err := rpcClient.Dispense("ablink")
	if err != nil {
		return nil, fmt.Errorf("failed to dispense plugin: %v", err)
	}

	plugin, ok := raw.(bplugin.ABLinkPlugin)
	if !ok {
		return nil, fmt.Errorf("unexpected plugin type")
	}

	return plugin, nil
}

func (m *Manager) CloseAll() {
	m.Lock()
	defer m.Unlock()

	for _, client := range m.plugins {
		client.Kill()
	}
}

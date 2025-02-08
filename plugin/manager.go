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
	plugins     map[string]*plugin.Client
	logger      hclog.Logger
	pluginCache map[string]bplugin.PluginInfo
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
		plugins:     make(map[string]*plugin.Client),
		logger:      logger,
		pluginCache: make(map[string]bplugin.PluginInfo),
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

	// Clear cache
	m.pluginCache = make(map[string]bplugin.PluginInfo)

	for _, client := range m.plugins {
		client.Kill()
	}
}

// Plugins

// Plugins returns a list of all loaded plugins
func (m *Manager) Plugins() []bplugin.PluginInfo {
	m.RLock()
	defer m.RUnlock()

	plugins := make([]bplugin.PluginInfo, 0, len(m.plugins))
	for name, client := range m.plugins {
		// Try to get from cache first
		if cached, ok := m.pluginCache[name]; ok {
			plugins = append(plugins, cached)
			continue
		}

		// If not in cache, get new instance
		rpcClient, err := client.Client()
		if err != nil {
			m.logger.Error("failed to get plugin client", "name", name, "error", err)
			continue
		}

		raw, err := rpcClient.Dispense("ablink")
		if err != nil {
			m.logger.Error("failed to dispense plugin", "name", name, "error", err)
			continue
		}

		if plugin, ok := raw.(bplugin.ABLinkPlugin); ok {

			info, err := plugin.Start(nil)
			if err != nil {
				m.logger.Error("failed to start plugin", "name", name, "error", err)
				continue
			}

			m.pluginCache[name] = info
			plugins = append(plugins, info)
		}
	}
	return plugins
}

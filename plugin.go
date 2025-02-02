package plugin

import (
	"encoding/gob"
	"errors"
	"net/rpc"
	"net/url"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "pkg.blksails.net/plugin/proto"
)

type ABLink struct {
	Id                  string
	Name                string
	BackSource          string
	ReviewUrl           string
	Domains             []string
	Alias               []string
	RoUrls              []*url.URL
	CoUrls              []*url.URL
	ReUrls              []*url.URL
	Tags                []string
	CompIds             []string
	WhiteCompIds        []string
	Review              bool
	WhiteComp           bool
	PrepareDomainSize   int32
	Mode                string
	RecentPickDomain    string
	JumpMode            string
	Scope               string
	ScopeName           string
	ProtectCode         string
	ProtectCodeJs       string
	InstallJs           string
	InstallChecked      bool
	Percentile          float64
	PercentileRate      float64
	DomainFactory       string
	DomainThreshold     int32
	BlockThreshold      int32
	DisableInjectjs     bool
	DisableReview       bool
	DomainFactoryConfig map[string]string
	EmailTo             []string
	Links               []string
	IpCities            []string
	IspBlocks           []string
	ReverseCity         bool
	Disable             bool
	CreatedAt           *timestamppb.Timestamp
	UpdatedAt           *timestamppb.Timestamp
	RefreshAt           *timestamppb.Timestamp
}

type PluginInfo struct {
	Name        string
	Description string
	Version     string
	Author      string
	Email       string
	Url         string
}

type Plugin interface {
	Start(config map[string]string) (PluginInfo, error)
	Stop() error
	// Config(ctx context.Context) map[string]interface{}
}

type ABLinkPlugin interface {
	Plugin

	Install(ablink *ABLink) ([]string, error)             // 安装 JS
	Enter(sessionId string, ablink *ABLink) (bool, error) // 进入链接
}

// HandshakeConfig 用于插件和宿主程序之间的握手配置
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ABLINK_PLUGIN",
	MagicCookieValue: "ablink",
}

// PluginMap 是可用插件的映射
var PluginMap = map[string]plugin.Plugin{
	"ablink_grpc": &ABLinkGRPCPlugin{},
	"ablink":      &ABLinkPluginImpl{},
}

// ABLinkPluginImpl 实现了 go-plugin 的接口
type ABLinkPluginImpl struct {
	Impl ABLinkPlugin
}

// Server 必须实现
func (p *ABLinkPluginImpl) Server(*plugin.MuxBroker) (interface{}, error) {
	return &DefaultABLinkPlugin{Impl: p.Impl}, nil
}

// Client 必须实现
func (p *ABLinkPluginImpl) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ABLinkPluginClient{client: c}, nil
}

// DefaultABLinkPlugin 提供默认实现
type DefaultABLinkPlugin struct {
	Impl ABLinkPlugin
}

func (p *DefaultABLinkPlugin) Start(config map[string]string, resp *PluginInfo) error {
	info, err := p.Impl.Start(config)
	if err != nil {
		return err
	}

	*resp = info
	return nil
}

func (p *DefaultABLinkPlugin) Stop(args interface{}, resp *string) error {
	return p.Impl.Stop()
}

func (p *DefaultABLinkPlugin) Install(linkId *ABLink, resp *[]string) error {
	scripts, err := p.Impl.Install(linkId)
	if err != nil {
		return err
	}
	*resp = scripts
	return nil
}

func (p *DefaultABLinkPlugin) Enter(args []interface{}, resp *bool) error {
	sessionId := args[0].(string)
	ablink, ok := args[1].(*ABLink)
	if !ok {
		return errors.New("invalid type ablink")
	}
	allowed, err := p.Impl.Enter(sessionId, ablink)
	if err != nil {
		return err
	}
	*resp = allowed
	return nil
}

// ABLinkPluginClient RPC 客户端实现
type ABLinkPluginClient struct {
	client *rpc.Client
}

func (p *ABLinkPluginClient) Start(config map[string]string) (PluginInfo, error) {
	var resp PluginInfo
	return resp, p.client.Call("Plugin.Start", config, &resp)
}

func (p *ABLinkPluginClient) Stop() error {
	var resp interface{}
	return p.client.Call("Plugin.Stop", new(interface{}), &resp)
}

func (p *ABLinkPluginClient) Install(ablink *ABLink) ([]string, error) {
	var resp []string
	err := p.client.Call("Plugin.Install", ablink, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *ABLinkPluginClient) Enter(sessionId string, ablink *ABLink) (bool, error) {
	var resp bool
	err := p.client.Call("Plugin.Enter", []interface{}{sessionId, ablink}, &resp)
	if err != nil {
		return false, err
	}

	return resp, nil
}

// ABLinkGRPCPlugin 实现了 ABLinkPlugin 接口
type ABLinkGRPCPlugin struct {
	plugin.Plugin
	Impl ABLinkPlugin
}

// func (p *ABLinkGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
// 	pb.RegisterPluginServiceServer(s, &GRPCServer{Impl: p.Impl})
// 	return nil
// }

// func (p *ABLinkGRPCPlugin) GRPCClient(broker *plugin.GRPCBroker, c *grpc.ClientConn) error {
// 	p.Impl = &GRPCClient{client: pb.NewPluginServiceClient(c)}
// 	return nil
// }

// 辅助函数：转换 ABLink 到 protobuf 格式
func convertToPBLink(link *ABLink) *pb.ABLink {
	if link == nil {
		return nil
	}

	roUrls := make([]string, len(link.RoUrls))
	for i, u := range link.RoUrls {
		roUrls[i] = u.String()
	}

	coUrls := make([]string, len(link.CoUrls))
	for i, u := range link.CoUrls {
		coUrls[i] = u.String()
	}

	reUrls := make([]string, len(link.ReUrls))
	for i, u := range link.ReUrls {
		reUrls[i] = u.String()
	}

	return &pb.ABLink{
		Id:                  link.Id,
		Name:                link.Name,
		BackSource:          link.BackSource,
		ReviewUrl:           link.ReviewUrl,
		Domains:             link.Domains,
		Alias:               link.Alias,
		RoUrls:              roUrls,
		CoUrls:              coUrls,
		ReUrls:              reUrls,
		Tags:                link.Tags,
		CompIds:             link.CompIds,
		WhiteCompIds:        link.WhiteCompIds,
		Review:              link.Review,
		WhiteComp:           link.WhiteComp,
		PrepareDomainSize:   link.PrepareDomainSize,
		Mode:                link.Mode,
		RecentPickDomain:    link.RecentPickDomain,
		JumpMode:            link.JumpMode,
		Scope:               link.Scope,
		ScopeName:           link.ScopeName,
		ProtectCode:         link.ProtectCode,
		ProtectCodeJs:       link.ProtectCodeJs,
		InstallJs:           link.InstallJs,
		InstallChecked:      link.InstallChecked,
		Percentile:          link.Percentile,
		PercentileRate:      link.PercentileRate,
		DomainFactory:       link.DomainFactory,
		DomainThreshold:     link.DomainThreshold,
		BlockThreshold:      link.BlockThreshold,
		DisableInjectjs:     link.DisableInjectjs,
		DisableReview:       link.DisableReview,
		DomainFactoryConfig: link.DomainFactoryConfig,
		EmailTo:             link.EmailTo,
		Links:               link.Links,
		IpCities:            link.IpCities,
		IspBlocks:           link.IspBlocks,
		ReverseCity:         link.ReverseCity,
		Disable:             link.Disable,
		CreatedAt:           link.CreatedAt,
		UpdatedAt:           link.UpdatedAt,
		RefreshAt:           link.RefreshAt,
	}
}

func convertToABLink(link *pb.ABLink) *ABLink {
	if link == nil {
		return nil
	}

	return &ABLink{
		Id:                  link.Id,
		Name:                link.Name,
		BackSource:          link.BackSource,
		ReviewUrl:           link.ReviewUrl,
		Domains:             link.Domains,
		Alias:               link.Alias,
		Tags:                link.Tags,
		CompIds:             link.CompIds,
		WhiteCompIds:        link.WhiteCompIds,
		Review:              link.Review,
		WhiteComp:           link.WhiteComp,
		PrepareDomainSize:   link.PrepareDomainSize,
		Mode:                link.Mode,
		RecentPickDomain:    link.RecentPickDomain,
		JumpMode:            link.JumpMode,
		Scope:               link.Scope,
		ScopeName:           link.ScopeName,
		ProtectCode:         link.ProtectCode,
		ProtectCodeJs:       link.ProtectCodeJs,
		InstallJs:           link.InstallJs,
		InstallChecked:      link.InstallChecked,
		Percentile:          link.Percentile,
		PercentileRate:      link.PercentileRate,
		DomainFactory:       link.DomainFactory,
		DomainThreshold:     link.DomainThreshold,
		BlockThreshold:      link.BlockThreshold,
		DisableInjectjs:     link.DisableInjectjs,
		DisableReview:       link.DisableReview,
		DomainFactoryConfig: link.DomainFactoryConfig,
		EmailTo:             link.EmailTo,
		Links:               link.Links,
		IpCities:            link.IpCities,
		IspBlocks:           link.IspBlocks,
		ReverseCity:         link.ReverseCity,
		Disable:             link.Disable,
		CreatedAt:           link.CreatedAt,
		UpdatedAt:           link.UpdatedAt,
		RefreshAt:           link.RefreshAt,
	}
}

func init() {
	gob.Register(&PluginInfo{})
	gob.Register(map[string]string{})
	gob.Register(&url.URL{})
	gob.Register(&timestamppb.Timestamp{})
	gob.Register(&ABLink{})
}

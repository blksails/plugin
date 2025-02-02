package plugin

// type GRPCClient struct {
// 	client pb.PluginServiceClient
// }

// func (c *GRPCClient) Start(config string) error {
// 	var configMap map[string]string
// 	err := json.Unmarshal([]byte(config), &configMap)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = c.client.Start(context.Background(), &pb.StartRequest{Config: configMap})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *GRPCClient) Stop() error {
// 	_, err := c.client.Stop(context.Background(), &pb.StopRequest{})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *GRPCClient) Install(linkId string) ([]string, error) {
// 	resp, err := c.client.Install(context.Background(), &pb.InstallRequest{LinkId: linkId})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp.Scripts, nil
// }

// func (c *GRPCClient) Enter(sessionId string, linkId string) (bool, error) {
// 	resp, err := c.client.Enter(context.Background(), &pb.EnterRequest{SessionId: sessionId, LinkId: linkId})
// 	if err != nil {
// 		return false, err
// 	}
// 	return resp.Allowed, nil
// }

// type GRPCServer struct {
// 	Impl ABLinkPlugin
// }

// func (s *GRPCServer) Start(ctx context.Context, req *pb.StartRequest) (*pb.StartResponse, error) {
// 	config, err := json.Marshal(req.Config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	info, err := s.Impl.Start(req.Config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &pb.StartResponse{Info: info}, nil
// }

// func (s *GRPCServer) Stop(ctx context.Context, req *pb.StopRequest) (*pb.StopResponse, error) {
// 	return &pb.StopResponse{}, s.Impl.Stop()
// }

// func (s *GRPCServer) Install(ctx context.Context, req *pb.InstallRequest) (*pb.InstallResponse, error) {
// 	scripts, err := s.Impl.Install(req.LinkId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &pb.InstallResponse{Scripts: scripts}, nil
// }

// func (s *GRPCServer) Enter(ctx context.Context, req *pb.EnterRequest) (*pb.EnterResponse, error) {
// 	allowed, err := s.Impl.Enter(req.SessionId, req.LinkId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &pb.EnterResponse{Allowed: allowed}, nil
// }

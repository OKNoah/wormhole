package session

import "github.com/superfly/wormhole/messages"

// Session hold information about connected client
type Session interface {
	ID() string
	BackendID() string
	NodeID() string
	Client() string
	Cluster() string
	Endpoint() string
	Key() string
	Release() *messages.Release
	RequireStream() error
	RequireAuthentication() error
	Close()
}

type baseSession struct {
	id           string `redis:"id,omitempty"`
	nodeID       string `redis:"node_id,omitempty"`
	backendID    string `redis:"backend_id,omitempty"`
	clientAddr   string `redis:"client_addr,omitempty"`
	EndpointAddr string `redis:"endpoint_addr,omitempty"`
	ClusterURL   string `redis:"cluster_url,omitempty"`

	release *messages.Release
	store   *RedisStore

	sessions map[string]Session
}

func (s *baseSession) ID() string {
	return s.id
}

func (s *baseSession) BackendID() string {
	return s.backendID
}

func (s *baseSession) NodeID() string {
	return s.nodeID
}

func (s *baseSession) Client() string {
	return s.clientAddr
}

func (s *baseSession) Cluster() string {
	return s.ClusterURL
}

func (s *baseSession) Endpoint() string {
	return s.EndpointAddr
}

func (s *baseSession) Key() string {
	return "session:" + s.id
}

func (s *baseSession) Release() *messages.Release {
	return s.release
}

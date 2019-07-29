package server

import (
	"github.com/firmeve/firmeve/container"
	_ "github.com/firmeve/firmeve/server/http"
	"sync"
)

type Server interface {
	Run()
}

type Manager struct {
	servers map[Type]Server
}

type ServiceProvider struct {
	Firmeve *container.Firmeve `inject:"firmeve"`
}

type Type string

const (
	Http Type = `http`
	Websocket Type = `websocket`
	Tcp	Type = `tcp`
)

var (
	manager *Manager
	mu      sync.Mutex
	once    sync.Once
)

func init()  {
	firmeve := container.GetFirmeve()
	firmeve.Register(`server`, firmeve.GetContainer().Resolve(new(ServiceProvider)).(*ServiceProvider))
}

func NewServer() *Manager {
	if manager != nil {
		return manager
	}

	once.Do(func() {
		manager = &Manager{
			servers: make(map[Type]Server, 0),
		}
	})

	return manager
}

func (m *Manager) Add(name Type, server Server) {
	mu.Lock()
	defer mu.Unlock()

	m.servers[name] = server
}

func (m *Manager) Get(name Type) Server {
	return m.servers[name]
}

func (m *Manager) Run() {
	for _, server := range m.servers {
		go server.Run()
	}

	for {
		select {

		}
	}
}

// ------------------------------- ServiceProvider ------------------------------

func (sp *ServiceProvider) Register() {
	sp.Firmeve.GetContainer().Bind(`server`, NewServer())
}

func (sp *ServiceProvider) Boot() {
	manager.Add(`http`, sp.Firmeve.GetContainer().Get(`http.server`).(Server))
}

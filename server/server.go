package server

import (
	"fmt"
	"github.com/firmeve/firmeve"
	"sync"
)

type Server interface {
	Run()
}

type Manager struct {
	servers map[Type]Server
}

type ServiceProvider struct {
	Firmeve *firmeve.Firmeve `inject:"firmeve"`
}

type Type uint

var (
	manager *Manager
	mu      sync.Mutex
	once    sync.Once
)

const (
	Http Type = iota
	Websocket
	Socket
)

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
		fmt.Println("GGG")
		go server.Run()
	}

	for {
		select {

		}
	}
}

// ------------------------------- ServiceProvider ------------------------------

func (sp *ServiceProvider) Register() {
	sp.Firmeve.Bind(`server`, NewServer())
}

func (sp *ServiceProvider) Boot() {
	manager := sp.Firmeve.Get(`server`).(*Manager)
	servers := map[Type]string{Http: `http.server`,}
	for key, server := range servers {
		if sp.Firmeve.Has(`http.server`) {
			manager.Add(key, sp.Firmeve.Get(server).(Server))
		}
	}
}

//<<<<<<< Updated upstream
//// Context
//=======
//import "context"
//
//type Server interface {
//	Run()
//}
//
//type Http struct {
//
//}
//
//func test2()  {
//	hash := make(map[string]func(context.Context))
//}
//
//func test()  {
//	route := gin.Default()
//	route.Any(`/*`, func(context *gin.Context) {
//		//context.Request.
//		var msg struct {
//			Name    string `json:"user"`
//			Message string
//			Number  int
//		}
//		msg.Name = "Lena"
//		msg.Message = "hey"
//		msg.Number = 123
//		context.JSON(200,msg)
//	})
//
//	route.Run()
//}
//>>>>>>> Stashed changes

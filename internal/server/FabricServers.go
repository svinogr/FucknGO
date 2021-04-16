package server

import (
	"FucknGO/config"
	"FucknGO/config/ui"
	"FucknGO/internal/handler"
	"FucknGO/internal/jwt"
	"FucknGO/internal/server/model"
	"FucknGO/log"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type fabricServers struct {
	servers []*server
}

var instance *fabricServers
var once sync.Once

const ApiSlave = "/api/slave"
const ApiMaster = "/api"

// FabricServer constructs singleton
func FabricServer() (*fabricServers, error) {
	var err error
	once.Do(func() {
		instance = &fabricServers{}
		//TODO сделать нормальное создание массива
		instance.servers = make([]*server, 0)
	})

	return instance, err
}

// GetNewMasterServer creates and returns new master servers
func (f *fabricServers) GetNewMasterServer(address string, port string) *server {
	idServer := time.Now().Unix()

	server := server{}
	server.setup(address, port, uint64(idServer), false)

	f.servers = append(f.servers, &server)

	return &server
}

// GetNewSlaveServer creates and returns new slave  servers
func (f *fabricServers) GetNewSlaveServer(address string, port string) (*server, error) {

	for _, el := range f.servers {
		if el.port == port {
			return nil, errors.New("port is uses yet")
		}
	}

	idServer := time.Now().Unix()

	server := server{}

	server.setup(address, port, uint64(idServer), true)

	f.servers = append(f.servers, &server)

	server.server.RegisterOnShutdown(func() {
		fmt.Print("dwddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
	})

	return &server, nil
}

// RemoveServer removes from slice! only
func (f *fabricServers) RemoveServer(server server) {
	for i, el := range f.servers {
		if el.port == server.port {
			f.servers = append(f.servers[:i], f.servers[i+1:]...)
			//	f.servers[i] = nil
		}
	}
}

func (f *fabricServers) DeleteSlaveServer(sM *model.ServerModel) {

	for _, el := range f.servers {
		if !el.isSlave {
			continue
		}

		if el.Id() == sM.Id {
			err := el.server.Shutdown(context.Background())
			sM.Port = el.Port()
			sM.Address = el.Address()
			sM.StaticResource = el.StaticResource()
			sM.Id = el.Id()
			sM.IsRun = false
			f.RemoveServer(*el)
			break
		}
	}
}

//setupStaticResource setup serverApi handler by type serverApi slave/master and auth
func setupHandlers(s *server) {
	fabric := NewFabric()
	if s.isSlave {
		for _, e := range fabric.Handlers {
			s.mux.HandleFunc(ApiSlave+e.GetHandler().Path, e.GetHandler().HandlerFunc).Methods(e.GetHandler().Method)
		}
	} else {
		for _, e := range fabric.Handlers {
			fh := http.HandlerFunc(e.GetHandler().HandlerFunc)

			switch e.GetHandler().TypeRequest {
			case handler.TypeWeb:
				if e.GetHandler().NeedAuthToken {
					s.mux.Handle(ApiMaster+e.GetHandler().Path, jwt.CheckTokensInCookie(jwt.AccessOrRefresh(fh))).Methods(e.GetHandler().Method)
				}
			case handler.TypeApi:
				if e.GetHandler().NeedAuthToken {
					s.mux.Handle(ApiMaster+e.GetHandler().Path, jwt.CheckTokensInCookie(jwt.AccessOrRefresh(fh))).Methods(e.GetHandler().Method)
					//s.mux.Handle(ApiMaster+e.GetHandler().Path, jwt.JwtVerifMiddleware.Handler(jwt.ParseJWT(fh))).Methods(e.GetHandler().Method)
				}
			}

			s.mux.HandleFunc(ApiMaster+e.GetHandler().Path, e.GetHandler().HandlerFunc).Methods(e.GetHandler().Method)
		}
	}
}

//setupStaticResource set static dir for serverApi
func setupStaticResource(server *server) {
	conf, err := config.GetConfig()

	if err != nil {
		log.NewLog().Fatal(err)
	}

	defaultStaticResource := conf.JsonStr.UiConfig.WWW.Static
	storageStaticResource := conf.JsonStr.UiConfig.WWW.Storage

	server.staticResource = defaultStaticResource

	_, err = os.Stat(defaultStaticResource)

	if err != nil {
		err := ui.CopyResource(storageStaticResource, defaultStaticResource)

		if err != nil {
			log.NewLog().Fatal(err)
		}
	}
	// staticResource = "./ui/web/static"
	fileServer := http.FileServer(http.Dir(defaultStaticResource + "/web/static"))
	// особенность при использовании gorilla
	server.mux.PathPrefix("/static/{rest}").Handler(
		http.StripPrefix("/static", fileServer))
}

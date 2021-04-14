package server

import (
	"FucknGO/config"
	"FucknGO/config/ui"
	"FucknGO/internal/handler"
	"FucknGO/internal/jwt"
	"FucknGO/log"
	"errors"
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

const API_SLAVE = "/api/slave"
const API_MASTER = "/api"

// FabricServer construct singleton
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

		/*	if el == nil {
			continue
		}*/

		if el.port == port {
			return nil, errors.New("port is uses yet")
		}
	}

	idServer := time.Now().Unix()

	server := server{}

	server.setup(address, port, uint64(idServer), true)

	setupStaticResource(&server)
	setupHandlers(&server)

	f.servers = append(f.servers, &server)

	return &server, nil
}

func (f *fabricServers) RemoveServer(server server) {
	for i, el := range f.servers {
		if el.port == server.port {
			f.servers = append(f.servers[:i], f.servers[i+1:]...)
			//	f.servers[i] = nil
		}
	}
}

//setupStaticResource setup serverApi handler by type serverApi slave/master and auth
func setupHandlers(s *server) {
	fabric := NewFabric()
	if s.isSlave {
		for _, e := range fabric.Handlers {
			s.mux.HandleFunc(API_SLAVE+e.GetHandler().Path, e.GetHandler().HandlerFunc).Methods(e.GetHandler().Method)
		}
	} else {
		for _, e := range fabric.Handlers {
			fh := http.HandlerFunc(e.GetHandler().HandlerFunc)

			switch e.GetHandler().TypeRequest {
			case handler.TypeWeb:
				if e.GetHandler().NeedAuthToken {
					s.mux.Handle(API_MASTER+e.GetHandler().Path, jwt.CheckTokensInCookie(jwt.AccessOrRefresh(fh))).Methods(e.GetHandler().Method)
				}
			case handler.TypeApi:
				if e.GetHandler().NeedAuthToken {
					s.mux.Handle(API_MASTER+e.GetHandler().Path, jwt.CheckTokensInCookie(jwt.AccessOrRefresh(fh))).Methods(e.GetHandler().Method)
					//s.mux.Handle(API_MASTER+e.GetHandler().Path, jwt.JwtVerifMiddleware.Handler(jwt.ParseJWT(fh))).Methods(e.GetHandler().Method)
				}
			}

			s.mux.HandleFunc(API_MASTER+e.GetHandler().Path, e.GetHandler().HandlerFunc).Methods(e.GetHandler().Method)
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
		/*	err = os.MkdirAll(defaultStaticResource, 0777)

			if err != nil {
				log.NewLog().Fatal(err)
			}
		*/
		err := ui.CopyResource(storageStaticResource, defaultStaticResource)

		if err != nil {
			log.NewLog().Fatal(err)
		}

	}
	// staticResource = "./ui/web/static"
	fileServer := http.FileServer(http.Dir(defaultStaticResource))

	//serverApi.mux.Handle("/static/js/jquery-3.6.0.min.js", http.StripPrefix("/static", fileServer))
	server.mux.PathPrefix("/static/{rest}").Handler(
		http.StripPrefix("/static", fileServer))
}

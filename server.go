package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"web_projekt/v6/app/controller"
	"web_projekt/v6/app/utils"
)

func main() {

	MuxRouter := mux.NewRouter()

	MuxRouter.HandleFunc(controller.INDEX_ADDRESS, controller.Index)
	MuxRouter.HandleFunc(controller.INDEX_ADDRESS_1, controller.Index)

	MuxRouter.HandleFunc(controller.ROOMS_ADDRESS, controller.Rooms)

	MuxRouter.HandleFunc(controller.ROOMS_ID_ADDRESS, controller.RoomsSession)

	MuxRouter.HandleFunc(controller.LOGIN_ADDRESS, controller.Login)
	MuxRouter.HandleFunc(controller.REGISTER_ADDRESS, controller.Register)
	MuxRouter.HandleFunc(controller.LOGOUT_ADDRESS, controller.Logout)

	//MuxRouter.HandlerFunc("/user", nil);

	MuxRouter.PathPrefix("/css/").Handler(http.StripPrefix("/css", http.FileServer(http.Dir(utils.GetLocalEnv()+utils.PathToCss))))
	MuxRouter.PathPrefix("/scripts/").Handler(http.StripPrefix("/scripts", http.FileServer(http.Dir(utils.GetLocalEnv()+utils.PathToScripts))))
	MuxRouter.PathPrefix("/images/").Handler(http.StripPrefix("/images", http.FileServer(http.Dir(utils.GetLocalEnv()+utils.PathToImages))))
	MuxRouter.PathPrefix("/models/").Handler(http.StripPrefix("/models", http.FileServer(http.Dir(utils.GetLocalEnv()+utils.PathToModels))))
	MuxRouter.PathPrefix("/svg/").Handler(http.StripPrefix("/svg", http.FileServer(http.Dir(utils.GetLocalEnv()+utils.PathToSvg))))

	http.Handle("/", MuxRouter)
	
	server := http.Server{
		Addr: ":8080",
	}
	fmt.Println("Server has been booted")
	fmt.Println("Listening on Port: ':8080'")
	//server.ListenAndServe()


	err := server.ListenAndServe()
	//("https-server.crt", "https-server.key")
	fmt.Println(err)

}

/**
  author: kevin
 */
package web

import (
	"net/http"
	"fmt"
	"github.com/hyperledger/Demo/web/controller"
)

func  WebStart(app *controller.Application)  {

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", app.HomeView)
	http.HandleFunc("/index.html", app.IndexView)
	http.HandleFunc("/setInfo.html", app.SetInfoView)
	http.HandleFunc("/setReq", app.SetInfo)
	http.HandleFunc("/queryReq", app.QueryInfo)
	http.HandleFunc("/createChannel", app.CreateChannel)
	http.HandleFunc("/changeChannel", app.ChangeChannel)

	fmt.Println("启动Web服务, 监听端口号: 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("启动Web服务错误%v", err)
	}

}
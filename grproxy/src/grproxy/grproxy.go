package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
	"os"
	"github.com/talbright/go-zookeeper/zk"
)

var logger *log.Logger
var zkLogger *log.Logger


func init() {
	logger = log.New(os.Stdout, "[SE2_GrlLogger] ", log.Ldate|log.Ltime)
	zkLogger = log.New(os.Stdout, "[SE2_ZKLogger] ", log.Ldate|log.Ltime)
}

func main() {
	var counter int = 1
	conn, _ := connectWithOptions()
	znodeBasicExample(conn)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		v := r.URL.String()
		remoteProxy2, _ := url.Parse("http://nginx:80")
		if strings.Contains(v, "library") {
			if counter == 1 {
				remoteProxy2, _ = url.Parse("http://gserve1:8080")
					counter = counter - 1
				}else{
					remoteProxy2, _ = url.Parse("http://gserve2:8080")
					counter = counter + 1
				}
			
		}
		proxy2 := httputil.NewSingleHostReverseProxy(remoteProxy2)
		proxy2.ServeHTTP(w, r)
	})
	http.ListenAndServe(":8080", nil)
}

func connectWithOptions() (*zk.Conn, <-chan zk.Event) {
	conn, events, err := zk.Connect([]string{"zookeeper"},
		time.Second*15,
		zk.WithLogger(zkLogger),
		zk.WithDialer(net.DialTimeout),
		zk.WithConnectTimeout(time.Second*15))
	if err != nil {
		panic(err)
	}
	for event := range events {
		if event.State == zk.StateHasSession {
			return conn, events
		}
	}
	return conn, events
}

func znodeBasicExample(conn *zk.Conn) {
	logger.Print("basic znode example")
	var path string
	var err error

	logger.Printf("creating node /tester")
	if path, err = conn.Create("/tester", []byte{}, zk.FlagPersistent, zk.WorldACL(zk.PermAll)); err != nil {
		panic(err)
	}
	logger.Print("La ruta es :"+path)

	logger.Printf("checking if node %s exists\n", path)
	if yes, _, err := conn.Exists(path); !yes || err != nil {
		panic(err)
	}


	logger.Printf("setting node %s data\n", path)
	if _, err = conn.Set(path, []byte("hello"), -1); err != nil {
		panic(err)
	}

	logger.Printf("getting node %s data\n", path)
	if data, _, err := conn.Get(path); err != nil || string(data) != "hello" {
		panic(err)
	} else {
		logger.Printf("node data: %v", string(data))
	}

	hijo2 , _,_ := conn.Children("/zookeeper")
	logger.Print(hijo2)

	path2 := "/zookeeper"

	logger.Printf("setting node %s data\n", path2)
	if _, err = conn.Set(path2, []byte("Hola Mundo!"), -1); err != nil {
		panic(err)
	}

	data, _, _ := conn.Get(path2);
	logger.Printf("node data2: %v", string(data))
	

	logger.Printf("deleting node %s\n", path)
	if err = conn.Delete(path, -1); err != nil {
		panic(err)
	}

}

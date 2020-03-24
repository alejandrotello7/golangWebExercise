package main

import (
	"log"
	"net"
	"time"
	"os"
	"github.com/talbright/go-zookeeper/zk"
	"net/http"
	"fmt"
	"bytes"
	 "encoding/json"
	 "io/ioutil"
	"html/template"
)

var logger *log.Logger
var zkLogger *log.Logger

type jsonStructure struct {
	Row []struct {
		Key  string `json:"key"`
		Cell []struct {
			Column string `json:"column"`
			Values string `json:"$"`
		} `json:"Cell"`
	} `json:"Row"`
}

type htmlStructure struct{
	KeyStruc string
	ColumnFirst string
	ColumnSecond string
	ValueFirst string
	ValueSecond string
	Instance string

}

func init() {
	logger = log.New(os.Stdout, "[SE2_GrlLogger] ", log.Ldate|log.Ltime)
	zkLogger = log.New(os.Stdout, "[ZK_Logger] ", log.Ldate|log.Ltime)
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
}//connectWithOptions

func znodeBasicExample(conn *zk.Conn) {

	logger.Print("basic znode example")
	//var path string
	//var err error
	path2, _ := conn.Create("/zookeeper/gserve1", []byte{},zk.FlagPersistent, zk.WorldACL(zk.PermAll))
	logger.Print(path2)
	

} //zNodeBasic Example

func webHandler(w http.ResponseWriter, r *http.Request) {
	logger.Print("test")
    switch r.Method {
    case "POST":
      	
   		 var t RowsType
    	decoder := json.NewDecoder(r.Body)
    	decoder.Decode(&t)


    	z := t.encode()
    	encodedJSON2, _ := json.Marshal(&z)

    	reqi, statuss := http.NewRequest("PUT", "http://zookeeper:8080/se2:library/fakerows",bytes.NewBuffer(encodedJSON2))
    	log.Println(statuss)
    	reqi.Header.Set("Content-Type", "application/json")
    	client := &http.Client{}
    	 resp, _ := client.Do(reqi)
    	 resp.Body.Close()
    	     fmt.Println("response Status:", resp.Status)
    		fmt.Println("response Headers:", resp.Header)
    		bodyss, _ := ioutil.ReadAll(resp.Body)
    		fmt.Println("response Body:", string(bodyss))

    	log.Println(encodedJSON2)

    case "GET":
    	w.Header().Set("Content-Type", "text/html; charset=utf-8")
    	sizeBatch := []byte("<Scanner batch=\"50\"/>")
    	reqi, errUrl := http.NewRequest("PUT", "http://zookeeper:8080/se2:library/scanner/", bytes.NewBuffer(sizeBatch))
    	reqi.Header.Set( "Accept" , "text/plain")
    	reqi.Header.Set("Content-Type", "text/xml")
    	client := &http.Client{}
    	resp, _ := client.Do(reqi)
    	resp.Body.Close()

    	fmt.Println("Error status:", errUrl)
    	fmt.Println("Response Status:", resp.Status)
    	fmt.Println("URL:", resp.Header.Get("Location"))
    	urlScanner := resp.Header.Get("Location")

    	reqBatch, errBatch := http.NewRequest("GET", urlScanner, nil)
    	reqBatch.Header.Set("Accept", "application/json")
    	clientBatch := &http.Client{}
    	respBatch, _ := clientBatch.Do(reqBatch)
		xmlData, _ := ioutil.ReadAll(respBatch.Body)    	
    	
    	defer respBatch.Body.Close()

    	fmt.Println("Error status:", errBatch)
    	fmt.Println("Response Status:", respBatch.Status)
    	fmt.Println("Text:", string(xmlData))

    	var encodedResponse EncRowsType
    	json.Unmarshal(xmlData, &encodedResponse)
    	unEncodedResponse, _ := encodedResponse.decode()

    	unEncodedJSONResponse, _ := json.Marshal(unEncodedResponse)
    	
    	//nuevo
    	var tester2 jsonStructure
    	json.Unmarshal(unEncodedJSONResponse, &tester2)
		fmt.Println("Text unEncodedJSONResponse: ", string(unEncodedJSONResponse)) 
		fmt.Println("# of Rows is: ", len(tester2.Row))
		fmt.Println("Key: ",tester2.Row[0].Key)

		sizeJSONRow := len(tester2.Row)

		for x :=0; x<sizeJSONRow; x++{
			fmt.Println(tester2.Row[x].Key)
			for z :=0; z<len(tester2.Row[x].Cell) ; z++ {
				fmt.Println(tester2.Row[x].Cell[z].Column)
				fmt.Println(tester2.Row[x].Cell[z].Values)
			}
		}

		tmpl := template.Must(template.ParseFiles("page.html"))
		for i := 0; i < sizeJSONRow; i++ {

		/*fmt.Fprintf(w, " <h1> Key: %s </h1>",tester2.Row[i].Key)
    	fmt.Fprintf(w, "<h3> Cell,Column 0: </h3> <body> %s </body>",tester2.Row[i].Cell[0].Column)
    	fmt.Fprintf(w,"<h3> Cell,Column 0: </h3> <body> %s </body>",tester2.Row[i].Cell[0].Values)
    	fmt.Fprintf(w, "<h3> Cell,Column 1: </h3> <body> %s </body>",tester2.Row[i].Cell[1].Column)
    	fmt.Fprintf(w, "<h3> Cell,Column 1: </h3> <body> %s </body>",tester2.Row[i].Cell[1].Values)*/
    	temporal := tester2.Row[i].Key
    	for z :=0; z<len(tester2.Row[i].Cell) - 1 ; z++ {
    		htmlTester := htmlStructure{KeyStruc: temporal , ColumnFirst: tester2.Row[i].Cell[z].Column, ValueFirst: tester2.Row[i].Cell[z].Values, ColumnSecond: tester2.Row[i].Cell[z+1].Column, ValueSecond: tester2.Row[i].Cell[z+1].Values, Instance: "proudly served by gserve1" }
    		tmpl.Execute(w, htmlTester)
			}
		}
	
	


    	
    default:
        fmt.Fprintf(w, "GET and POST methods only")
    }
}


func main() {
	conn, _ := connectWithOptions()
	znodeBasicExample(conn)
	http.HandleFunc("/", webHandler)
	http.ListenAndServe(":8080", nil)
}//Main

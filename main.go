package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
	"strconv"
)

type counters struct {
	sync.Mutex
	view  map[string]int
	click map[string]int
}

var (
	c = counters{}
	lastKey = ""
	content = []string{"sports", "entertainment", "business", "education"}
	title = []string{}
	record = []string{}
)

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to EQ Works 😎")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	data := content[rand.Intn(len(content))]

	key := data +  ":" + strconv.Itoa(time.Now().Year()) + "-" + strconv.Itoa(int(time.Now().Month())) + "-" + strconv.Itoa(time.Now().Day()) + "  " + strconv.Itoa(time.Now().Hour()) + ":" +strconv.Itoa(time.Now().Minute())
	lastKey = key
	c.Lock()
	c.view[key]++
	c.Unlock()
	title = append(title, key)

	err := processRequest(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// simulate random click call
	if rand.Intn(100) < 50 {
		processClick(key)
	}


	value := "{views:" + strconv.Itoa(c.view[key]) + "clicks:" + strconv.Itoa(c.click[key]) + "}"
	record = append(record, value)
	fmt.Println(key)
	fmt.Println(value)

}

func processRequest(r *http.Request) error {
	time.Sleep(time.Duration(rand.Int31n(50)) * time.Millisecond)
	return nil
}

func processClick(data string) error {
	c.Lock()
	c.click[data]++
	c.Unlock()

	return nil
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if !isAllowed() {
		w.WriteHeader(429)
		return
	}
}

func isAllowed() bool {
	return true
}

func uploadCounters(f *os.File) error {
	for isAllowed(){
		time.Sleep(5 * time.Second)
		//messenge := "{views:" + string(c.view) + "clicks:" + string(c.click) + "}"
		//fmt.Fprintln(f, lastKey)
		//fmt.Fprintln(f, "{views:" , c.view[lastKey] , "clicks:" , c.click[lastKey] , "}")
		for i := range title{
			fmt.Fprintln(f, title[i])
			fmt.Fprintln(f, record[i])
		}
		title = []string{}
		record = []string{}
	}
	return nil
}

func main() {
	fmt.Println("the upload function will only upload the data that has not been loaded, that means there will not be any duplicate data in the MockStore.txt file")
	c.view = make(map[string]int)
	c.click = make(map[string]int)
	f, err := os.Create("MockStore.txt")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/stats/", statsHandler)
	go uploadCounters(f)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
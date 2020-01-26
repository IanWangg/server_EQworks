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
	view  int
	click int
}

var (
	c = counters{}

	content = []string{"sports", "entertainment", "business", "education"}
)

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to EQ Works 😎")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	data := content[rand.Intn(len(content))]

	c.Lock()
	c.view++
	c.Unlock()

	err := processRequest(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// simulate random click call
	if rand.Intn(100) < 50 {
		processClick(data)
	}

	key := data +  ":" + time.Now().String()
	value := "views:" + strconv.Itoa(c.view) + "clicks:" + strconv.Itoa(c.click)
	fmt.Println(key)
	fmt.Println(value)

}

func processRequest(r *http.Request) error {
	time.Sleep(time.Duration(rand.Int31n(50)) * time.Millisecond)
	return nil
}

func processClick(data string) error {
	c.Lock()
	c.click++
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

		fmt.Fprintln(f, "{views:" , c.view , "clicks:" , c.click , "}")
	}
	return nil
}

func main() {
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
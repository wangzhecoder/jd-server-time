package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type jdTime struct {
	Current string `json:"currentTime"`
}

var lastTime string

func getTime() {
	req, _ := http.NewRequest(
		"GET",
		"https://api.m.jd.com/client.action?functionId=queryMaterialProducts&client=wh5",
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)

	if err == nil {
		body, errB := ioutil.ReadAll(response.Body)
		if errB != nil {
			panic(errB.Error())
		}
		jd := &jdTime{}
		json.Unmarshal(body, jd)
		if len(jd.Current) != 0 && jd.Current != lastTime {
			fmt.Println(jd.Current)
			lastTime = jd.Current
		}

		// return body
	}
	defer response.Body.Close()
	// panic(err.Error())
}

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	for {
		select {
		case s := <-done:
			fmt.Println(s)
			goto END
		default:
			getTime()
			time.Sleep(time.Microsecond)
		}
	}
END:
	fmt.Println("Bye")
}

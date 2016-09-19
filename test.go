package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- string) {
	os.Chdir("gittest")
	for j := range jobs {
		results <- runCommand("git show " + strconv.Itoa(j) + ":test.txt")
	}
}

func testWorkers() {
	start := time.Now()
	//In order to use our pool of workers we need to send them work and collect their results. We make 2 channels for this.
	jobs := make(chan int, 100)
	results := make(chan string, 100)
	//This starts up 50 workers, initially blocked because there are no jobs yet.
	for w := 1; w <= 50; w++ {
		go worker(w, jobs, results)
	}
	//Here we send 9 jobs and then close that channel to indicate that’s all the work we have.
	for j := 0; j < 100; j++ {
		jobs <- j
	}
	close(jobs)
	//Finally we collect all the results of the work.
	for a := 0; a < 100; a++ {
		<-results
	}
	elapsed := time.Since(start)
	log.Printf("testWorkers took %s", elapsed/101)
}

func createGithubRepo(username string, password string, reponame string) (bool, string) {
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	body := strings.NewReader(`{"name":"` + reponame + `"}`)
	req, err := http.NewRequest("POST", "https://api.github.com/user/repos", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(username, password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
	response := string(jsonDataFromHttp)
	if strings.Contains(response, "Bad credentials") || strings.Contains(response, "Validation Failed") {
		return false, response
	}
	return true, response
}

func main() {
	// createBranches()
	// readBranches()
	// testWorkers()
	fmt.Println(runCommand("git push origin master"))
}

func readBranches() {
	os.Chdir("gittest")
	start := time.Now()
	for i := 0; i < 100; i++ {
		runCommand("git show " + strconv.Itoa(i) + ":test.txt")
	}
	elapsed := time.Since(start)
	log.Printf("readBranches took %s", elapsed/101)
	fmt.Println("Done")

}

func createBranches() {
	os.RemoveAll("./gittest")
	os.Mkdir("gittest", 0644)
	os.Chdir("gittest")
	runCommand("git init")
	d1 := []byte("hello, world")
	err := ioutil.WriteFile("test.txt", d1, 0644)
	if err != nil {
		log.Fatal(err)
	}
	runCommand("git add test.txt")
	runCommand("git commit -am 'added test.txt'")

	start := time.Now()
	for i := 0; i < 100; i++ {
		runCommand("git checkout --orphan " + strconv.Itoa(i))
		d1 = []byte("hello, world branch #" + strconv.Itoa(i))
		err = ioutil.WriteFile("test.txt", d1, 0644)
		if err != nil {
			log.Fatal(err)
		}
		runCommand("git add test.txt")
		runCommand("git commit -am 'added test.txt'")
		// 3 commands = 115 ms / command
	}

	elapsed := time.Since(start)
	log.Printf("createBranches took %s", elapsed)
	fmt.Println("Done")
}

func runCommand(fullCommand string) string {
	var (
		cmdOut []byte
		err    error
	)
	splitCommand := strings.Split(fullCommand, " ")
	if len(splitCommand) > 1 {
		if cmdOut, err = exec.Command(splitCommand[0], splitCommand[1:]...).Output(); err != nil {
			log.Println(splitCommand)
			log.Fatal(err)
		}
	} else {
		if cmdOut, err = exec.Command(splitCommand[0]).Output(); err != nil {
			log.Println(splitCommand)
			log.Fatal(err)
		}
	}
	return string(cmdOut)
}

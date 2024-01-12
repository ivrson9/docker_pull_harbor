package main

import (
	"context"
	"fmt"
	"time"
	"log"
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"runtime"
	
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {

	// OS check
	log.Printf("Running in %s",runtime.GOOS)
	if runtime.GOOS == "windows" {
		setUpInWindows()
	} else if runtime.GOOS == "linux"{		
		setUpInLinux()
	}

	// load command line arguments
	name := flag.String("name","world","name to print")
	flag.Parse()
	log.Printf("Starting sleepservice for %s",*name)

	// setup signal catching
	sigs := make(chan os.Signal, 1)
	// catch all signals since not explicitly listing
	//signal.Notify(sigs)
	signal.Notify(sigs, syscall.SIGQUIT)

	// method invoked upon seeing signal
	go func() {
		s := <-sigs
		log.Printf("RECEIVED SIGNAL: %s",s)
		AppCleanup()
		os.Exit(1)
	}()

	// infinite print loop
	for {
		log.Printf("hello %s",*name)
		
		// wait random number of milliseconds
		Nsecs := rand.Intn(3000)
		log.Printf("About to sleep %dms before looping again",Nsecs)
		time.Sleep(time.Millisecond * time.Duration(Nsecs))
	}
}

func AppCleanup() {
	log.Println("CLEANUP APP BEFORE EXIT!!!")
}

func setUpInWindows() {
	//Docker Daemon에 연결부분
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
		fmt.Println("Nim?")
	}

    //Docker api 호출하여 사용
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
		fmt.Println("Nim??")
	}

    //결과값 출력
	for _, container := range containers {
		fmt.Println(container.ID)
	}
}

func setUpInLinux() {
	//Docker Daemon에 연결부분
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
		fmt.Println("Nim?")
	}

    //Docker api 호출하여 사용
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
		fmt.Println("Nim??")
	}

    //결과값 출력
	for _, container := range containers {
		fmt.Println(container.ID)
	}
}

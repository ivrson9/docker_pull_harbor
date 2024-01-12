package main

import (
	"context"
	"fmt"
	"time"
	"flag"
	"math/rand"
	"io"
	"os"
	"os/signal"
	"syscall"
	"runtime"
	"encoding/base64"
	"encoding/json"
	
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
)

func main() {

	// OS check
	fmt.Printf("Running in %s\n",runtime.GOOS)
	if runtime.GOOS == "windows" {
		setUpInWindows()
	} else if runtime.GOOS == "linux"{		
		setUpInLinux()
	}

	// load command line arguments
	name := flag.String("name","world","name to print")
	flag.Parse()
	fmt.Printf("Starting sleepservice for %s\n",*name)

	// setup signal catching
	sigs := make(chan os.Signal, 1)
	// catch all signals since not explicitly listing
	//signal.Notify(sigs)
	signal.Notify(sigs, syscall.SIGQUIT)

	// method invoked upon seeing signal
	go func() {
		s := <-sigs
		fmt.Printf("RECEIVED SIGNAL: %s\n",s)
		AppCleanup()
		os.Exit(1)
	}()

	// infinite print loop
	for {
		fmt.Printf("hello %s",*name)
		
		// wait random number of milliseconds
		Nsecs := rand.Intn(3000)
		fmt.Printf("About to sleep %dms before looping again\n",Nsecs)
		time.Sleep(time.Millisecond * time.Duration(Nsecs))
	}
}

func AppCleanup() {
	fmt.Println("CLEANUP APP BEFORE EXIT!!!")
}

func setUpInWindows() {
	//Docker Daemon에 연결부분
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

    //Docker api 호출하여 사용
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

    //결과값 출력
	for _, container := range containers {
		fmt.Println(container.ID)
	}

	ImagePullHarbor("sdffffff")
}

func setUpInLinux() {
	//Docker Daemon에 연결부분
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

    //Docker api 호출하여 사용
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

    //결과값 출력
	for _, container := range containers {
		fmt.Println(container.ID)
	}
}


func ImagePullHarbor(img string) error {
	ctx := context.Background()
  
	authConfig := registry.AuthConfig{
	  Username:      "harbor_registry_username",
	  Password:      "harbor_registry_password",
	  ServerAddress: "harbor_registry_url",
	}
  
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
	  panic(err)
	}
  
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	
	// docker client is use to pull the image.
	dockerCli, err := newDockerClient()
	if err != nil {
	  log.Error().Msgf("failed to create docker client: %w", err)
	  return err
	}
  
	out, err := dockerCli.ImagePull(ctx, img, types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
	  return fmt.Errorf("failed to pull image: %w", err)
	}
	defer out.Close()
  
	_, err = io.Copy(os.Stdout, out)
	if err != nil {
	  return fmt.Errorf("failed to read image logs: %w", err)
	}
	return err
}

func newDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client 1: %w", err)
	}
	return cli, nil
}
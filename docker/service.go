package docker

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

// Service ...
type Service struct {
	Client *docker.Client
	// DeployerPayload
}

// NewService ...
func NewService(endpoint string) (*Service, error) {
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return nil, err
	}

	return &Service{
		Client: client,
		// DeployerPayload: payload,
	}, nil
}

// DeployCompose ...
func (ds *Service) DeployCompose(payload DeployerPayload) error {
	// TODO: integrate db
	payload.ID = strconv.FormatInt(time.Now().UnixNano(), 10)

	// save compose
	path, err := payload.SaveComposeFile()
	if err != nil {
		return err
	}
	fmt.Printf("Data: %+v\n", payload)

	// TODO : DOCKER LOGIN TO PRIVATE DOCKER HUB
	// err = ds.PullImage("docker", "1.12.0")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err = ds.PullImage("docker/compose", "1.8.0")
	if err != nil {
		return err
	}

	//  Run docker log in... pass in credentials....
	// usr, err := user.Current()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = ds.RunContainer(
	// 	docker.Config{
	// 		Image: "docker/compose:1.8.0",
	// 		Cmd:   []string{"pull"},
	// 		// Volumes: map[string]struct{}{
	// 		// 	// "/tmp/projects:/tmp/projects": {},
	// 		// },
	// 		Env:        data.CreateEnvs(),
	// 		WorkingDir: path,
	// 	},
	// 	docker.HostConfig{
	// 		Binds: []string{
	// 			"/tmp/projects:/tmp/projects",
	// 			"/var/run/docker.sock:/var/run/docker.sock",
	// 			// usr.HomeDir + "/.docker:/root/.docker",
	// 		},
	// 	},
	// )

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// docker run -w '/tmp/projects/230945890/Test App' -v '/tmp/projects:/tmp/projects' -v '/var/run/docker.sock:/var/run/docker.sock' docker/compose:1.8.0 up -d
	return ds.RunContainer(
		docker.Config{
			Image:      "docker/compose:1.8.0",
			Cmd:        []string{"up", "-d"},
			Env:        payload.CreateEnvs(),
			WorkingDir: path,
		},
		docker.HostConfig{
			Binds: []string{
				"/tmp/projects:/tmp/projects",
				"/var/run/docker.sock:/var/run/docker.sock",
			},
		},
	)
}

// RunContainer ...
func (ds *Service) RunContainer(config docker.Config, hostConfig docker.HostConfig) error {
	opts := docker.CreateContainerOptions{
		Config:     &config,
		HostConfig: &hostConfig,
	}

	cont, err := ds.Client.CreateContainer(opts)
	if err != nil {
		return err
	}
	log.Printf("Container created: %s\n", cont.ID)

	err = ds.Client.StartContainer(cont.ID, &docker.HostConfig{})
	if err != nil {
		return err
	}
	log.Printf("Waiting for container: %s\n", cont.ID)
	code, err := ds.Client.WaitContainer(cont.ID)
	if err != nil {
		return err
	}

	log.Printf("Container finished with code: %d\n", code)
	if code == 0 {
		log.Printf("Removing container  ID: %s\n", cont.ID)
		return ds.Client.RemoveContainer(docker.RemoveContainerOptions{
			ID: cont.ID,
		})
	}
	return errors.New("container exited with error")
}

// PullImage ...
func (ds *Service) PullImage(repo, tag string) error {
	auth := docker.AuthConfiguration{}
	fmt.Println("Pulling docker image", repo, tag)

	return ds.Client.PullImage(
		docker.PullImageOptions{
			Repository: repo,
			Tag:        tag},
		auth)
}

// DeployerPayload ...
type DeployerPayload struct {
	ID          string            `json:"id"`
	Project     string            `json:"project"`
	Registry    Registry          `json:"registry"`
	ComposeFile string            `json:"compose_file"`
	Extra       map[string]string `json:"extra"`
}

// Registry ...
type Registry struct {
	URL      string `json:"url"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// ComposePath ...
func (payload *DeployerPayload) ComposePath() (string, error) {
	// if id or project empty return error
	path := fmt.Sprintf("/tmp/projects/%s/%s", payload.ID, payload.Project)
	return path, nil
}

// SaveComposeFile ...
func (payload *DeployerPayload) SaveComposeFile() (string, error) {
	path, err := payload.ComposePath()
	if err != nil {
		return path, err
	}
	fmt.Printf("Path: %s\n", path)
	os.MkdirAll(path, 0777)
	sDec, err := base64.StdEncoding.DecodeString(payload.ComposeFile)
	if err != nil {
		return path, err
	}

	filePath := fmt.Sprintf("%s/docker-compose.yml", path)

	f, err := os.Create(filePath)
	if err != nil {
		return path, err
	}
	defer f.Close()

	_, err = f.Write(sDec)
	if err != nil {
		return path, err
	}
	return path, f.Sync()
}

// CreateEnvs ...
func (payload *DeployerPayload) CreateEnvs() []string {
	envs := make([]string, len(payload.Extra))
	i := 0
	for k, v := range payload.Extra {
		envs[i] = fmt.Sprintf("%s=%s", k, v)
		i++
	}
	return envs
}

package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

const (
	DEFAULT_TMP               = "tmp"
	CONTAINER_RUNNING_TIMEOUT = 5
)

type Client struct {
	Client              *client.Client
	Image               string
	TmpDirPathContainer string
	TmpDirPathHost      string
}

func NewClient(image string) (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	pwd, _ := os.Getwd()
	tmpDirPathContainer := pwd + string(os.PathSeparator) + DEFAULT_TMP
	tmpDirPathHost := os.Getenv("TMP_DIR")
	_ = os.Mkdir(tmpDirPathContainer, 077)
	return &Client{
		Client:              cli,
		Image:               image,
		TmpDirPathContainer: tmpDirPathContainer,
		TmpDirPathHost:      tmpDirPathHost,
	}, nil
}

func (c *Client) NameContainer() string {
	timestamp := time.Now().Format("20060102150405")
	containerName := fmt.Sprintf("%s-%s", "cntnr", timestamp)
	return containerName
}

func (c *Client) CreateContainer(ctx context.Context, cmd string) (container.CreateResponse, error) {
	timeout := CONTAINER_RUNNING_TIMEOUT
	config := &container.Config{
		Image:       c.Image,
		Cmd:         []string{"sh", "-c", cmd},
		StopTimeout: &timeout,
	}
	hostConfig := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:%s", c.TmpDirPathHost, "/code"),
		},
		Resources: container.Resources{
			Memory: 1024 * 1024 * 512,
		},
	}
	resp, err := c.Client.ContainerCreate(ctx, config, hostConfig, &network.NetworkingConfig{}, nil, c.NameContainer())
	if err != nil {
		return container.CreateResponse{}, err
	}
	return resp, nil
}

func (c *Client) RunContainer(ctx context.Context, resp container.CreateResponse) error {
	return c.Client.ContainerStart(ctx, resp.ID, container.StartOptions{})
}

func (c *Client) NameFile() string {
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s-%s", "code", timestamp)
	return fileName
}

func (c *Client) CreateCodeFile(code []byte, ext string) (string, error) {
	fileName := c.NameFile()
	codeFilePath := c.PathFileInContainer(fileName, ext)
	file, err := os.Create(codeFilePath)
	if err != nil {
		return "", err
	}
	if _, err := file.Write(code); err != nil {
		return "", err
	}
	file.Close()
	return fileName, nil
}

func (c *Client) PathFileInContainer(fileName string, ext string) string {
	if ext == "" {
		return fmt.Sprintf("%s%s%s", c.TmpDirPathContainer, string(os.PathSeparator), fileName)
	}
	return fmt.Sprintf("%s%s%s.%s", c.TmpDirPathContainer, string(os.PathSeparator), fileName, ext)
}

func (c *Client) RemoveFiles(fileName string, ext string) {
	codePath := c.PathFileInContainer(fileName, ext)
	execPath := c.PathFileInContainer(fileName, "")
	os.Remove(codePath)
	os.Remove(execPath)
}

func (c *Client) StopContainer(ctx context.Context, cntnr container.CreateResponse) error {
	timeout := CONTAINER_RUNNING_TIMEOUT
	return c.Client.ContainerStop(ctx, cntnr.ID, container.StopOptions{
		Signal:  "SIGTERM",
		Timeout: &timeout,
	})
}

func (c *Client) removeNullBytes(data []byte) []byte {
	noNullBytes := make([]byte, 0, len(data))
	for _, b := range data {
		if b != 0 {
			noNullBytes = append(noNullBytes, b)
		}
	}
	return noNullBytes
}

func (c *Client) GetLogsContainer(ctx context.Context, cntnr container.CreateResponse) ([]byte, error) {
	out, err := c.Client.ContainerLogs(ctx, cntnr.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return nil, err
	}
	ignore := make([]byte, 8)
	out.Read(ignore)
	output, err := io.ReadAll(out)
	if err != nil {
		return nil, err
	}
	return c.removeNullBytes(output), nil
}

func (c *Client) Cmd(translator string, fileName string) (string, string, error) {
	var compileCmd, execCmd string
	switch translator {
	case "python3":
		compileCmd = ""
		execCmd = fmt.Sprintf("python /code/%s.%s", fileName, translator)
	case "c":
		compileCmd = fmt.Sprintf("gcc /code/%s.%s -o /code/%s", fileName, translator, fileName)
		execCmd = fmt.Sprintf("./code/%s", fileName)
	default:
		return "", "", errors.New("not supported this translator")
	}
	return compileCmd, execCmd, nil
}

func (c *Client) ExecCode(translator string, code []byte) ([]byte, error) {
	ctx := context.Background()
	fileExt := translator
	fileName, err := c.CreateCodeFile(code, fileExt)
	if err != nil {
		return nil, err
	}
	defer c.RemoveFiles(fileName, fileExt)

	compileCmd, execCmd, err := c.Cmd(translator, fileName)
	if err != nil {
		return []byte(err.Error()), nil
	}
	var cmd string
	if compileCmd == "" {
		cmd = execCmd
	} else {
		cmd = fmt.Sprintf("%s && %s", compileCmd, execCmd)
	}
	cntnr, err := c.CreateContainer(ctx, cmd)
	defer c.Client.ContainerRemove(ctx, cntnr.ID, container.RemoveOptions{})
	if err != nil {
		return nil, err
	}

	if err := c.RunContainer(ctx, cntnr); err != nil {
		return nil, err
	}

	if err := c.StopContainer(ctx, cntnr); err != nil {
		return nil, err
	}

	logs, err := c.GetLogsContainer(ctx, cntnr)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

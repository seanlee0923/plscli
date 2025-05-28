package plscli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func (c *PlsClient) Run() error {

	if err := register(c); err != nil {
		return fmt.Errorf("failed to register command: %w", err)
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	return unregister(c)
}

func register(c *PlsClient) error {
	req := RegisterRequest{
		DeployName: c.DeployName,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := c.Http.Post(c.Url+"/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("can not register client")
	}

	var r RegisterResponse
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	c.ClientId = r.ClientId

	return nil
}

func unregister(c *PlsClient) error {
	req := DeleteRequest{
		ClientId:   c.ClientId,
		DeployName: c.DeployName,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := c.Http.Post(c.Url+"/unregister", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to unregister client")
	}

	var r DeleteResponse
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	if c.ClientId != r.ClientId {
		return fmt.Errorf("invalid client id")
	}

	return nil
}

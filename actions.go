package plscli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (c *PlsClient) RunWithContext(ctx context.Context) error {
	retryInterval := 5 * time.Second

	for {
		err := register(c)
		if err == nil {
			break
		}

		fmt.Printf("register failed: %v. retrying in %s...\n", err, retryInterval)

		select {
		case <-time.After(retryInterval):
			continue
		case <-ctx.Done():
			return fmt.Errorf("context canceled before successful registration")
		}
	}

	<-ctx.Done()

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

func (c *PlsClient) Registered() bool {
	return c.ClientId != ""
}

func (c *PlsClient) IsLeader() (bool, error) {
	req := LeaderRequest{
		ClientId:   c.ClientId,
		DeployName: c.DeployName,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return false, err
	}

	resp, err := c.Http.Post(c.Url+"/leader", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to check leader")
	}

	var r LeaderResponse
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return false, err
	}

	if r.ClientId != c.ClientId {
		return false, nil
	}

	return true, nil
}

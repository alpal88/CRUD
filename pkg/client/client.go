package client

import (
	"Desktop/golangProjects/CRUD/pkg"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	errNotStatusOk = errors.New("status not 200")
)

type Client struct {
	baseUrl string
}

func New(baseUrl string) *Client {
	if baseUrl == "" {
		baseUrl = pkg.REGULARURL
	}
	return &Client{
		baseUrl: baseUrl,
	}
}

func (c *Client) CreateUser(name string, age int) error {
	user := pkg.HttpData{
		Name: name,
		Age:  age,
	}
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("%s%s", c.baseUrl, pkg.CREATEADDRROUTE), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error with the post %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return errNotStatusOk
	}
	fmt.Printf("user %s age %d sucesfully created \n", name, age)
	return nil
}

func (c *Client) ReadUser(name string) error {
	user := pkg.HttpData{
		Name: name,
	}
	info, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error marshalling user data: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s%s", c.baseUrl, pkg.USERADDROUTE, name), bytes.NewBuffer(info))
	if err != nil {
		return fmt.Errorf("error creating http request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error in http call: %w", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	var responseData pkg.HttpData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return fmt.Errorf("error unmarshalling response: %w", err)
	}

	fmt.Printf("The user %s's age is: %d\n", responseData.Name, responseData.Age)
	return nil
}

func (c *Client) UpdateUser(name string, age int) error {
	user := pkg.HttpData{
		Name: name,
		Age:  age,
	}
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s%s%s", c.baseUrl, pkg.USERADDROUTE, name), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errNotStatusOk
	}
	fmt.Printf("user %s's age sucesfully updated to be %d \n", name, age)
	return nil
}

func (c *Client) DeleteUser(name string) error {
	user := pkg.HttpData{
		Name: name,
	}
	body, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error marshalling user data: %w", err)
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s%s%s", c.baseUrl, pkg.USERADDROUTE, name), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating http request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error in http call: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return errNotStatusOk
	}
	fmt.Printf("user %s succesfully deleted \n", name)
	return nil
}

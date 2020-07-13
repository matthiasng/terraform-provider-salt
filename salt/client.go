package salt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type client struct {
	address  string
	eauth    string
	username string
	password string
}

func (c *client) Run(target, targetType, function string) (commandResult, error) {

	// curl -H 'Content-Type: application/json' -d '[{"eauth":"pam","username":"vagrant","password":"vagrant","client":"local","tgt":"*","fun":"pillar.items"}]' localhost:8081/run

	payload := map[string]string{
		"eauth":    c.eauth,
		"username": c.username,
		"password": c.password,
		"client":   "local", // #todo ???
		"tgt":      target,
		"tgt_type": targetType,
		"fun":      function,
	}

	jsonReq, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/run", c.address), bytes.NewBuffer(jsonReq))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, err
	}

	result, err := parseCommandResult(data)
	if err != nil {
		return nil, err
	}
	
	return result, nil
}

type commandResult map[string]interface{}

func parseCommandResult(data map[string]interface{}) (commandResult, error) {
	_, ok := data["return"]
	if !ok {
		return nil, errors.New("'return' not set")
	}

	results, ok := data["return"].([]interface{}) // #todo check
	if !ok {
		return nil, errors.New("'return' invalid type")
	}
	if len(results) == 0 {
		return nil, errors.New("no command results")
	}

	grains, ok := results[0].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid result type")
	}

	r := commandResult{}
	for name, items := range grains {
		r[name] = items
	}

	return r, nil
}

package agify

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetAge(name string) (int, error) {
	resp, err := http.Get("https://api.agify.io?name=" + name)
	if err != nil {
		return 0, fmt.Errorf("failed to get age: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get age: status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	var age struct {
		Age int `json:"age"`
	}
	err = json.NewDecoder(resp.Body).Decode(&age)
	if err != nil {
		return 0, fmt.Errorf("failed to decode age: %w", err)
	}
	return age.Age, nil
}

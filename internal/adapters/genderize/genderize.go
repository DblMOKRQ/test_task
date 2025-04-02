package genderize

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetGender(name string) (string, error) {
	resp, err := http.Get("https://api.genderize.io/?name=" + name)
	if err != nil {
		return "", fmt.Errorf("failed to get gender: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get gender: status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	var gender struct {
		Gender string `json:"gender"`
	}
	err = json.NewDecoder(resp.Body).Decode(&gender)
	if err != nil {
		return "", fmt.Errorf("failed to decode gender: %w", err)
	}
	return gender.Gender, nil
}

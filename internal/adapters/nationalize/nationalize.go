package nationalize

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Nationalize struct {
	Name    string `json:"name"`
	Country []struct {
		Country_id  string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

func GetNationality(name string) (string, error) {
	resp, err := http.Get("https://api.nationalize.io/?name=" + name)
	if err != nil {
		return "", fmt.Errorf("failed to get nationality: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get nationality: status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	var nationalize Nationalize
	err = json.NewDecoder(resp.Body).Decode(&nationalize)
	if err != nil {
		return "", fmt.Errorf("failed to decode nationalize response: %w", err)
	}
	if len(nationalize.Country) == 0 {
		return "", fmt.Errorf("no country found for name %s", name)
	}
	maxProbability := 0.0
	countryID := ""
	for _, country := range nationalize.Country {
		if country.Probability > maxProbability {
			maxProbability = country.Probability
			countryID = country.Country_id
		}

	}
	return countryID, nil
}

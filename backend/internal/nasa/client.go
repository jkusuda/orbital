package nasa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	key        string
	baseURL    string
	httpClient *http.Client
}

type Asteroid struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Hazardous   bool              `json:"is_potentially_hazardous_asteroid"`
	Diameter    EstimatedDiameter `json:"estimated_diameter"`
	OrbitalData OrbitalData       `json:"orbital_data"`
}

type EstimatedDiameter struct {
	Kilometers DiameterRange `json:"kilometers"`
}

type DiameterRange struct {
	Min float64 `json:"estimated_diameter_min"`
	Max float64 `json:"estimated_diameter_max"`
}

type OrbitalData struct {
	Epoch         float64 `json:"epoch_osculation,string"`
	Eccentricity  float64 `json:"eccentricity,string"`
	SemiMajorAxis float64 `json:"semi_major_axis,string"`
	Inclination   float64 `json:"inclination,string"`
	AscendingNode float64 `json:"ascending_node_longitude,string"`
	PerihelionArg float64 `json:"perihelion_argument,string"`
	MeanAnomaly   float64 `json:"mean_anomaly,string"`
}

func NewClient(key string) *Client {
	return &Client{
		key:     key,
		baseURL: "https://api.nasa.gov/neo/rest/v1",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetNEO(id string) (*Asteroid, error) {
	params := url.Values{}
	params.Set("api_key", c.key)
	fullURL := c.baseURL + "/neo/" + id + "?" + params.Encode()

	resp, err := c.httpClient.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nasa: unexpected status %d", resp.StatusCode)
	}

	var asteroid Asteroid
	if err := json.NewDecoder(resp.Body).Decode(&asteroid); err != nil {
		return nil, err
	}
	return &asteroid, nil
}

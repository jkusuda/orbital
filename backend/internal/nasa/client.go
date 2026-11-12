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

type feedResponse struct {
	NearEarthObjects map[string][]Asteroid `json:"near_earth_objects"`
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
	// https://api.nasa.gov/neo/rest/v1/neo/3542519?api_key=DEMO_KEY
	params := url.Values{}
	params.Set("api_key", c.key)
	fullURL := c.baseURL + "/neo/" + id + "?" + params.Encode()

	resp, err := c.httpClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("nasa: get neo %s: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nasa: get neo %s: unexpected status %d", id, resp.StatusCode)
	}

	var asteroid Asteroid
	if err := json.NewDecoder(resp.Body).Decode(&asteroid); err != nil {
		return nil, fmt.Errorf("nasa: decode neo %s: %w", id, err)
	}
	return &asteroid, nil
}

func (c *Client) GetFeed(start string, end string) ([]Asteroid, error) {
	// https://api.nasa.gov/neo/rest/v1/feed?start_date=START_DATE&end_date=END_DATE&api_key=API_KEY
	params := url.Values{}
	params.Set("start_date", start)
	params.Set("end_date", end)
	params.Set("api_key", c.key)
	fullURL := c.baseURL + "/feed?" + params.Encode()

	resp, err := c.httpClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("nasa: get feed %s to %s: %w", start, end, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nasa: get feed %s to %s: unexpected status %d", start, end, resp.StatusCode)
	}

	var feed feedResponse
	if err := json.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("nasa: decode feed %s to %s: %w", start, end, err)
	}

	var asteroids []Asteroid
	for _, day := range feed.NearEarthObjects {
		asteroids = append(asteroids, day...)
	}
	return asteroids, nil
}

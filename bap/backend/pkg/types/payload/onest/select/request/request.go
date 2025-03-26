package request

import "time"

type SeekerSelectPayload struct {
	WorkerID string              `json:"worker_id,omitempty"`
	ProviderID string            `json:"provider_id,omitempty"`
	JobID    string              `json:"job_id,omitempty"`
	BppID   string              `json:"bpp_id,omitempty"`
	BppURI  string              `json:"bpp_uri,omitempty"`
	Location SelectedJobLocation `json:"location,omitempty"`
}

type SelectedJobLocation struct {
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
	AreaCode    string      `json:"area_code"`
	Coordinates Coordinates `json:"coordinates"`
}

type Coordinates struct {
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitude float64 `bson:"longitude" json:"longitude"`
}

type SelectRequest struct {
	Context Context `json:"context"`
	Message Message `json:"message"`
}
type City struct {
	Code string `json:"code"`
}
type Country struct {
	Code string `json:"code"`
}
type Location struct {
	City    City    `json:"city"`
	Country Country `json:"country"`
}
type Context struct {
	Domain        string    `json:"domain"`
	Action        string    `json:"action"`
	Version       string    `json:"version"`
	BapID         string    `json:"bap_id"`
	BapURI        string    `json:"bap_uri"`
	BppID         string    `json:"bpp_id"`
	BppURI        string    `json:"bpp_uri"`
	TransactionID string    `json:"transaction_id"`
	MessageID     string    `json:"message_id"`
	Location      Location  `json:"location"`
	Timestamp     time.Time `json:"timestamp"`
	TTL           string    `json:"ttl"`
}
type Provider struct {
	ID string `json:"id"`
}
type Fulfillments struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
type Descriptor struct {
	Code string `json:"code"`
}
type List struct {
	Code  string `json:"code"`
	Value string `json:"value"`
}
type Tags struct {
	Descriptor Descriptor `json:"descriptor"`
	List       []List     `json:"list"`
}
type Items struct {
	ID   string `json:"id"`
	Tags []Tags `json:"tags"`
}
type Order struct {
	Provider     Provider       `json:"provider"`
	Fulfillments []Fulfillments `json:"fulfillments"`
	Items        []Items        `json:"items"`
}
type Message struct {
	Order Order `json:"order"`
}

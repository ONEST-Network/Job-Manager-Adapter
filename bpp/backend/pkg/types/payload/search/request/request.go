package request

type SearchRequest struct {
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
	Domain        string   `json:"domain"`
	Action        string   `json:"action"`
	Version       string   `json:"version"`
	BapID         string   `json:"bap_id"`
	BapURI        string   `json:"bap_uri"`
	TransactionID string   `json:"transaction_id"`
	MessageID     string   `json:"message_id"`
	Location      Location `json:"location"`
	Timestamp     string   `json:"timestamp"`
	TTL           string   `json:"ttl"`
}
type PaymentDescriptor struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
type PaymentList struct {
	Code  string `json:"code"`
	Value string `json:"value"`
}
type Payment struct {
	Descriptor PaymentDescriptor `json:"descriptor"`
	List       []PaymentList     `json:"list"`
}
type ItemDescriptor struct {
	Name string `json:"name"`
}
type Item struct {
	Descriptor ItemDescriptor `json:"descriptor"`
}
type TagsDescriptor struct {
	Code string `json:"code"`
}
type List struct {
	Descriptor TagsDescriptor `json:"descriptor"`
	Value      string         `json:"value"`
}
type Tags struct {
	Descriptor TagsDescriptor `json:"descriptor"`
	List       []List         `json:"list"`
}
type Intent struct {
	Payment Payment `json:"payment"`
	Item    Item    `json:"item"`
	Tags    []Tags  `json:"tags"`
}
type Message struct {
	Intent Intent `json:"intent"`
}

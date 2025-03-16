package searchresponse

import (
	searchresponse "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/types/payload/onest/search/response"
)

type SearchJobResponse struct {
	ID            string                          `bson:"id"`
	TransactionID string                          `bson:"transaction_id"`
	JobsResponse  []searchresponse.SearchResponse `bson:"jobs_response"`
}

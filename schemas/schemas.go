package schemas

import (
	"encoding/json"
)

type Source struct {
	URL string `json:"url" validate:"required,url"`
}

type AgentResponse struct {
	Answer  string   `json:"answer" validate:"required"`
	Sources []Source `json:"sources"`
}

// Example: Marshal AgentResponse to JSON
func (ar *AgentResponse) ToJSON() (string, error) {
	b, err := json.Marshal(ar)
	return string(b), err
}

func GetFormatInstructions() string {
	return `{
    "answer": string, // The final answer to the input question
    "sources": [      // List of sources used to generate the answer
        {
        "url": string // The URL of the source
        }
    ]
    }`
}

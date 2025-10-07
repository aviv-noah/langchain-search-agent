package main

import (
	"context"
	"helpers"
	"os"
	"tavily-go-api/pkg/tavily"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
)

func main() {
	_ = helpers.LoadDotEnv(".env")
	openai_api_key, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		panic("OPENAI_API_KEY environment variable not set")
	}
	var opts []openai.Option
	opts = append(opts, openai.WithModel("gpt-4"))
	opts = append(opts, openai.WithToken(openai_api_key))

	// Create your LLM and tools
	llm, err := openai.New()
	if err != nil {
		panic(err)
	}

	tsearch := tavily.New("http://localhost:8080")

	agentTools := []tools.Tool{
		tsearch,
	}

	// Create a ReAct-style agent
	agent := agents.NewOneShotAgent(
		llm,
		agentTools,
	)

	// Create the executor (AgentExecutor equivalent)
	executor := agents.NewExecutor(agent)

	// Run the agent
	result, err := executor.Call(
		context.Background(),
		map[string]any{
			"input": "Search for 3 job postings for an ai engineer using langchain in tlv area on linkedin with their urls and list their details",
		},
	)
	if err != nil {
		panic(err)
	}
	println(result["output"].(string))
}

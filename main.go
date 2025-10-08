package main

import (
	"context"
	"encoding/json"
	"helpers"
	"os"
	"search-agent/prompt"
	"search-agent/schemas"
	"tavily-go-api/pkg/tavily"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
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
	ReActPrompt := prompt.REACT_PROMPT_WITH_FORMAT_INSTRUCTIONS
	promptTemplate := prompts.NewPromptTemplate(
		ReActPrompt,
		[]string{"tool_names", "input", "agent_scratchpad"},
	)
	// Allow us to preformat our prompt, if we already have some variables
	// that we know will always be the same.
	promptTemplate.PartialVariables = map[string]any{
		"format_instructions": schemas.GetFormatInstructions(),
	}

	// Create a ReAct-style agent
	agent := agents.NewOneShotAgent(
		llm,
		agentTools,
		agents.WithPrompt(promptTemplate),
	)

	// Create the executor (AgentExecutor equivalent)
	executor := agents.NewExecutor(agent)

	// Run the agent
	result, err := executor.Call(
		context.Background(),
		map[string]any{
			"input":      "Search for 3 job postings for an ai engineer using langchain in tel aviv area on linkedin and list their details",
			"tool_names": "tavily",
			"tools":      "tavily",
		},
	)
	if err != nil {
		panic(err)
	}

	// Validate and parse output
	var agentResp schemas.AgentResponse
	err = json.Unmarshal([]byte(result["output"].(string)), &agentResp)
	if err != nil {
		panic(err)
	}
	println(agentResp.ToJSON())
	//println(agentResp.Answer)
}

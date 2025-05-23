package gonkaopenai_test

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	gonkaopenai "github.com/libermans/gonka-openai/go"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func TestExampleUsage(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Log("Note: .env file not found or could not be loaded. Proceeding with existing environment variables.")
	}

	if os.Getenv(gonkaopenai.EnvPrivateKey) == "" {
		t.Log("Missing required environment variable: GONKA_PRIVATE_KEY")
		t.Skip("Skipping test: Missing GONKA_PRIVATE_KEY") // Use t.Skip for conditional test skipping
		return
	}

	t.Log("\n------ Test Environment ------") // Use t.Log for test output
	baseURL := gonkaopenai.GonkaBaseURL(nil)
	t.Log("Using Gonka Base URL:", baseURL)

	// The APIKey is often a mock or test-specific key in test environments
	client, err := gonkaopenai.NewGonkaOpenAI(gonkaopenai.Options{
		GonkaPrivateKey: os.Getenv(gonkaopenai.EnvPrivateKey),
	})
	if err != nil {
		t.Fatalf("Error creating client: %v", err) // Use t.Fatalf to fail the test on critical errors
		return
	}

	t.Log("\nSending request...")
	resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: "Qwen/QwQ-32B", // Model as a string, consistent with gonkaopenai.go
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Hello! Tell me a short joke for a test."), // Using UserMessage helper
		},
	})
	if err != nil {
		t.Fatalf("Error during API call: %v", err)
		return
	}

	if len(resp.Choices) == 0 {
		t.Fatal("Expected at least one choice in the response, got none.")
		return
	}

	t.Log("\nResponse:")
	t.Log(resp.Choices[0].Message.Content)
}

// TestDirectOpenAIUsage tests the OpenAI API directly without the Gonka wrapper.
func TestDirectOpenAIUsage(t *testing.T) {
	// Attempt to load .env file
	err := godotenv.Load()
	if err != nil {
		t.Log("Note: .env file not found or could not be loaded. Proceeding with existing environment variables.")
	}

	// 1. Get Gonka Private Key
	gonkaPrivateKey := os.Getenv(gonkaopenai.EnvPrivateKey)
	if gonkaPrivateKey == "" {
		t.Log("Missing required environment variable: GONKA_PRIVATE_KEY for Gonka HTTP client setup.")
		t.Skip("Skipping test: Missing GONKA_PRIVATE_KEY")
		return
	}

	t.Log("\n------ Test Manually Configured Gonka Client (using openai.Client) ------")

	// 2. Determine Gonka Address (mirroring logic from NewGonkaOpenAI)
	gonkaAddress := os.Getenv(gonkaopenai.EnvAddress)
	if gonkaAddress == "" {
		addr, errAddr := gonkaopenai.GonkaAddress(gonkaPrivateKey)
		if errAddr == nil {
			gonkaAddress = addr
		} else {
			t.Logf("Could not derive GonkaAddress automatically (err: %v), fallback might be incomplete without GonkaChainID", errAddr)
		}
	}
	if gonkaAddress == "" {
		t.Log("Warning: GonkaAddress could not be determined. GonkaHTTPClient might fail or use defaults.")
	}

	// 3. Get Gonka Base URL
	baseURL := gonkaopenai.GonkaBaseURL(nil) // Assuming no specific endpoints for this test
	t.Log("Using Gonka Base URL for manual client:", baseURL)
	t.Log("Using Gonka Private Key (for HTTP client):", gonkaPrivateKey[:5]+"...") // Log a snippet for verification
	t.Log("Using Gonka Address (for HTTP client):", gonkaAddress)

	// 4. Create Gonka HTTP Client
	customHTTPClient := gonkaopenai.GonkaHTTPClient(gonkaopenai.HTTPClientOptions{
		PrivateKey: gonkaPrivateKey,
		Address:    gonkaAddress,
		Client:     nil, // No base client override for this test
	})
	t.Log("Custom Gonka HTTP Client configured.")

	// 5. Initialize OpenAI Client with Gonka settings
	// Using the "mock-api-key" as per your previous change.
	// If you want to use a real key for this Gonka setup, change it here or load from env.
	clientAPIKey := "mock-api-key"
	var clientOptions []option.RequestOption
	clientOptions = append(clientOptions, option.WithBaseURL(baseURL))
	clientOptions = append(clientOptions, option.WithHTTPClient(customHTTPClient))
	clientOptions = append(clientOptions, option.WithAPIKey(clientAPIKey))

	client := openai.NewClient(clientOptions...)
	t.Log("Manually configured Gonka-like client created with API Key:", clientAPIKey)

	t.Log("Sending request with manually configured Gonka-like client...")
	resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: "Qwen/QwQ-32B", // Using the same model as TestExampleUsage
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Hello Gonka-configured client! Tell me a very short story."),
		},
	})
	if err != nil {
		t.Fatalf("Error during API call with Gonka-configured client: %v", err)
		return
	}

	if len(resp.Choices) == 0 {
		t.Fatal("Expected at least one choice in the response, got none.")
		return
	}

	t.Log("\nResponse from Gonka-configured client:")
	t.Log(resp.Choices[0].Message.Content)
}

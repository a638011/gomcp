package prompts

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// DebugHelpPromptArgs defines arguments for the debug assistance prompt
type DebugHelpPromptArgs struct {
	Error    string `json:"error" jsonschema:"required,the error message or stack trace"`
	Code     string `json:"code" jsonschema:"related code snippet (optional)"`
	Language string `json:"language" jsonschema:"programming language"`
	Context  string `json:"context" jsonschema:"additional context about when error occurs"`
}

// GetDebugHelpPrompt generates a debugging assistance prompt
func GetDebugHelpPrompt(
	ctx context.Context,
	req *mcp.GetPromptRequest,
	args DebugHelpPromptArgs,
) (*mcp.GetPromptResult, error) {
	language := args.Language
	if language == "" {
		language = "not specified"
	}

	codeSection := ""
	if args.Code != "" {
		codeSection = fmt.Sprintf(`

**Related Code:**
`+"```"+`
%s
`+"```"+``, args.Code)
	}

	contextSection := ""
	if args.Context != "" {
		contextSection = fmt.Sprintf("\n\n**Context:** %s", args.Context)
	}

	promptText := fmt.Sprintf(`Help me debug this error:

**Language:** %s

**Error Message:**
`+"```"+`
%s
`+"```"+`
%s%s

Please provide:
1. **Root Cause Analysis:** Explain what's causing this error
2. **Solution Steps:** Clear steps to fix the issue
3. **Prevention:** How to avoid this error in the future
4. **Best Practices:** Related best practices for this scenario

Be specific and provide code examples where helpful.`, language, args.Error, codeSection, contextSection)

	return &mcp.GetPromptResult{
		Description: "Debugging assistance with root cause analysis and solutions",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: promptText},
			},
		},
	}, nil
}

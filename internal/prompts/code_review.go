package prompts

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// CodeReviewPromptArgs defines arguments for the code review prompt
type CodeReviewPromptArgs struct {
	Code     string `json:"code" jsonschema:"required,the source code to review"`
	Language string `json:"language" jsonschema:"programming language (default: autodetect)"`
	Focus    string `json:"focus" jsonschema:"specific focus area (e.g., security, performance)"`
}

// GetCodeReviewPrompt returns a structured code review prompt
func GetCodeReviewPrompt(
	ctx context.Context,
	req *mcp.GetPromptRequest,
	args CodeReviewPromptArgs,
) (*mcp.GetPromptResult, error) {
	// Default language to autodetect
	language := args.Language
	if language == "" {
		language = "autodetect"
	}

	// Build focus section
	focusSection := ""
	if args.Focus != "" {
		focusSection = fmt.Sprintf("\n\n**Special Focus:** %s", args.Focus)
	}

	// Construct the prompt
	promptText := fmt.Sprintf(`You are an expert code reviewer. Please review the following code:

**Language:** %s
**Code:**
`+"```"+`
%s
`+"```"+`

**Review Criteria:**
- Code quality and readability
- Potential bugs or issues
- Best practices adherence
- Performance considerations
- Security concerns
- Documentation quality%s

Please provide:
1. Overall assessment (Good/Fair/Needs Improvement)
2. Specific issues found (with line numbers if applicable)
3. Recommendations for improvement
4. Positive aspects worth noting

Be constructive and specific in your feedback.`, language, args.Code, focusSection)

	return &mcp.GetPromptResult{
		Description: "Code review prompt with structured analysis criteria",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: promptText},
			},
		},
	}, nil
}

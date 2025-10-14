package prompts

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GitCommitPromptArgs defines arguments for the git commit message prompt
type GitCommitPromptArgs struct {
	Diff         string `json:"diff" jsonschema:"required,the git diff to summarize"`
	Type         string `json:"type" jsonschema:"commit type (feat/fix/docs/etc)"`
	Scope        string `json:"scope" jsonschema:"commit scope (optional)"`
	Conventional bool   `json:"conventional" jsonschema:"use conventional commits format"`
}

// GetGitCommitPrompt generates a git commit message prompt
func GetGitCommitPrompt(
	ctx context.Context,
	req *mcp.GetPromptRequest,
	args GitCommitPromptArgs,
) (*mcp.GetPromptResult, error) {
	formatInstructions := ""
	if args.Conventional {
		formatInstructions = `
**Format:** Use Conventional Commits format:
- type(scope): subject
- Blank line
- Body explaining what and why
- Footer with breaking changes if any

**Types:** feat, fix, docs, style, refactor, test, chore`
	} else {
		formatInstructions = `
**Format:** 
- Clear, concise subject line (50 chars or less)
- Blank line
- Detailed description if needed
- Use imperative mood ("Add feature" not "Added feature")`
	}

	typeScope := ""
	if args.Type != "" {
		typeScope = fmt.Sprintf("\n**Suggested Type:** %s", args.Type)
		if args.Scope != "" {
			typeScope += fmt.Sprintf("\n**Suggested Scope:** %s", args.Scope)
		}
	}

	promptText := fmt.Sprintf(`Generate a git commit message for the following changes:

**Git Diff:**
`+"```"+`diff
%s
`+"```"+`
%s
%s

Please provide a commit message that:
1. Clearly describes what changed
2. Explains why the change was made
3. Follows best practices for commit messages
4. Is clear and concise`, args.Diff, typeScope, formatInstructions)

	return &mcp.GetPromptResult{
		Description: "Git commit message generator based on diff analysis",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: promptText},
			},
		},
	}, nil
}

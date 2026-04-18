package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/matsen/bipartite/internal/flow"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(commentsSinceCmd)
}

var commentsSinceCmd = &cobra.Command{
	Use:   "comments-since <owner/repo> <RFC3339-timestamp>",
	Short: "List GitHub issue comments updated since a timestamp (JSON)",
	Long: `Print issue comments on <owner/repo> updated at or after the
given RFC3339 timestamp as a JSON array.

Wraps internal/flow.FetchIssueComments so external tools (e.g. the nexus
monitor watcher) can consume bip's GitHub fetcher without re-implementing
the gh-CLI pagination dance.

Output is the raw GitHubComment shape (id, body, user, created_at,
updated_at, html_url, issue_url, pull_request_url).

Examples:
  bip comments-since settylab/dotto-nexus 2026-04-17T20:00:00Z
  bip comments-since matsen/bipartite "$(date -u -d '1 hour ago' +%FT%TZ)"`,
	Args: cobra.ExactArgs(2),
	RunE: runCommentsSince,
}

func runCommentsSince(cmd *cobra.Command, args []string) error {
	repo := args[0]
	since, err := time.Parse(time.RFC3339, args[1])
	if err != nil {
		return fmt.Errorf("invalid RFC3339 timestamp %q: %w", args[1], err)
	}

	comments, err := flow.FetchIssueComments(repo, since)
	if err != nil {
		exitWithError(ExitError, "fetching comments: %v", err)
	}

	if comments == nil {
		comments = []flow.GitHubComment{}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(comments); err != nil {
		return fmt.Errorf("encoding output: %w", err)
	}
	return nil
}

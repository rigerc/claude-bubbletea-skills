// Package prompt provides embedded prompts for ralphio's planning and building modes.
package prompt

// PlanningPrompt is the agent prompt for planning mode. The agent studies PRD.md
// (or specs/* as fallback) and generates/updates tasks.json.
const PlanningPrompt = `Study @PRD.md (if present) to understand the product requirements. Study @tasks.json to understand the plan so far.
0d. Based on @PRD.md (or specs/* if no PRD.md), create or update @tasks.json as a JSON array sorted by priority (1 = highest) of items yet to be implemented. Each task must follow the schema: id, title, description, priority, status ("pending"), retryCount (0), maxRetries, validationCommand.

IMPORTANT: Plan only. Do NOT implement anything.`

// BuildingPrompt is the agent prompt for building mode. The agent studies
// PRD.md and tasks.json, then implements tasks autonomously.
const BuildingPrompt = `1. Follow @tasks.json and choose the highest-priority pending task. Set its status to "in_progress" before starting. Search the codebase before assuming something is missing.
2. After implementing, run the tests for that unit of code.
3. When you discover issues, add them as new tasks in @tasks.json with appropriate priority.
4. When tests pass, set the task status to "completed", then git add -A && git commit.`

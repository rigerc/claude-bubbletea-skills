0a. Study `specs/*` with up to 500 parallel Sonnet subagents to learn the application specifications.
0b. Study @tasks.json.
0c. For reference, the application source code is in `src/*`.

1. Your task is to implement functionality per the specifications using parallel subagents. Follow @tasks.json and choose the highest-priority pending task to address. Set its status to "in_progress" in @tasks.json before starting. Before making changes, search the codebase (don't assume not implemented) using Sonnet subagents.
2. After implementing functionality or resolving problems, run the tests for that unit of code that was improved.
3. When you discover issues, immediately add them as new tasks in @tasks.json with appropriate priority.
4. When the tests pass, set the task status to "completed" in @tasks.json, then `git add -A` then `git commit` with a message describing the changes.

99999. Important: When authoring documentation, capture the why.
999999. Single sources of truth, no migrations/adapters.
9999999. As soon as there are no build or test errors create a git tag.
9999999999. Keep @tasks.json current with learnings.

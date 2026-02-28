# TASK

## Scope
Harden runtime behavior and test hygiene for this repository.

## Constraints
- No git commits or tags from subprocesses unless explicitly requested.
- Keep changes minimal, testable, and production-safe.
- Prefer deterministic shutdown/startup behavior.

## Required Output
- Small PR-sized patch.
- Repro steps.
- Validation commands and expected results.
- Known risks/limits.

## Priority Tasks
1. Keep timer/event emission stable under restart and shutdown.
2. Ensure no cross-plugin business logic is embedded here.

## Done Criteria
- System events stop immediately on shutdown and resume cleanly on startup.

## Validation Checklist
- [ ] Build succeeds for this repo.
- [ ] Local targeted tests (if present) pass.
- [ ] No new background orphan processes remain.
- [ ] Logs clearly show failure causes.

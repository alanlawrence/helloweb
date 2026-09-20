---
name: browser-check
description: This skill should be used when verifying a helloweb UI/JS change in index.html — checking that a form resizes/renders/behaves correctly in an actual browser, not just via go build/vet/test. Applies whenever the user asks to "check the UI", "test in a browser", "verify visually", or CLAUDE.md's UI-testing guidance would otherwise require a manual check. Provides a headless-Chromium-via-Playwright setup for this specific dev environment, since no browser tooling exists in the base sandbox.
version: 1.0.0
---

# Browser check (Playwright, headless Chromium)

This environment has no Node.js and no MCP browser tooling. The working
setup for this project is a Python virtualenv running Playwright's
headless Chromium, installed at `~/.claude-tools/playwright-venv`. It
persists across sessions on this machine — usually there is nothing to
install, just use it.

## 0. Check it's there before reinstalling

```bash
test -x ~/.claude-tools/playwright-venv/bin/python && echo present || echo missing
```

If `missing`, run `scripts/setup.sh` in this skill directory (idempotent —
safe to re-run). It needs ~300MB of disk and network access to
`cdn.playwright.dev`/`registry.npmjs.org`-equivalent Playwright CDN, and
uses `apt`/`sudo` for OS-level Chromium dependencies (fonts, `libnss3`,
`xvfb`, etc.), so it's worth flagging to the user before running blind if
it's ever missing.

## 1. Start the app

From the `helloweb/` directory (containing `go.mod`):

```bash
go run helloweb/webserver > /tmp/webserver.log 2>&1 &
sleep 1
curl -s localhost:8080/hello   # sanity check it's up
```

Remember to stop it afterwards — `go run` spawns a child binary, so
`pkill -f "go run"` alone won't kill it. Find the real PID via the
listening socket and kill that:

```bash
ss -ltnp | grep 8080          # -> pid=NNNN
kill NNNN
```

## 2. Drive it with Playwright

Write a short script and run it with the venv's Python directly — no
need to activate the venv:

```bash
~/.claude-tools/playwright-venv/bin/python /path/to/script.py
```

See `scripts/example_check.py` in this skill directory for a working
template (navigate, type into a field, read back computed style /
`innerHTML`, screenshot). Key patterns:

- **Assert on real state, not just visuals**: `page.locator(sel).evaluate("el => getComputedStyle(el).height")` catches exactly the kind of bug manual testing found in this project before (a `\n\n` suffix silently collapsed by the browser's default HTML whitespace handling, invisible unless you check the rendered DOM/screenshot rather than the raw response body).
- **`inner_html()` / `inner_text()`** on the answer `<span>` to confirm what the AJAX handler actually injected.
- **Screenshots**: `page.screenshot(path="...")`, then view with the Read tool (it renders images directly) — use `full_page=True` for anything below the fold.
- Keep viewport reasonably wide (e.g. `900x700`) — this project's forms aren't responsive, so a narrow viewport gives a misleading wrap.

## 3. Clean up

Kill the dev server (step 1) when done. Screenshots/scripts belong in the
scratchpad directory, not this skill directory.

## Why this exists

Before this was set up, UI changes to `index.html` (issue #57's growing
textarea, issue #60's `<pre>`-wrapping bug) could only be verified by
asking the user to manually test in their own browser. This skill exists
so Claude can do that verification itself — go build/vet/test only prove
the Go code compiles and its own logic is correct; they say nothing about
how the browser actually renders the resulting HTML/JS.

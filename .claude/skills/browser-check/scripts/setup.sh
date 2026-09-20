#!/bin/bash
# Idempotent install of the Playwright/Chromium browser-check tooling for
# this dev environment. Safe to re-run: pip/playwright install are no-ops
# if already satisfied.
#
# Usage: ./setup.sh
set -e

VENV=~/.claude-tools/playwright-venv

if [ -x "$VENV/bin/python" ]; then
    echo "venv already present at $VENV"
else
    echo "Creating venv at $VENV ..."
    mkdir -p ~/.claude-tools
    python3 -m venv "$VENV"
fi

echo "Installing/upgrading playwright in the venv ..."
"$VENV/bin/pip" install --upgrade pip playwright

echo "Installing Chromium + OS deps (uses apt/sudo, ~300MB download) ..."
"$VENV/bin/playwright" install --with-deps chromium

echo "Smoke test ..."
"$VENV/bin/python" - <<'EOF'
from playwright.sync_api import sync_playwright
with sync_playwright() as p:
    b = p.chromium.launch()
    pg = b.new_page()
    pg.set_content("<p id='x'>ok</p>")
    assert pg.inner_text("#x") == "ok"
    b.close()
print("Playwright/Chromium smoke test passed.")
EOF

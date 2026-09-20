#!/usr/bin/env python3
"""Template for checking a helloweb UI form with headless Chromium.

Run with the browser-check venv's Python, e.g.:
    ~/.claude-tools/playwright-venv/bin/python example_check.py

Assumes `go run helloweb/webserver` is already running on localhost:8080
(see SKILL.md step 1). Copy/adapt this rather than editing it in place.
"""
from playwright.sync_api import sync_playwright

with sync_playwright() as p:
    browser = p.chromium.launch()
    page = browser.new_page(viewport={"width": 900, "height": 700})
    page.goto("http://localhost:8080/")

    # Jump to the section under test.
    page.click("a[href='#idxHyphenate']")

    # Interact with the control.
    box = page.locator("#hyphenateStr")
    box.click()
    box.type("this is a fairly long string to force the box to grow past "
              "six hundred pixels wide and start wrapping onto more than "
              "one line")
    page.wait_for_timeout(200)

    # Assert on real rendered state, not just the raw response text --
    # this is what catches bugs like a browser silently collapsing
    # whitespace that would look fine in curl's raw output.
    width = box.evaluate("el => getComputedStyle(el).width")
    height = box.evaluate("el => getComputedStyle(el).height")
    print("box width:", width, "height:", height)

    page.screenshot(path="/tmp/browser_check_before_submit.png")

    page.keyboard.press("Enter")
    page.wait_for_timeout(300)

    answer_html = page.locator("#Hyphenate").inner_html()
    print("answer innerHTML:", repr(answer_html))

    page.screenshot(path="/tmp/browser_check_after_submit.png",
                     full_page=True)

    browser.close()

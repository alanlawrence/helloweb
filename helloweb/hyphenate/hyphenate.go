// Package hyphenate provides string hyphenation and the HTML fragment
// generation used by the /hyphen webserver endpoint.
package hyphenate

import (
    "fmt"
    "strings"
    "unicode"
    "unicode/utf8"
)

// maxLength is the issue #56 acceptance-criteria limit on submitted
// strings (criterion 2.2 / criterion 4).
const maxLength = 256

// tooLongMsg is the exact error message required by issue #56 criterion 4.
const tooLongMsg = "Strings submitted for hyphenation must be less than " +
    "or equal to 256 characters"

// isSpecial reports whether r is one of the characters treated as a
// separator: any Unicode white space character (issue #58 -- covers
// space, tab, newline, carriage return, and the rest of unicode.IsSpace,
// not just the literal ' ' a UI might otherwise restrict entry to),
// underscore, and the punctuation listed in issue #56's "Special
// characters" section. Note that this set includes every character HTML
// treats specially when injected via innerHTML (< > & " '), so hyphenated
// output is inherently free of those characters -- no separate escaping
// is needed before it is written into the response fragment.
func isSpecial(r rune) bool {
    if unicode.IsSpace(r) {
        return true
    }
    switch r {
    case '_', '£', '$', '%', '^', '&', '*', '(', ')', '+', '=',
        '\\', '/', '|', '!', '"', '{', '}', '[', ']', ';', ':', '#',
        ',', '<', '>', '?', '¬', '`', '\'':
        return true
    }
    return false
}

// Hyphenate returns s with every special character (see isSpecial)
// replaced by a hyphen, runs of hyphens (whether from replacement or
// already literally present in s) collapsed to one, and any leading or
// trailing hyphen removed.
// Design note: the substitution is done long hand to reduce the risk of
//              error associated with dealing with Go raw strings not
//              permitting internal backticks and the regex package using
//              RE2 rather than POSIX ERE like sed. In RE2 backslashes are
//              active and the resulting escaping of special characters is
//              hard to read.
//              This long hand version is quite efficient also at O(n) which
//              aligns with the lean execution philosophy of this webserver.
func Hyphenate(s string) string {
    replaced := strings.Map(func(r rune) rune {
        if isSpecial(r) {
            return '-'
        }
        return r
    }, s)

    var b strings.Builder
    prevHyphen := false
    for _, r := range replaced {
        if r == '-' {
            if prevHyphen {
                continue
            }
            prevHyphen = true
        } else {
            prevHyphen = false
        }
        b.WriteRune(r)
    }

    return strings.Trim(b.String(), "-")
}

// GenerateHtml formats the hyphenation result as the HTML fragment the
// /hyphen endpoint responds with.
func GenerateHtml(hyphenated string) string {
    return fmt.Sprintf("%v", hyphenated)
}

// CalculateHtml hyphenates s and returns the result as an HTML fragment.
// Strings longer than maxLength return the issue #56 error fragment
// instead (criterion 4).
func CalculateHtml(s string) string {
    if utf8.RuneCountInString(s) > maxLength {
        return tooLongMsg
    }
    return GenerateHtml(Hyphenate(s))
}

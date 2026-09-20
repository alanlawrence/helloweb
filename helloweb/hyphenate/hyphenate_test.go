package hyphenate

import (
    "strings"
    "testing"
)

// TestHyphenateOverviewExamples covers every example listed in the issue
// #56 "Overview" section.
func TestHyphenateOverviewExamples(t *testing.T) {
    cases := []struct{ in, want string }{
        {"My file", "My-file"},
        {"My - hypen - and - spaces - string", "My-hypen-and-spaces-string"},
        {"My --- hyphens --- string", "My-hyphens-string"},
        {"My trailing hyphens---", "My-trailing-hyphens"},
        {"---My leading hyphens", "My-leading-hyphens"},
        {"My$file", "My-file"},
        {"My_file_", "My-file"},
        {"My_-_file", "My-file"},
        {"My__file", "My-file"},
        {"File on 2026/05/31", "File-on-2026-05-31"},
    }
    for _, c := range cases {
        if got := Hyphenate(c.in); got != c.want {
            t.Errorf("Hyphenate(%q) = %q, want %q", c.in, got, c.want)
        }
    }
}

// TestHyphenateSedTranscriptExamples covers the "More examples" sed
// transcript from issue #56, which is the acceptance-criteria reference
// implementation (criterion 1).
func TestHyphenateSedTranscriptExamples(t *testing.T) {
    cases := []struct{ in, want string }{
        {"My file", "My-file"},
        {"---My file---", "My-file"},
        {"---My --- file---", "My-file"},
        {"---My -_-- file---", "My-file"},
        {"My%file", "My-file"},
        {"My", "My"},
        {"My£file", "My-file"},
        {"My$file", "My-file"},
        {"My$¬`!\"£$%^&*()_+={}[];:@#~\\|/?<>,.file", "My-@-~-.file"},
        {"My@file", "My@file"},
        {"My~special@file.txt", "My~special@file.txt"},
        {"-My~special@file.txt", "My~special@file.txt"},
    }
    for _, c := range cases {
        if got := Hyphenate(c.in); got != c.want {
            t.Errorf("Hyphenate(%q) = %q, want %q", c.in, got, c.want)
        }
    }
}

// TestHyphenateAllSpecialCharsAreCollapsed confirms every character named
// in the issue's "Special characters" section is treated as a separator.
func TestHyphenateAllSpecialCharsAreCollapsed(t *testing.T) {
    special := "|\\/!\"£$%^&*()_+=[]{};'#,<>?¬`"
    for _, r := range special {
        in := "a" + string(r) + "b"
        if got := Hyphenate(in); got != "a-b" {
            t.Errorf("Hyphenate(%q) = %q, want %q", in, got, "a-b")
        }
    }
}

// TestHyphenateEmptyAndAllSpecial covers degenerate inputs.
func TestHyphenateEmptyAndAllSpecial(t *testing.T) {
    cases := []struct{ in, want string }{
        {"", ""},
        {"---", ""},
        {"   ", ""},
        {"___", ""},
    }
    for _, c := range cases {
        if got := Hyphenate(c.in); got != c.want {
            t.Errorf("Hyphenate(%q) = %q, want %q", c.in, got, c.want)
        }
    }
}

// TestHyphenateWhitespaceCharacters covers issue #58: any white space
// character submitted to the endpoint (not just the literal space a UI
// might otherwise restrict entry to) must be treated as a separator.
func TestHyphenateWhitespaceCharacters(t *testing.T) {
    cases := []struct{ in, want string }{
        {"My\nfile", "My-file"},     // shift+enter / line feed
        {"My\r\nfile", "My-file"},   // carriage return + line feed
        {"My\tfile", "My-file"},     // tab
        {"My\vfile", "My-file"},     // vertical tab
        {"My\ffile", "My-file"},     // form feed
        {"My\u00A0file", "My-file"}, // non-breaking space
        {"\n\tMy file\r\n", "My-file"},
    }
    for _, c := range cases {
        if got := Hyphenate(c.in); got != c.want {
            t.Errorf("Hyphenate(%q) = %q, want %q", c.in, got, c.want)
        }
    }
}

func TestGenerateHtml(t *testing.T) {
    if got := GenerateHtml("My-file"); got != "My-file" {
        t.Errorf("GenerateHtml(%q) = %q, want %q", "My-file", got, "My-file")
    }
}

func TestCalculateHtmlHyphenates(t *testing.T) {
    if got := CalculateHtml("My file"); got != "My-file" {
        t.Errorf("CalculateHtml(%q) = %q, want %q", "My file", got, "My-file")
    }
}

func TestCalculateHtmlAcceptsMaxLength(t *testing.T) {
    in := strings.Repeat("a", maxLength)
    if got := CalculateHtml(in); got != in {
        t.Errorf("CalculateHtml(256 char string) = %q, want unchanged input", got)
    }
}

func TestCalculateHtmlRejectsOverMaxLength(t *testing.T) {
    in := strings.Repeat("a", maxLength+1)
    got := CalculateHtml(in)
    if got != tooLongMsg {
        t.Errorf("CalculateHtml(257 char string) = %q, want %q", got, tooLongMsg)
    }
}

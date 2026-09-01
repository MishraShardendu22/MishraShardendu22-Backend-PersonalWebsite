package web

import (
	"strings"
	"testing"
)

func TestStatusPageRendersShellAndGroups(t *testing.T) {
	out := StatusPage("production", "go1.25.0")

	for _, want := range []string{
		"<!doctype html>",
		`<html lang="en"`,
		"<title>Shardendu Mishra | Portfolio API</title>",
		`<meta name="viewport"`,
		`<link rel="canonical"`,
		"Content <em>API</em>",
		"production",
		"go1.25.0",
		"<footer>",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("status page is missing %q", want)
		}
	}

	for _, g := range APIGroups {
		if !strings.Contains(out, g.Title) {
			t.Errorf("status page is missing group %q", g.Title)
		}
		for _, e := range g.Endpoints {
			if !strings.Contains(out, e.Path) {
				t.Errorf("status page is missing endpoint %q", e.Path)
			}
		}
	}
}

func TestStatusPageHasExactlyOneH1(t *testing.T) {
	if n := strings.Count(StatusPage("development", "go1.25.0"), "<h1"); n != 1 {
		t.Errorf("expected exactly one h1, found %d", n)
	}
	if n := strings.Count(NotFoundPage("/nope"), "<h1"); n != 1 {
		t.Errorf("expected exactly one h1 on the 404 page, found %d", n)
	}
}

func TestNotFoundPageEscapesThePath(t *testing.T) {
	out := NotFoundPage(`/<script>alert("x")</script>`)

	if strings.Contains(out, "<script>alert") {
		t.Error("the requested path was interpolated without escaping")
	}
	if !strings.Contains(out, "&lt;script&gt;") {
		t.Error("expected the path to be HTML-escaped")
	}
}

func TestEveryTableCellHasAHeaderRow(t *testing.T) {
	out := StatusPage("production", "go1.25.0")

	if strings.Count(out, "<thead>") != len(APIGroups) {
		t.Errorf("expected one thead per group (%d), found %d",
			len(APIGroups), strings.Count(out, "<thead>"))
	}
	if !strings.Contains(out, `<th scope="col">`) {
		t.Error("table headers must declare scope")
	}
}

func TestNoEmojiInRenderedOutput(t *testing.T) {
	pages := map[string]string{
		"status":   StatusPage("production", "go1.25.0"),
		"notfound": NotFoundPage("/missing"),
	}
	for name, out := range pages {
		for _, r := range out {
			if r >= 0x1F300 && r <= 0x1FAFF {
				t.Errorf("%s page contains an emoji rune %U", name, r)
			}
		}
	}
}

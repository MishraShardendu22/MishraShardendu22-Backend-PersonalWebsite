// Package web renders the small set of human-facing HTML surfaces this API
// serves: a status landing page at the root and a not-found page for browser
// requests. Everything else on this service returns JSON.
//
// The markup follows the Meridian design system shared across the
// mishrashardendu22 ecosystem. Because an API host must not depend on an
// external stylesheet, the token subset used here is inlined verbatim from
// the canonical meridian.tokens.css. Keep the literal values in sync with
// https://github.com/MishraShardendu22/agent-skills
// (.agents/skills/meridian-design-system/assets/meridian.tokens.css).
package web

import (
	"fmt"
	"html"
	"strings"
)

// Endpoint describes one publicly documented route on this service.
type Endpoint struct {
	Method string
	Path   string
	Auth   bool
	Note   string
}

// Group is a titled collection of endpoints rendered as one table.
type Group struct {
	Title       string
	Description string
	Endpoints   []Endpoint
}

const (
	siteName    = "Portfolio API"
	siteOwner   = "Shardendu Mishra"
	portfolio   = "https://mishrashardendu22.is-a.dev"
	blog        = "https://blogs.mishrashardendu22.is-a.dev"
	systemsLab  = "https://github.mishrashardendu22.is-a.dev"
	sourceRepo  = "https://github.com/MishraShardendu22/MishraShardendu22-Backend-PersonalWebsite"
	description = "Read-only content API behind the Shardendu Mishra portfolio: projects, experiences, skills, certifications, volunteer work, timeline, search and public developer statistics."
)

// APIGroups is the documented surface of this service. It is derived from the
// routes registered in route/*.go; when a route changes, update this list in
// the same commit.
var APIGroups = []Group{
	{
		Title:       "Content",
		Description: "Portfolio entities. Read endpoints are public; writes require a bearer token.",
		Endpoints: []Endpoint{
			{Method: "GET", Path: "/api/projects"},
			{Method: "GET", Path: "/api/projects/kanban"},
			{Method: "GET", Path: "/api/projects/:id"},
			{Method: "POST", Path: "/api/projects", Auth: true},
			{Method: "POST", Path: "/api/projects/updateOrder", Auth: true},
			{Method: "PUT", Path: "/api/projects/:id", Auth: true},
			{Method: "DELETE", Path: "/api/projects/:id", Auth: true},
			{Method: "GET", Path: "/api/experiences"},
			{Method: "GET", Path: "/api/experiences/:id"},
			{Method: "POST", Path: "/api/experiences", Auth: true},
			{Method: "PUT", Path: "/api/experiences/:id", Auth: true},
			{Method: "DELETE", Path: "/api/experiences/:id", Auth: true},
			{Method: "GET", Path: "/api/volunteer/experiences"},
			{Method: "GET", Path: "/api/volunteer/experiences/:id"},
			{Method: "POST", Path: "/api/volunteer/experiences", Auth: true},
			{Method: "PUT", Path: "/api/volunteer/experiences/:id", Auth: true},
			{Method: "DELETE", Path: "/api/volunteer/experiences/:id", Auth: true},
			{Method: "GET", Path: "/api/certifications"},
			{Method: "GET", Path: "/api/certifications/:id"},
			{Method: "POST", Path: "/api/certifications", Auth: true},
			{Method: "PUT", Path: "/api/certifications/:id", Auth: true},
			{Method: "DELETE", Path: "/api/certifications/:id", Auth: true},
			{Method: "GET", Path: "/api/skills"},
			{Method: "POST", Path: "/api/skills", Auth: true},
			{Method: "GET", Path: "/api/timeline", Note: "Experience timeline, chronologically ordered"},
		},
	},
	{
		Title:       "Search",
		Description: "Full-text search across every content collection.",
		Endpoints: []Endpoint{
			{Method: "GET", Path: "/api/search/", Note: "Query parameter: q"},
			{Method: "GET", Path: "/api/search/suggestions"},
		},
	},
	{
		Title:       "Developer statistics",
		Description: "Cached upstream statistics. These endpoints use a separate, stricter rate limiter.",
		Endpoints: []Endpoint{
			{Method: "GET", Path: "/api/github"},
			{Method: "GET", Path: "/api/github/stars"},
			{Method: "GET", Path: "/api/github/commits"},
			{Method: "GET", Path: "/api/github/languages"},
			{Method: "GET", Path: "/api/github/top-repos"},
			{Method: "GET", Path: "/api/github/calendar"},
			{Method: "GET", Path: "/api/leetcode"},
		},
	},
	{
		Title:       "Administration",
		Description: "Token issuance and verification for the Console.",
		Endpoints: []Endpoint{
			{Method: "POST", Path: "/api/admin/auth", Note: "Exchange the admin password for a token"},
			{Method: "GET", Path: "/api/admin/auth", Auth: true, Note: "Verify a token"},
		},
	},
}

// styles is the inlined Meridian subset. Token values are copied verbatim from
// the canonical stylesheet; the component rules below are the minimum needed
// to render these two pages.
const styles = `
:root{color-scheme:dark;
--ink-900:#121211;--ink-800:#171716;--ink-700:#1d1d1b;--ink-600:#232321;
--ink-500:#2e2e2b;--ink-400:#3b3b37;--ink-200:#6e6b65;--ink-150:#8f8b83;
--ink-100:#b7b3aa;--ink-0:#f5f3ee;
--iris-500:#8b7cff;--iris-400:#a396ff;--iris-wash:rgba(139,124,255,.10);
--iris-edge:rgba(139,124,255,.32);
--positive-500:#4caf7d;--positive-wash:rgba(76,175,125,.10);
--caution-500:#d8a84a;--caution-wash:rgba(216,168,74,.10);
--critical-500:#e06060;--critical-wash:rgba(224,96,96,.10);
--font-display:"Instrument Serif","Iowan Old Style","Palatino",Georgia,serif;
--font-text:Inter,ui-sans-serif,system-ui,-apple-system,"Segoe UI",sans-serif;
--font-mono:"IBM Plex Mono",ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;
--radius-xs:4px;--radius-sm:6px;--radius-md:10px;--radius-lg:14px;--radius-full:999px;
--gutter:clamp(1.25rem,4vw,3rem)}
*,*::before,*::after{box-sizing:border-box}
*{margin:0;padding:0}
html{-webkit-text-size-adjust:100%}
body{background:var(--ink-900);color:var(--ink-0);font-family:var(--font-text);
font-size:.9375rem;line-height:1.6;-webkit-font-smoothing:antialiased;
min-height:100vh;display:flex;flex-direction:column;
background-image:radial-gradient(circle at 82% -12%,rgba(139,124,255,.09),transparent 30rem);
background-repeat:no-repeat;background-attachment:fixed}
a{color:inherit;text-decoration:none}
a:hover{color:var(--iris-400)}
:focus-visible{outline:2px solid var(--iris-500);outline-offset:2px;border-radius:var(--radius-xs)}
::selection{background:rgba(139,124,255,.18)}
main{flex:1;width:100%;max-width:calc(72rem + var(--gutter)*2);margin-inline:auto;
padding:2.5rem var(--gutter) 5rem}
.kicker{font-size:.6875rem;font-weight:700;text-transform:uppercase;
letter-spacing:.14em;color:var(--ink-150);margin-bottom:.75rem}
h1{font-family:var(--font-display);font-weight:400;
font-size:clamp(2.25rem,1.6rem + 3.1vw,3.5rem);line-height:1.02;
letter-spacing:-.028em;margin-bottom:1rem;text-wrap:balance}
h1 em{font-style:italic;color:var(--iris-500)}
h2{font-family:var(--font-display);font-weight:400;font-size:1.75rem;
line-height:1.14;letter-spacing:-.018em}
.lead{font-size:1.0625rem;color:var(--ink-100);max-width:58ch}
.wordmark{display:inline-flex;align-items:center;gap:.75rem}
.mark{display:grid;place-items:center;width:1.75rem;height:1.75rem;
border:1px solid var(--iris-edge);border-radius:var(--radius-sm);
background:var(--iris-wash);color:var(--iris-400);font-family:var(--font-display);
font-size:.8125rem;line-height:1}
.wordmark b{font-size:.8125rem;font-weight:600;letter-spacing:-.01em}
.wordmark span{font-size:.6875rem;font-weight:600;letter-spacing:.14em;
text-transform:uppercase;color:var(--ink-150)}
.topbar{display:flex;align-items:center;justify-content:space-between;gap:1rem;
flex-wrap:wrap;height:auto;min-height:56px;padding:.75rem var(--gutter);
background:var(--ink-800);border-bottom:1px solid var(--ink-500)}
.topnav{display:flex;gap:.25rem;flex-wrap:wrap}
.topnav a{padding:.5rem .75rem;border-radius:var(--radius-sm);font-size:.8125rem;
font-weight:500;color:var(--ink-150)}
.topnav a:hover{color:var(--ink-0);background:rgba(245,243,238,.04)}
.status{display:inline-flex;align-items:center;gap:.5rem;padding:.25rem .75rem;
border:1px solid rgba(76,175,125,.26);border-radius:var(--radius-full);
background:var(--positive-wash);color:var(--positive-500);font-size:.75rem;
font-weight:600}
.dot{width:7px;height:7px;border-radius:var(--radius-full);
background:var(--positive-500);box-shadow:0 0 0 3px var(--positive-wash)}
.grid{display:grid;gap:1rem;grid-template-columns:repeat(auto-fit,minmax(11rem,1fr));
margin-block:2.5rem}
.metric{display:flex;flex-direction:column;gap:.5rem;padding:1.25rem;
background:var(--ink-700);border:1px solid var(--ink-500);border-radius:var(--radius-lg)}
.metric dt{font-size:.6875rem;font-weight:700;text-transform:uppercase;
letter-spacing:.14em;color:var(--ink-150)}
.metric dd{font-family:var(--font-display);font-size:1.75rem;line-height:1;
letter-spacing:-.028em;font-variant-numeric:tabular-nums}
.metric dd.text{font-family:var(--font-mono);font-size:.9375rem;line-height:1.4;
letter-spacing:0;color:var(--ink-100);word-break:break-word}
section{margin-top:3rem}
.sechead{display:flex;align-items:baseline;justify-content:space-between;gap:1rem;
flex-wrap:wrap;padding-bottom:1rem;border-bottom:1px solid var(--ink-500);
margin-bottom:1.5rem}
.sechead p{font-size:.8125rem;color:var(--ink-150);max-width:52ch}
.panel{background:var(--ink-800);border:1px solid var(--ink-500);
border-radius:var(--radius-lg);overflow-x:auto;-webkit-overflow-scrolling:touch}
table{width:100%;min-width:34rem;border-collapse:collapse;font-size:.8125rem}
th{text-align:left;padding:.75rem 1rem;background:var(--ink-700);
border-bottom:1px solid var(--ink-500);font-size:.6875rem;font-weight:700;
text-transform:uppercase;letter-spacing:.09em;color:var(--ink-150);white-space:nowrap}
td{padding:.75rem 1rem;border-bottom:1px solid rgba(245,243,238,.06);
color:var(--ink-100);vertical-align:top}
tbody tr:last-child td{border-bottom:none}
tbody tr:hover{background:rgba(245,243,238,.04)}
code,.mono{font-family:var(--font-mono);font-size:.9em}
.verb{display:inline-flex;align-items:center;padding:.1rem .4rem;
border:1px solid var(--ink-500);border-radius:var(--radius-xs);
background:var(--ink-600);font-family:var(--font-mono);font-size:.75rem;
color:var(--ink-100)}
.verb-get{border-color:rgba(139,124,255,.32);background:var(--iris-wash);color:var(--iris-400)}
.verb-post{border-color:rgba(76,175,125,.26);background:var(--positive-wash);color:var(--positive-500)}
.verb-put{border-color:rgba(216,168,74,.26);background:var(--caution-wash);color:var(--caution-500)}
.verb-delete{border-color:rgba(224,96,96,.26);background:var(--critical-wash);color:var(--critical-500)}
.lock{display:inline-flex;align-items:center;padding:.1rem .4rem;
border:1px solid var(--ink-500);border-radius:var(--radius-xs);
font-family:var(--font-mono);font-size:.75rem;color:var(--ink-150)}
.state{display:flex;flex-direction:column;align-items:center;text-align:center;
gap:1rem;max-width:34rem;margin:min(18vh,8rem) auto}
.state .icon{display:grid;place-items:center;width:2.75rem;height:2.75rem;
border:1px solid var(--ink-500);border-radius:var(--radius-lg);
background:var(--ink-600);color:var(--iris-400);font-family:var(--font-display);
font-size:1.25rem}
.state p{color:var(--ink-100);max-width:46ch}
.btn{display:inline-flex;align-items:center;justify-content:center;gap:.5rem;
height:2.5rem;padding:0 1.25rem;border:1px solid var(--ink-500);
border-radius:var(--radius-md);background:var(--ink-700);color:var(--ink-0);
font-size:.8125rem;font-weight:600;transition:background-color .14s,border-color .14s}
.btn:hover{background:var(--ink-600);border-color:var(--iris-edge);color:var(--ink-0)}
.btn-primary{background:var(--ink-0);border-color:var(--ink-0);color:var(--ink-900)}
.btn-primary:hover{background:var(--iris-500);border-color:var(--iris-500);color:var(--ink-900)}
footer{border-top:1px solid var(--ink-500);background:var(--ink-800);margin-top:auto}
.footinner{width:100%;max-width:calc(72rem + var(--gutter)*2);margin-inline:auto;
padding:2.5rem var(--gutter);display:grid;gap:2rem;
grid-template-columns:minmax(0,1.4fr) repeat(auto-fit,minmax(9rem,1fr))}
.footinner h2{font-family:var(--font-text);font-size:.6875rem;font-weight:700;
text-transform:uppercase;letter-spacing:.14em;color:var(--ink-150);margin-bottom:.25rem}
.footcol{display:flex;flex-direction:column;gap:.75rem;min-width:0}
.footcol a{font-size:.8125rem;color:var(--ink-100)}
.footnote{font-size:.8125rem;color:var(--ink-150);max-width:40ch}
.footbar{border-top:1px solid rgba(245,243,238,.06)}
.footbar div{width:100%;max-width:calc(72rem + var(--gutter)*2);margin-inline:auto;
padding:1.25rem var(--gutter);display:flex;justify-content:space-between;
gap:1rem;flex-wrap:wrap;font-size:.75rem;color:var(--ink-150)}
@media (max-width:640px){
.footinner{grid-template-columns:1fr}
.panel{background:transparent;border:none;overflow:visible}
table{min-width:0;border-collapse:separate;border-spacing:0}
thead{position:absolute;width:1px;height:1px;overflow:hidden;clip-path:inset(50%)}
tbody{display:grid;gap:1rem}
tbody tr{display:block;background:var(--ink-700);border:1px solid var(--ink-500);
border-radius:var(--radius-lg);overflow:hidden}
tbody tr:hover{background:var(--ink-700)}
td{display:flex;align-items:baseline;justify-content:space-between;gap:1rem;
min-height:44px;padding:.625rem .875rem;text-align:right;white-space:normal;
overflow-wrap:anywhere}
tbody tr td:last-child{border-bottom:none}
td::before{content:attr(data-label);flex:0 0 auto;text-align:left;
font-size:.6875rem;font-weight:700;letter-spacing:.09em;text-transform:uppercase;
color:var(--ink-150)}
}
@media (prefers-reduced-motion:reduce){*{transition-duration:.01ms!important}}
`

func verbClass(method string) string {
	switch method {
	case "GET":
		return "verb verb-get"
	case "POST":
		return "verb verb-post"
	case "PUT", "PATCH":
		return "verb verb-put"
	case "DELETE":
		return "verb verb-delete"
	default:
		return "verb"
	}
}

func shell(title, description, canonical, body string) string {
	var b strings.Builder
	b.WriteString(`<!doctype html><html lang="en" dir="ltr"><head><meta charset="utf-8">`)
	b.WriteString(`<meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover">`)
	fmt.Fprintf(&b, `<title>%s</title>`, html.EscapeString(title))
	fmt.Fprintf(&b, `<meta name="description" content="%s">`, html.EscapeString(description))
	b.WriteString(`<meta name="color-scheme" content="dark">`)
	b.WriteString(`<meta name="theme-color" content="#121211">`)
	b.WriteString(`<meta name="robots" content="noindex,follow">`)
	if canonical != "" {
		fmt.Fprintf(&b, `<link rel="canonical" href="%s">`, html.EscapeString(canonical))
	}
	fmt.Fprintf(&b, `<meta property="og:type" content="website"><meta property="og:title" content="%s">`, html.EscapeString(title))
	fmt.Fprintf(&b, `<meta property="og:description" content="%s">`, html.EscapeString(description))
	fmt.Fprintf(&b, `<meta property="og:site_name" content="%s">`, html.EscapeString(siteName))
	b.WriteString(`<meta name="twitter:card" content="summary">`)
	b.WriteString(`<link rel="preconnect" href="https://fonts.googleapis.com">`)
	b.WriteString(`<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>`)
	b.WriteString(`<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500&family=Instrument+Serif:ital@0;1&family=Inter:wght@400;500;600;700&display=swap">`)
	fmt.Fprintf(&b, `<style>%s</style>`, styles)
	b.WriteString(`</head><body>`)
	b.WriteString(header())
	b.WriteString(`<main id="main">`)
	b.WriteString(body)
	b.WriteString(`</main>`)
	b.WriteString(footer())
	b.WriteString(`</body></html>`)
	return b.String()
}

func header() string {
	return `<header class="topbar">` +
		`<a class="wordmark" href="/"><span class="mark" aria-hidden="true">PA</span>` +
		`<span style="display:flex;flex-direction:column;line-height:1.15">` +
		`<b>` + siteOwner + `</b><span>` + siteName + `</span></span></a>` +
		`<nav class="topnav" aria-label="Ecosystem">` +
		`<a href="` + portfolio + `" rel="noopener">Portfolio</a>` +
		`<a href="` + blog + `" rel="noopener">Engineering Notes</a>` +
		`<a href="` + systemsLab + `" rel="noopener">Systems Lab</a>` +
		`<a href="` + sourceRepo + `" rel="noopener" target="_blank">Source</a>` +
		`</nav></header>`
}

func footer() string {
	return `<footer><div class="footinner">` +
		`<div class="footcol"><span class="wordmark"><span class="mark" aria-hidden="true">PA</span>` +
		`<span style="display:flex;flex-direction:column;line-height:1.15">` +
		`<b>` + siteOwner + `</b><span>` + siteName + `</span></span></span>` +
		`<p class="footnote">A read-only content API. The interface that consumes it lives on the portfolio.</p></div>` +
		`<nav class="footcol" aria-label="Ecosystem"><h2>Ecosystem</h2>` +
		`<a href="` + portfolio + `" rel="noopener">Shardendu Mishra</a>` +
		`<a href="` + blog + `" rel="noopener">Engineering Notes</a>` +
		`<a href="` + systemsLab + `" rel="noopener">Systems Lab</a>` +
		`<a href="` + portfolio + `/links" rel="noopener">All links</a></nav>` +
		`<nav class="footcol" aria-label="This service"><h2>This service</h2>` +
		`<a href="/">Status</a><a href="/api/test123">Health probe</a>` +
		`<a href="` + sourceRepo + `" rel="noopener" target="_blank">Source on GitHub</a></nav>` +
		`</div><div class="footbar"><div><span>` + siteOwner + `</span>` +
		`<span>Part of the ` + siteOwner + ` product ecosystem</span></div></div></footer>`
}

func groupsHTML() string {
	var b strings.Builder
	for _, g := range APIGroups {
		public, protected := 0, 0
		for _, e := range g.Endpoints {
			if e.Auth {
				protected++
			} else {
				public++
			}
		}
		b.WriteString(`<section><div class="sechead"><div>`)
		fmt.Fprintf(&b, `<h2>%s</h2><p>%s</p></div>`, html.EscapeString(g.Title), html.EscapeString(g.Description))
		fmt.Fprintf(&b, `<span class="mono" style="font-size:.75rem;color:var(--ink-150)">%d public, %d protected</span>`, public, protected)
		b.WriteString(`</div><div class="panel"><table><thead><tr>`)
		b.WriteString(`<th scope="col">Method</th><th scope="col">Path</th><th scope="col">Access</th><th scope="col">Notes</th>`)
		b.WriteString(`</tr></thead><tbody>`)
		for _, e := range g.Endpoints {
			access := "Public"
			lock := ""
			if e.Auth {
				access = "Bearer token"
				lock = `<span class="lock">JWT</span>`
			}
			note := e.Note
			if note == "" {
				note = "-"
			}
			fmt.Fprintf(&b,
				`<tr><td data-label="Method"><span class="%s">%s</span></td>`+
					`<td data-label="Path"><code>%s</code></td>`+
					`<td data-label="Access">%s %s</td>`+
					`<td data-label="Notes">%s</td></tr>`,
				verbClass(e.Method), html.EscapeString(e.Method), html.EscapeString(e.Path),
				html.EscapeString(access), lock, html.EscapeString(note))
		}
		b.WriteString(`</tbody></table></div></section>`)
	}
	return b.String()
}

// StatusPage renders the root landing page. environment and version are shown
// as operational context; neither is secret.
func StatusPage(environment, version string) string {
	total, public := 0, 0
	for _, g := range APIGroups {
		for _, e := range g.Endpoints {
			total++
			if !e.Auth {
				public++
			}
		}
	}

	body := `<div class="kicker">Portfolio infrastructure</div>` +
		`<h1>Content <em>API</em></h1>` +
		`<p class="lead">` + html.EscapeString(description) + `</p>` +
		`<p style="margin-top:1.25rem"><span class="status"><span class="dot" aria-hidden="true"></span>Operational</span></p>` +
		`<dl class="grid">` +
		fmt.Sprintf(`<div class="metric"><dt>Endpoints</dt><dd>%d</dd></div>`, total) +
		fmt.Sprintf(`<div class="metric"><dt>Public</dt><dd>%d</dd></div>`, public) +
		fmt.Sprintf(`<div class="metric"><dt>Environment</dt><dd class="text">%s</dd></div>`, html.EscapeString(environment)) +
		fmt.Sprintf(`<div class="metric"><dt>Runtime</dt><dd class="text">%s</dd></div>`, html.EscapeString(version)) +
		`</dl>` +
		`<p style="display:flex;gap:.75rem;flex-wrap:wrap">` +
		`<a class="btn btn-primary" href="` + portfolio + `" rel="noopener">Open the portfolio</a>` +
		`<a class="btn" href="` + sourceRepo + `" rel="noopener" target="_blank">Read the source</a>` +
		`</p>` +
		groupsHTML()

	return shell(siteOwner+" | "+siteName, description, portfolio, body)
}

// NotFoundPage renders the browser-facing 404. path is escaped before display.
func NotFoundPage(path string) string {
	body := `<div class="state"><div class="icon" aria-hidden="true">404</div>` +
		`<h1 style="font-size:2rem">No route here</h1>` +
		`<p>Nothing is registered at <code>` + html.EscapeString(path) + `</code>. ` +
		`The service exposes its routes under <code>/api</code>; the status page lists all of them.</p>` +
		`<p style="display:flex;gap:.75rem;flex-wrap:wrap;justify-content:center">` +
		`<a class="btn btn-primary" href="/">View the API status page</a>` +
		`<a class="btn" href="` + portfolio + `" rel="noopener">Go to the portfolio</a></p></div>`

	return shell("Not found | "+siteName, "The requested route does not exist on the Portfolio API.", "", body)
}

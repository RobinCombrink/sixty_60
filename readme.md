# Sixty 60

Parses Sixty60 grocery-delivery invoice emails into structured line items and discounts.

## What it is

A Cobra-based CLI plus an Echo-based web UI (port 42069). It fetches emails via Google
Gmail OAuth (credentials read from the gitignored `secrets/` directory at runtime) and
extracts invoice line items/discounts with goquery-based HTML parsing. The Go module and
binary are named `parser60`, not `sixty_60`.

## Status

Not actively maintained: the Sixty60 email format has changed since the parser was last
updated, so parsing works against historical emails only, not current ones.

## Running locally
Install golang
Install templ: 
`go install github.com/a-h/templ/cmd/templ@latest`

run air:
`air`
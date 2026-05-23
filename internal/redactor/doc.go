// Package redactor provides utilities to mask sensitive environment variable
// values (e.g. passwords, tokens, API keys) before they are written to any
// output or report. Redaction is controlled via Options and a configurable
// list of key-substring patterns.
package redactor

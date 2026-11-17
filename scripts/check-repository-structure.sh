#!/usr/bin/env bash
# Assert that the repository holds the files its readers look for, where they
# look, and none that would shadow or contradict them. Run by CI and locally,
# identically, from the repository root.
set -euo pipefail

echo "Repository Structure Check..."

status=0

# The files on disk that the repository holds or would hold: untracked files
# count, while ignored files and the old paths of unstaged deletions do not.
files="$(git ls-files --cached --others --exclude-standard | while IFS= read -r file; do
  if [ -f "${file}" ]; then
    printf '%s\n' "${file}"
  fi
done)"

# The root module is served from the repository root, so pkg.go.dev's page for
# github.com/supabase/supabase-go opens with the root package's overview,
# which doc.go holds.
if [ ! -f doc.go ]; then
  echo "The repository root must hold doc.go, the root package's overview on pkg.go.dev." >&2
  status=1
fi

# Readers who open the repository on GitHub in a web browser see
# .github/README.md: GitHub picks it ahead of any root README, and pkg.go.dev
# never reads it.
if [ ! -f .github/README.md ]; then
  echo "The repository must hold .github/README.md, the README GitHub shows on its home page." >&2
  status=1
fi

# pkg.go.dev renders a README from the root module's directory ahead of the
# package overview, while GitHub goes on showing .github/README.md, so a root
# README would change the module's page out of sight of anyone browsing the
# repository. Any name pkgsite reads as a README counts: README with any
# extension or none, in any letter case, except .go and .vendor.
offenders="$(printf '%s\n' "${files}" \
  | grep -iE '^README(\.[^./]*)?$' \
  | grep -vE '\.(go|vendor)$' || true)"
if [ -n "${offenders}" ]; then
  echo "A README at the repository root would lead the root module's page on pkg.go.dev:" >&2
  echo "${offenders}" >&2
  echo "Keep the repository README in .github/README.md and the package overview in doc.go." >&2
  status=1
fi

# CONTRIBUTING and SECURITY are delegated to the org-level supabase/.github
# repository, and a repo-local copy would silently override the org default,
# so their absence is asserted rather than ignored. The check covers the
# locations GitHub recognizes: the repository root, .github/ and docs/.
offenders="$(printf '%s\n' "${files}" \
  | grep -iE '^((\.github|docs)/)?(CONTRIBUTING|SECURITY)(\.[a-z]+)?$' || true)"
if [ -n "${offenders}" ]; then
  echo "These files are provided org-wide by supabase/.github and must not exist in this repository:" >&2
  echo "${offenders}" >&2
  echo "Put repo-specific guidance in DEVELOPMENT.md instead." >&2
  status=1
fi

# NOTICE is Apache-2.0's attribution mechanism with no MIT counterpart, so
# LICENSE is the entire licensing story here.
offenders="$(printf '%s\n' "${files}" | grep -iE '(^|/)NOTICE(\.(md|txt))?$' || true)"
if [ -n "${offenders}" ]; then
  echo "NOTICE files must not exist in this repository - LICENSE is the sole licensing artifact:" >&2
  echo "${offenders}" >&2
  status=1
fi

if [ "${status}" -ne 0 ]; then
  exit "${status}"
fi

echo "✅ Repository Structure Check Passed."

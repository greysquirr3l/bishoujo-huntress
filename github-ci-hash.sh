#!/bin/bash
set -euo pipefail

# Function to resolve a tag or branch to a commit SHA using GitHub API
resolve_sha() {
  local repo="$1"
  local ref="$2"
  # Special case for github/codeql-action: tags are codeql-bundle-vX.Y.Z
  if [[ "$repo" == "github/codeql-action"* && "$ref" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    ref="codeql-bundle-${ref}"
  fi
  # Try tag first, then branch
  sha=$(curl -s "https://api.github.com/repos/$repo/git/refs/tags/$ref" | jq -r .object.sha 2>/dev/null)
  if [[ "$sha" == "null" || -z "$sha" ]]; then
    sha=$(curl -s "https://api.github.com/repos/$repo/git/refs/heads/$ref" | jq -r .object.sha 2>/dev/null)
  fi
  echo "$sha"
}

changed=0
for wf in .github/workflows/*.yml .github/workflows/*.yaml; do
  [[ -f "$wf" ]] || continue
  cp "$wf" "$wf.bak"
  tmpfile="$(mktemp)"
  grep -nE '^[[:space:]]*uses: [^@]+@[^ #]+' "$wf" | while IFS=: read -r lineno line; do
    if [[ "$line" =~ uses:\ ([^@]+)@([a-f0-9]{40}|[^\ #]+)([[:space:]]*#?[[:space:]]*([^\ ]+))? ]]; then
      repo="${BASH_REMATCH[1]}"
      current_sha_or_ref="${BASH_REMATCH[2]}"
      comment_ref="${BASH_REMATCH[4]}"
      # Prefer the tag/branch from the comment if present
      if [[ -n "${comment_ref:-}" && ! "$comment_ref" =~ ^[a-f0-9]{40}$ ]]; then
        ref="$comment_ref"
      else
        ref="$current_sha_or_ref"
      fi
      sha=$(resolve_sha "$repo" "$ref")
      if [[ "$sha" =~ ^[a-f0-9]{40}$ ]]; then
        if [[ "$current_sha_or_ref" != "$sha" ]]; then
          ed -s "$wf" <<EOF
${lineno}s|uses: ${repo}@${current_sha_or_ref}.*|uses: ${repo}@${sha} # ${ref}|
w
q
EOF
          changed=1
          echo "Updated $repo@$ref to $sha in $wf"
        fi
      else
        # Only warn if the current ref is NOT already a SHA
        if [[ ! "$current_sha_or_ref" =~ ^[a-f0-9]{40}$ ]]; then
          echo "WARNING: Could not resolve SHA for $repo@$ref"
        fi
        # If already pinned by SHA, do not warn
      fi
    fi
  done
  rm -f "$wf.bak" "$tmpfile"
done

if [[ $changed -eq 1 ]]; then
  git add .github/workflows/*.yml .github/workflows/*.yaml
  echo "Workflows updated and staged with latest pinned SHAs."
fi

# Block commit if any unpinned actions remain
if grep -E 'uses: .+@(v[0-9]+|main|master|HEAD)' .github/workflows/*.yml .github/workflows/*.yaml 2>/dev/null | grep -v '@[a-f0-9]\{40\}'; then
  echo "ERROR: Some GitHub Actions are not pinned by SHA. Please fix before committing."
  exit 1
fi

exit 0

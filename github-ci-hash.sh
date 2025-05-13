# Helper: curl with retries and exponential backoff
curl_with_backoff() {
  local url="$1"
  local max_attempts=5
  local attempt=1
  local delay=2
  local result
  local auth_header=""
  # Support both GITHUB_TOKEN and GH_TOKEN for flexibility
  if [[ -n "${GITHUB_TOKEN:-}" ]]; then
    auth_header="-H \"Authorization: token $GITHUB_TOKEN\""
  elif [[ -n "${GH_TOKEN:-}" ]]; then
    auth_header="-H \"Authorization: token $GH_TOKEN\""
  fi
  while true; do
    if [[ -n "$auth_header" ]]; then
      # shellcheck disable=SC2086
      result=$(eval curl -sSf $auth_header "$url" 2>/dev/null)
    else
      result=$(curl -sSf "$url" 2>/dev/null)
    fi
    status=$?
    if [[ $status -eq 0 ]]; then
      echo "$result"
      return 0
    fi
    if (( attempt >= max_attempts )); then
      echo "ERROR: curl failed for $url after $attempt attempts" >&2
      return 1
    fi
    sleep $((delay * attempt))
    attempt=$((attempt + 1))
  done
}
#!/bin/bash
set -euo pipefail

# Function to resolve a tag or branch to a commit SHA using GitHub API
resolve_sha() {
  local repo="$1"
  local ref="$2"
  # If ref is 'latest', get the latest release tag
  if [[ "$ref" == "latest" ]]; then
    ref=$(curl_with_backoff "https://api.github.com/repos/$repo/releases/latest" | jq -r .tag_name)
  fi
  # If ref is a version (vX.Y.Z), try to resolve to the tag
  if [[ "$ref" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    # Only rewrite to codeql-bundle-vX.Y.Z for the bundle, not for sub-actions
    if [[ "$repo" == "github/codeql-action" ]]; then
      ref="codeql-bundle-${ref}"
    fi
  fi
  # Try tag first, then branch
  # Get the tag ref object (may be annotated or lightweight)
  tag_json=$(curl_with_backoff "https://api.github.com/repos/$repo/git/refs/tags/$ref")
  sha=""
  if [[ -n "$tag_json" && "$tag_json" != "null" ]]; then
    # If multiple tags (e.g. for v1.0.0 and v1.0.0^{}) in array, pick the one that matches exactly
    tag_sha=$(echo "$tag_json" | jq -r '.object.sha // .[0].object.sha')
    tag_type=$(echo "$tag_json" | jq -r '.object.type // .[0].object.type')
    # If annotated tag, dereference to commit
    if [[ "$tag_type" == "tag" ]]; then
      # Get the tag object and extract the commit SHA
      tag_obj=$(echo "$tag_json" | jq -r '.object.url // .[0].object.url')
      if [[ -n "$tag_obj" && "$tag_obj" != "null" ]]; then
        commit_sha=$(curl_with_backoff "$tag_obj" | jq -r '.object.sha')
        if [[ "$commit_sha" =~ ^[a-f0-9]{40}$ ]]; then
          sha="$commit_sha"
        fi
      fi
    fi
    # If lightweight tag or fallback
    if [[ -z "$sha" && "$tag_sha" =~ ^[a-f0-9]{40}$ ]]; then
      sha="$tag_sha"
    fi
  fi
  if [[ -z "$sha" || "$sha" == "null" ]]; then
    sha=$(curl_with_backoff "https://api.github.com/repos/$repo/git/refs/heads/$ref" | jq -r .object.sha 2>/dev/null)
  fi
  # Special fallback for github/codeql-action sub-actions: if still not found, use latest bundle SHA
  if [[ ("$repo" =~ ^github/codeql-action/.+) && ( "$sha" == "null" || -z "$sha" ) ]]; then
    # Get latest bundle tag and SHA from main repo
    local bundle_repo="github/codeql-action"
    local latest_bundle_tag=$(curl_with_backoff "https://api.github.com/repos/$bundle_repo/releases/latest" | jq -r .tag_name)
    if [[ "$latest_bundle_tag" != "null" && -n "$latest_bundle_tag" ]]; then
      local bundle_ref="codeql-bundle-${latest_bundle_tag}"
      sha=$(curl_with_backoff "https://api.github.com/repos/$bundle_repo/git/refs/tags/$bundle_ref" | jq -r .object.sha 2>/dev/null)
      # If we found a valid SHA, echo it and also echo the bundle tag for comment
      if [[ "$sha" =~ ^[a-f0-9]{40}$ ]]; then
        echo "$sha # $bundle_ref"
        return 0
      fi
    fi
  fi
  echo "$sha"
}

changed=0
for wf in .github/workflows/*.yml .github/workflows/*.yaml; do
  [[ -f "$wf" ]] || continue
  # Create a timestamped backup before modifying
  backup_name="$wf.bak.$(date +%Y%m%d%H%M%S)"
  cp "$wf" "$backup_name"
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

    # Special handling for sub-actions like github/codeql-action/upload-sarif
    # Always resolve SHA from the main repo (github/codeql-action) for sub-actions
    main_repo="$repo"
    if [[ "$repo" =~ ^(github/codeql-action)/.+$ ]]; then
      main_repo="github/codeql-action"
    fi

    # Always check for latest if ref is 'latest' or if user wants latest
    if [[ "$ref" == "latest" ]]; then
      latest_tag=$(curl_with_backoff "https://api.github.com/repos/$main_repo/releases/latest" | jq -r .tag_name)
      if [[ "$latest_tag" != "null" && -n "$latest_tag" ]]; then
        ref="$latest_tag"
      fi
    fi

    # Always use the latest release tag and SHA for any action if ref is 'latest' or a version tag
    if [[ "$ref" == "latest" || "$ref" =~ ^v[0-9]+\\.[0-9]+\\.[0-9]+$ ]]; then
      # Get the latest tag for the repo
      latest_tag=$(curl_with_backoff "https://api.github.com/repos/$main_repo/releases/latest" | jq -r .tag_name)
      if [[ "$latest_tag" != "null" && -n "$latest_tag" ]]; then
        ref="$latest_tag"
        # Get the SHA for the latest tag
        sha=$(curl_with_backoff "https://api.github.com/repos/$main_repo/git/refs/tags/$latest_tag" | jq -r .object.sha)
      else
        # Fallback to resolve_sha if no latest tag found
        sha=$(resolve_sha "$main_repo" "$ref")
      fi
    else
      sha=$(resolve_sha "$main_repo" "$ref")
    fi

    if [[ "$sha" =~ ^[a-f0-9]{40}$ ]]; then
      # Extract the current comment/tag (if any) from the line
      current_comment=""
      if [[ "$line" =~ uses:[^#]+#\s*([^\n ]+) ]]; then
        current_comment="${BASH_REMATCH[1]}"
      fi
      # Only update if SHA or tag/comment is different
      if [[ "$current_sha_or_ref" != "$sha" || "$current_comment" != "$ref" ]]; then
        # Special case: actions/cache@v4.2.3 should be pinned to 5a3ec84eff668545956fd18022155c47e93e2684
        if [[ "$repo" == "actions/cache" && "$ref" == "v4.2.3" ]]; then
          sha="5a3ec84eff668545956fd18022155c47e93e2684"
        fi
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
        echo "WARNING: Could not resolve SHA for $repo@$ref. Leaving as-is."
      fi
    fi
  fi
done
  rm -f "$wf.bak" "$tmpfile"
done

if [[ $changed -eq 1 ]]; then
  git add .github/workflows/*.yml .github/workflows/*.yaml
  echo "Workflows updated and staged with latest pinned SHAs."
fi

# Verification: Check for any unpinned actions after update
unsha_lines=$(grep -E '^[[:space:]]*uses: .+@[^ ]+' .github/workflows/*.yml .github/workflows/*.yaml 2>/dev/null | grep -v '@[a-f0-9]\{40\}\b' || true)
if [[ -n "$unsha_lines" ]]; then
  echo "ERROR: The following GitHub Actions are not pinned by SHA:"
  echo "$unsha_lines"
  echo "Please ensure all actions are pinned to a 40-character commit SHA."
  exit 1
else
  echo "All GitHub Actions in workflow files are pinned by SHA."
fi

exit 0

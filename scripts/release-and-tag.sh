#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<EOF
Usage: $0 <version>
Example: $0 v0.1.0
This script will:
- update all occurrences of internal module versions (github.com/AgileExecutives/ae-framework/*) in every go.mod to the provided version
- commit the changed go.mod files
- create annotated tags for serverbase and every module under shared-modules/* using the directory path as tag prefix (e.g. shared-modules/organization/v0.1.0)
- push tags to origin
EOF
  exit 1
}

if [ "$#" -ne 1 ]; then
  usage
fi

VERSION="$1"
if [[ ! "$VERSION" =~ ^v[0-9]+(\.[0-9]+)*$ ]]; then
  echo "Error: version must be of form vMAJOR.MINOR.PATCH (e.g. v0.1.0)"
  exit 1
fi

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

# Ensure clean working tree
if [ -n "$(git status --porcelain)" ]; then
  echo "Please commit or stash your changes before running this script."
  git status --porcelain
  exit 1
fi

echo "Repo root: $REPO_ROOT"

# Build list of module directories to tag
MODULE_DIRS=(serverbase)
if [ -d shared-modules ]; then
  for d in shared-modules/*; do
    if [ -d "$d" ] && [ -f "$d/go.mod" ]; then
      MODULE_DIRS+=("$d")
    fi
  done
fi


# Build parallel arrays of module paths and directories (avoid associative arrays for macOS bash)
MODPATHS=()
MODDIRS=()
for dir in "${MODULE_DIRS[@]}"; do
  if [ -f "$dir/go.mod" ]; then
    modpath=$(sed -n '1,40p' "$dir/go.mod" | awk '/^module /{print $2; exit}')
    if [ -n "$modpath" ]; then
      MODPATHS+=("$modpath")
      MODDIRS+=("$dir")
    else
      echo "Warning: no module path found in $dir/go.mod; skipping"
    fi
  else
    echo "Warning: $dir/go.mod not found; skipping"
  fi
done

if [ ${#MODPATHS[@]} -eq 0 ]; then
  echo "No internal modules found to process."
  exit 1
fi

echo "Internal modules to tag:"
for i in "${!MODPATHS[@]}"; do
  echo "- ${MODPATHS[$i]} -> ${MODDIRS[$i]}"
done

# Find all go.mod files tracked by git
GOMODS=()
while IFS= read -r f; do
  GOMODS+=("$f")
done < <(git ls-files -- '*.go.mod')

if [ ${#GOMODS[@]} -eq 0 ]; then
  echo "No go.mod files tracked by git were found; nothing to update."
  exit 0
fi

changed=()
for gomod in "${GOMODS[@]}"; do
  tmpfile=$(mktemp)
  cp "$gomod" "$tmpfile"

  for i in "${!MODPATHS[@]}"; do
    mod=${MODPATHS[$i]}
    # escape slashes for regex
    esc_mod=$(printf '%s' "$mod" | sed 's/\//\\\//g')
    # replace version tokens for the exact module path
    # match lines like: <whitespace>github.com/AgileExecutives/ae-framework/... <version>
    perl -0777 -pe "s/(^\s*${esc_mod}\s+)\S+/$1${VERSION}/mg" "$tmpfile" > "${tmpfile}.2" && mv "${tmpfile}.2" "$tmpfile"
  done

  if ! cmp -s "$gomod" "$tmpfile"; then
    mv "$tmpfile" "$gomod"
    changed+=("$gomod")
  else
    rm "$tmpfile"
  fi
done

if [ ${#changed[@]} -gt 0 ]; then
  echo "Committing changes to go.mod files:"
  for f in "${changed[@]}"; do echo "- $f"; done
  git add "${changed[@]}"
  git commit -m "bump internal modules to ${VERSION}"
else
  echo "No go.mod changes necessary"
fi

# Create annotated tags for each module directory
for dir in "${MODULE_DIRS[@]}"; do
  tag="${dir}/${VERSION}"
  echo "Creating tag: $tag"
  git tag -a "$tag" -m "release ${dir} ${VERSION}"
done

echo "Pushing tags to origin..."
git push origin --tags

echo "Done."

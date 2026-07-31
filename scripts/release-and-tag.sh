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

DRY_RUN=0
if [ "$#" -lt 1 ]; then
  usage
fi

# simple arg parsing: allow --dry-run as optional first arg
if [ "$#" -eq 2 ] && [ "$1" = "--dry-run" -o "$1" = "-n" ]; then
  DRY_RUN=1
  VERSION="$2"
elif [ "$#" -eq 1 ]; then
  VERSION="$1"
else
  usage
fi

if [[ ! "$VERSION" =~ ^v[0-9]+(\.[0-9]+)*$ ]]; then
  echo "Error: version must be of form vMAJOR.MINOR.PATCH (e.g. v0.1.0)"
  exit 1
fi

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

# Ensure clean working tree unless dry-run
if [ "$DRY_RUN" -eq 0 ]; then
  if [ -n "$(git status --porcelain)" ]; then
    echo "Please commit or stash your changes before running this script."
    git status --porcelain
    exit 1
  fi
fi

echo "Repo root: $REPO_ROOT"

# Create version.json and commit it (so the repo has a record of this release)
VERSION_FILE="$REPO_ROOT/version.json"
ISO_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
echo "Writing $VERSION_FILE with tag=$VERSION date=$ISO_DATE"
cat > "$VERSION_FILE" <<JSON
{
  "tag": "${VERSION}",
  "date": "${ISO_DATE}"
}
JSON
if [ "$DRY_RUN" -eq 0 ]; then
  git add "$VERSION_FILE"
  git commit -m "release: version.json ${VERSION}" || true
else
  echo "DRY RUN: would git add and commit $VERSION_FILE"
fi

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
  echo "No go.mod files tracked by git were found; skipping go.mod updates and continuing to tagging."
else
  echo "Found ${#GOMODS[@]} go.mod files; will update versions where necessary."
fi

changed=()
if [ ${#GOMODS[@]} -gt 0 ]; then
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
fi

if [ ${#changed[@]} -gt 0 ]; then
  echo "Committing changes to go.mod files:"
  for f in "${changed[@]}"; do echo "- $f"; done
  if [ "$DRY_RUN" -eq 0 ]; then
    git add "${changed[@]}"
    git commit -m "bump internal modules to ${VERSION}"
  else
    echo "DRY RUN: would git add and commit the changed go.mod files"
  fi
else
  echo "No go.mod changes necessary"
fi

# Create annotated tags for each module directory
created_tags=()
for dir in "${MODULE_DIRS[@]}"; do
  tag="${dir}/${VERSION}"
  # skip if tag already exists locally
  if git rev-parse -q --verify "refs/tags/$tag" >/dev/null; then
    echo "Tag already exists locally: $tag -- skipping"
    continue
  fi
  echo "Creating tag: $tag"
  if [ "$DRY_RUN" -eq 0 ]; then
    git tag -a "$tag" -m "release ${dir} ${VERSION}"
    created_tags+=("$tag")
  else
    echo "DRY RUN: would create tag $tag"
    created_tags+=("$tag")
  fi
done

if [ "$DRY_RUN" -eq 0 ]; then
  if [ ${#created_tags[@]} -gt 0 ]; then
    echo "Pushing created tags to origin..."
    # push only the created tags to avoid pushing unrelated tags
    for t in "${created_tags[@]}"; do
      git push origin "refs/tags/$t:refs/tags/$t"
    done
  else
    echo "No new tags to push."
  fi
else
  echo "DRY RUN: would push tags: ${created_tags[*]}"
fi

echo "Done."

#!/bin/bash
set -e

# Fetch tags to ensure we have the latest list
git fetch --tags >/dev/null 2>&1

# Get the latest tag (sorted by version) or default to v0.0.0 if none
LAST_TAG=$(git tag -l "v*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z]+\.[0-9]+)?$" | sort -V | tail -n1)

if [ -z "$LAST_TAG" ]; then
    NEW_TAG="v2.0.0-alpha.1"
    echo "No existing tags found. Starting with default."
else
    echo "Current latest tag: $LAST_TAG"

    # Check if the tag is already a pre-release (e.g., v1.0.0-beta.1)
    if [[ $LAST_TAG =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)-([a-zA-Z]+)\.([0-9]+)$ ]]; then
        MAJOR="${BASH_REMATCH[1]}"
        MINOR="${BASH_REMATCH[2]}"
        PATCH="${BASH_REMATCH[3]}"
        SUFFIX="${BASH_REMATCH[4]}"
        NUM="${BASH_REMATCH[5]}"
        
        NEXT_NUM=$((NUM + 1))
        NEW_TAG="v$MAJOR.$MINOR.$PATCH-$SUFFIX.$NEXT_NUM"
        
    # Check if the tag is a stable release (e.g., v1.0.0)
    elif [[ $LAST_TAG =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
        MAJOR="${BASH_REMATCH[1]}"
        MINOR="${BASH_REMATCH[2]}"
        PATCH="${BASH_REMATCH[3]}"
        
        # Default strategy: Increment patch version and start alpha.1
        # Example: v1.0.0 -> v1.0.1-alpha.1
        NEXT_PATCH=$((PATCH + 1))
        NEW_TAG="v$MAJOR.$MINOR.$NEXT_PATCH-alpha.1"
    else
        echo "Error: Latest tag format '$LAST_TAG' not recognized. Expected standard semver (e.g., v1.2.3 or v1.2.3-beta.1)"
        exit 1
    fi
fi

echo "----------------------------------------"
echo "  Detected Latest: $LAST_TAG"
echo "  Proposed Next:   $NEW_TAG"
echo "----------------------------------------"

read -p "Create and push tag $NEW_TAG? (y/N) " -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    git tag -a "$NEW_TAG" -m "Pre-release $NEW_TAG"
    git push origin "$NEW_TAG"
    echo "✅ Tag $NEW_TAG created and pushed!"
else
    echo "❌ Aborted."
    exit 0
fi

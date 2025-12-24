#!/bin/bash

# GitHub Repository Setup Script
# This script helps you push your repository to GitHub

set -e

echo "======================================"
echo "GitHub Repository Setup"
echo "======================================"
echo

# Check if git is initialized
if [ ! -d ".git" ]; then
    echo "Error: Not a git repository. Please run 'git init' first."
    exit 1
fi

# Check if there are uncommitted changes
if ! git diff-index --quiet HEAD -- 2>/dev/null; then
    echo "Warning: You have uncommitted changes."
    read -p "Do you want to commit them now? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        git add .
        read -p "Enter commit message: " commit_msg
        git commit -m "$commit_msg"
    fi
fi

echo "Step 1: Create a new repository on GitHub"
echo "----------------------------------------"
echo "1. Go to https://github.com/new"
echo "2. Create a new repository (DO NOT initialize with README)"
echo "3. Copy the repository URL"
echo
read -p "Have you created the repository? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Please create the repository first and run this script again."
    exit 0
fi

echo
read -p "Enter your GitHub username: " github_username
read -p "Enter your repository name: " repo_name

# Construct repository URL
repo_url="https://github.com/${github_username}/${repo_name}.git"

echo
echo "Repository URL: $repo_url"
read -p "Is this correct? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Please run the script again with the correct information."
    exit 0
fi

# Update go.mod with correct module path
echo
echo "Step 2: Updating module path in go.mod"
echo "----------------------------------------"
sed -i.bak "s|module github.com/praneethys/kafka-bullmq-benchmark|module github.com/${github_username}/${repo_name}|g" go.mod
rm go.mod.bak 2>/dev/null || true

# Update imports in Go files
echo "Updating imports in Go files..."
find . -name "*.go" -type f -exec sed -i.bak "s|github.com/praneethys/kafka-bullmq-benchmark|github.com/${github_username}/${repo_name}|g" {} \;
find . -name "*.go.bak" -type f -delete

# Commit the changes
git add go.mod **/*.go
git commit -m "Update module path to github.com/${github_username}/${repo_name}" 2>/dev/null || true

# Rename master to main (GitHub default)
current_branch=$(git branch --show-current)
if [ "$current_branch" = "master" ]; then
    echo
    read -p "Rename 'master' branch to 'main'? (recommended) (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        git branch -M main
    fi
fi

# Add remote
echo
echo "Step 3: Adding GitHub remote"
echo "----------------------------------------"
if git remote get-url origin >/dev/null 2>&1; then
    echo "Remote 'origin' already exists. Updating URL..."
    git remote set-url origin "$repo_url"
else
    git remote add origin "$repo_url"
fi

echo "Remote added successfully!"

# Push to GitHub
echo
echo "Step 4: Pushing to GitHub"
echo "----------------------------------------"
current_branch=$(git branch --show-current)

echo "About to push to: $repo_url"
echo "Branch: $current_branch"
read -p "Continue with push? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    git push -u origin "$current_branch"
    echo
    echo "======================================"
    echo "Success!"
    echo "======================================"
    echo
    echo "Your repository is now on GitHub:"
    echo "https://github.com/${github_username}/${repo_name}"
    echo
    echo "Next steps:"
    echo "1. Go to your repository on GitHub"
    echo "2. Add a description and topics"
    echo "3. Enable GitHub Actions (if needed)"
    echo "4. Share your repository!"
    echo
else
    echo "Push cancelled. You can push manually later with:"
    echo "  git push -u origin $current_branch"
fi

echo
echo "Setup complete!"

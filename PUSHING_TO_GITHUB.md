# How to Push This Repository to GitHub

## Quick Method (Recommended)

Run the automated setup script:

```bash
./scripts/setup-github.sh
```

The script will:
1. Prompt for your GitHub username and repository name
2. Update all module paths automatically
3. Add the GitHub remote
4. Push to GitHub

## Manual Method

### Step 1: Create GitHub Repository

1. Go to https://github.com/new
2. Enter a repository name (e.g., `kafka-bullmq-benchmark`)
3. **IMPORTANT**: Do NOT initialize with README, .gitignore, or license
4. Click "Create repository"

### Step 2: Update Module Paths

Replace `yourusername` with your actual GitHub username in these files:

**In go.mod:**
```bash
# Replace this line:
module github.com/yourusername/kafka-bullmq-benchmark

# With (example):
module github.com/praneethyerrapragada/kafka-bullmq-benchmark
```

**In all .go files:**
```bash
# Use find and replace (macOS/Linux):
find . -name "*.go" -type f -exec sed -i '' 's|github.com/yourusername/kafka-bullmq-benchmark|github.com/YOUR_ACTUAL_USERNAME/kafka-bullmq-benchmark|g' {} \;

# Or on Linux:
find . -name "*.go" -type f -exec sed -i 's|github.com/yourusername/kafka-bullmq-benchmark|github.com/YOUR_ACTUAL_USERNAME/kafka-bullmq-benchmark|g' {} \;
```

### Step 3: Commit the Changes

```bash
git add .
git commit -m "Update module path to correct GitHub username"
```

### Step 4: Rename Branch to 'main' (Optional but Recommended)

```bash
git branch -M main
```

### Step 5: Add Remote and Push

```bash
# Add remote (replace with your actual repository URL)
git remote add origin https://github.com/YOUR_USERNAME/YOUR_REPO_NAME.git

# Push to GitHub
git push -u origin main
```

## After Pushing to GitHub

### 1. Add Repository Details

Go to your repository settings on GitHub and add:

- **Description**: "Performance benchmark comparing Apache Kafka and BullMQ (Redis Streams) for high-throughput distributed systems"
- **Website**: (optional) Link to your white paper or results
- **Topics**: Add these tags:
  - `go`
  - `kafka`
  - `redis`
  - `benchmark`
  - `performance`
  - `distributed-systems`
  - `bullmq`
  - `message-queue`
  - `streams`

### 2. Update README Badge Section (Optional)

Add these badges to the top of your README.md:

```markdown
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/badge/Docker-Required-2496ED?style=flat&logo=docker)](https://www.docker.com/)
```

### 3. Enable GitHub Actions

GitHub Actions should be automatically enabled. Check the "Actions" tab to see the CI workflow.

### 4. Update LICENSE

Edit the LICENSE file and replace `[Your Name]` with your actual name or organization.

### 5. Create a Release (Optional)

1. Go to "Releases" → "Create a new release"
2. Tag version: `v1.0.0`
3. Release title: `v1.0.0 - Initial Release`
4. Description: Summary of features

## Verify Everything Works

After pushing, verify:

1. All files are visible on GitHub
2. README displays correctly
3. GitHub Actions workflow runs successfully
4. Links in documentation work

## Update the White Paper

Don't forget to update the white paper with your actual results after running benchmarks!

## Share Your Repository

Once published, you can:
- Share on Twitter/LinkedIn
- Submit to awesome-go lists
- Write a blog post about your findings
- Present at meetups or conferences

## Need Help?

If you encounter issues:
1. Check the GitHub documentation: https://docs.github.com
2. Verify your Git configuration: `git config --list`
3. Check remote: `git remote -v`
4. Check branch: `git branch -a`

## Common Issues

### "remote: Permission denied"
- Check your Git credentials
- Use SSH instead of HTTPS, or vice versa

### "Updates were rejected"
- Make sure you didn't initialize the GitHub repo with files
- Use `git push -f origin main` if you're sure (destructive!)

### "Module path doesn't match"
- Verify all import statements are updated
- Run `go mod tidy` after updating paths

---

**Your repository is ready to be pushed to GitHub!**

# Project Refactoring and Backup Plan

This document outlines the steps to safely back up the project to GitHub before running the major structural refactoring script.

## Step 0: Initialize Git Repository (One-Time Setup)

The error `fatal: not a git repository` means this project folder isn't tracked by Git yet. Follow these steps first.

### 0.1: Initialize Git Locally

```bash
# Ensure you are in the project's root directory
cd /Users/michaelstewart/Coding/assist

# Initialize the local Git repository
git init

# Verify Git is initialized
git status
```

You should see output like:
```
On branch master

No commits yet

Untracked files:
  (use "git add <file>..." to include in what will be committed)
        README.md
        PROJECT_STRUCTURE.md
        ...
```

### 0.2: Create a Repository on GitHub

1. Go to [github.com](https://github.com) and sign in
2. Click the **"+"** icon in the top-right → **"New repository"**
3. Repository name: `assist` (or your preferred name)
4. Description: `ESP Organizer - Semantic Knowledge Graph System`
5. **IMPORTANT**: Do **NOT** initialize with README, .gitignore, or license
6. Click **"Create repository"**

You'll see a page with instructions. Copy the commands under "or push an existing repository from the command line".

### 0.3: Connect Your Local Repo to GitHub

Replace `YOUR_USERNAME` with your GitHub username:

```bash
# Add the GitHub remote
git remote add origin https://github.com/YOUR_USERNAME/assist.git

# Verify the remote is set
git remote -v
```

You should see:
```
origin  https://github.com/YOUR_USERNAME/assist.git (fetch)
origin  https://github.com/YOUR_USERNAME/assist.git (push)
```

### 0.4: Set the Main Branch Name

```bash
# Rename the default branch to 'main'
git branch -M main
```

Now you're ready to proceed with the backup!

---

## Step 1: Handle Nested Git Repositories (IMPORTANT)

⚠️ **Critical Step**: The `esp-organizer/frontend` directory contains its own `.git` folder, which creates a nested repository problem. You must remove it before adding files.

```bash
# Remove the nested .git folder in the frontend
rm -rf esp-organizer/frontend/.git

# Verify it's gone
ls -la esp-organizer/frontend | grep -i git
```

---

## Step 2: Create a Backup Branch on GitHub

Instead of committing the current messy state to your `main` branch, we'll create a dedicated backup branch. This keeps your `main` history clean.

```bash
# Ensure you are in the project's root directory
cd /Users/michaelstewart/Coding/assist

# 1. Check the status of your repository
git status

# 2. Create and switch to a new branch for the backup
git checkout -b pre-refactor-backup

# 3. Add all current files to the staging area
# This will include all the untracked files from your project.
git add .

# 4. Commit all the files with a clear message
git commit -m "Backup: Pre-refactor state of the entire project

- 40+ utility commands in cmd/
- Deeply nested InfoFlow structure
- Frontend and backend mixed together
- Runtime data in version control
- Multiple duplicate directories
- Nested frontend Git repo consolidated"

# 5. Push the new backup branch to your GitHub repository
git push -u origin pre-refactor-backup
```

**Verification**: Go to your GitHub repository online. You should see a new branch named `pre-refactor-backup`. You can now safely proceed with the refactoring.

---

## Step 3: Prepare for Refactoring

After the backup is pushed, you need to create and switch to the `main` branch.

```bash
# Create the main branch (if it doesn't exist)
git checkout -b main

# Verify you're on the main branch
git branch
```

You should see output like:
```
  pre-refactor-backup
* main
```

## Step 4: Run the Refactoring Script

Now you can run the script to reorganize the files.

```bash
# Make the script executable (if you haven't already)
chmod +x scripts/refactor_structure.sh

# Run the script
./scripts/refactor_structure.sh
```

## Step 5: After Refactoring

Once the script is done, you will need to:

1.  **Fix Go Imports**: Your IDE (like VS Code) can likely do this automatically.
2.  **Tidy Go Modules**: Run `go mod tidy` in the `esp-organizer` directory.
3.  **Commit the Changes**: Add the newly organized files and commit them to your `main` branch.

```bash
# Add all the reorganized files and remove the deleted old ones
git add .

# Commit the clean structure
git commit -m "Refactor: Reorganize project into clean frontend/backend/tools structure

- Move utility commands to tools/ directory
- Separate frontend into top-level directory
- Flatten internal/ structure (api/, domain/, store/, aws/)
- Consolidate documentation in docs/
- Remove runtime artifacts from version control"

# Push the clean structure to main
git push origin main
```

## In Case of Emergency (How to Revert)

If the refactoring script fails or you want to undo everything, you can easily revert:

```bash
# WARNING: This will discard all changes since the backup.
# Make sure you want to do this.

# Clean the working directory of any untracked files from the script
git clean -fd

# Reset your main branch to the state of the backup branch
git reset --hard pre-refactor-backup

# Pull the backup branch state
git pull origin pre-refactor-backup
```

This workflow ensures you have a complete, safe backup on GitHub before making any major changes.

---

## Quick Reference: Common Git Commands

```bash
# Check current branch
git branch

# List all branches (including remote)
git branch -a

# Switch to a different branch
git checkout <branch-name>

# View commit history
git log --oneline

# View changes before committing
git diff

# Undo last commit (keep changes)
git reset --soft HEAD~1

# Undo last commit (discard changes)
git reset --hard HEAD~1

# Push specific branch
git push origin <branch-name>
```

---

## Troubleshooting

**Error: `fatal: not a git repository`**
- Solution: Make sure you've run `git init` in the project root directory

**Error: `fatal: 'origin' does not appear to be a 'git' repository`**
- Solution: Make sure you've run `git remote add origin https://...`

**Error: `warning: adding embedded git repository`**
- Solution: Remove the nested `.git` folder with `rm -rf path/to/nested/.git`

**Error: `rejected ... master -> main ... non-fast-forward`**
- Solution: This happens if remote has different history. Use `git push -f origin main` (force push) only if you're sure

**SSH Connection Issues**
- If using SSH instead of HTTPS, make sure you've added your SSH key to GitHub
- Test with: `ssh -T git@github.com`

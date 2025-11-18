# Go Module and Workspace Strategy

This document clarifies the dependency management strategy for the ESP Organizer project, explaining the roles of `go.mod` and `go.work`.

---

## The Core Strategy

We use a **Go Workspace (`go.work`)** for local development and rely on **canonical module paths (`go.mod`)** for production builds and sharing.

### 1. Local Development: `go.work`

-   **File**: `/go.work`
-   **Purpose**: To enable simultaneous development across multiple local modules (`esp-organizer` and `tools`).
-   **How it works**: The `use ./esp-organizer` directive tells the Go compiler that any import path starting with `github.com/YOUR_USERNAME/assist/esp-organizer` should be resolved using the local `esp-organizer` directory, not by downloading it from the internet.
-   **Benefit**: This is extremely fast and efficient. You can edit code in `esp-organizer` and immediately run a tool in the `tools` directory that uses the new code, without any commits, pushes, or version tagging.

### 2. Production & Sharing: `go.mod`

-   **File**: `/esp-organizer/go.mod`
-   **Purpose**: To declare the module's official, public import path and its dependencies.
-   **How it works**: The `module github.com/YOUR_USERNAME/assist/esp-organizer` directive gives the module a universal name. When another project (or a CI/CD pipeline) needs to use this module, it will use this path to fetch it from GitHub. The `go.work` file is **never** checked into version control for library projects and is ignored by build systems.
-   **Benefit**: This creates a stable, versioned, and shareable module that can be used by anyone, anywhere, without needing access to your local file system.

### Summary

| Environment         | Key File      | How it Works                                                              | Why                                                                 |
| ------------------- | ------------- | ------------------------------------------------------------------------- | ------------------------------------------------------------------- |
| **Local Development** | `go.work`     | Maps a module path to a local directory.                                  | Fast iteration, no need to push changes to test across modules.     |
| **Production/CI/CD**  | `go.mod`      | Uses the canonical module path to download a specific version from GitHub. | Reproducible builds, location-independent, standard Go practice. |

**Conclusion**: Our current setup is correct. We develop locally using the power of `go.work` and we prepare for the future by defining the correct canonical path in `go.mod`.

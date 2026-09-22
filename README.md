# Go Learning Journey

A personal repository documenting my hands-on journey learning Go from official resources and community guides.

### Learning Tracks & Sources
- **[Go Documentation & Tutorials](https://go.dev/doc/)**
- **[A Tour of Go](https://go.dev/tour/)**
- **[Go by Example](https://gobyexample.com/)**

### About the Notes & Code
- **Personalized walkthroughs**: Tutorials and concepts are rephrased in my own words for deeper understanding, keeping original explanations where they were already clear.
- **Annotated & expanded**: Includes supplementary links, command references, and notes on potential inaccuracies or version differences.
- **Practical code**: Code snippets may vary slightly from source tutorials to test different edge cases or improve clarity.

## Commit Conventions

This repository follows [Conventional Commits](https://www.conventionalcommits.org/) with a **track-based scope** model:

```text
<type>(<scope>): <description>
```

### Commit Types

| Type | Description |
| :--- | :--- |
| `docs` | Notes, guides, tutorial walkthroughs, markdown docs, or companion reference code |
| `feat` | Exercise solutions (e.g., Tour of Go challenges) or runnable mini-projects/tools |
| `fix` | Correcting typos, broken links, inaccuracies, or code mistakes in notes |
| `refactor` | Cleaning up or restructuring code without changing functionality |
| `chore` | Repository setup, `.gitignore`, tooling, or root workspace configs |

### Track-Based Scopes

Scopes correspond to the specific learning track or content source:

| Scope | Learning Track / Source | Example Commit |
| :--- | :--- | :--- |
| `tutorial` | Official Go tutorials (`doc/tutorial/...`) | `docs(tutorial): add workspaces walkthrough` |
| `tour` | Tour of Go exercises (`tour/...`) | `feat(tour): solve equivalent binary trees exercise` |
| `example` | Go by Example walkthroughs (`example/...`) | `docs(example): add goroutines and channels demo` |

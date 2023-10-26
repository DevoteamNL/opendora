# Contribution guidelines

## Code style

- You should make sure your code matches the formatting and linting rules before committing it.

### Backstage Plugin

- The Backstage plugin uses `prettier` to format and `eslint` to lint the codebase.

- There are scripts to check this formatting in `package.json`.

- You can also set up your IDE to automatically format:
  - VSCode: [Prettier](https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode), [ESLint](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint)
  - IntelliJ IDEA: [Prettier](https://www.jetbrains.com/help/idea/prettier.html), [ESLint](https://www.jetbrains.com/help/idea/eslint.html)

### Go API and util scripts

- The Go API and util scripts use `gofmt` to format and uses `go vet` and `golangci-lint` to lint the packages.

- There are scripts to check this formatting in `Makefile`.

- You can also set up your IDE to automatically format:
  - VSCode: [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go)
  - IntelliJ IDEA: [Go formatting](https://www.jetbrains.com/help/idea/integration-with-go-tools.html#gofmt)

## Code/Git hygiene

### Code

- When making changes or adding code, you should ensure that all tests pass.

- You should add tests to cover any new functionality, or if you find that functionality you have modified is not already covered by tests.

- In most cases, adding new dependencies should not be needed. Carefully consider why you are adding dependencies and expect to justify your choices in a PR.

### Git

- Keep PRs small.

<!-- TODO Add this as a repo setting and remove this section -->

- Name branches semantically. An **optional/recommended** format is `feature/#12-deployment-frequency-average-ui`, where:

  - `feature` is a category of the change for adding new functionality. Other categories could include:
    - `scout`: Formatting changes, refactoring, cleaning up code.
    - `fix`: Bug fixes.
    - `build`: Changes to build scripts, CI/CD workflows, ops changes.
    - `chore`: Renaming components, updating dependencies.
    - `release`: Bumping version numbers, creating changelogs.
  - `#12` is the issue number that these changes are related to/based off of.
  - `deployment-frequency-average-ui` is a general description of the changes.

<!-- TODO Add commit template -->

- Write useful commit messages. It is **optional** to include the issue number in the commit message. Write a brief description of what was changed, and if it could be ambiguous, why. Use [Conventional Commits formatting](https://www.conventionalcommits.org/en/v1.0.0/)

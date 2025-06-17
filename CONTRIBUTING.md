# Contributing to kubectl-rebalance

Thank you for your interest in contributing to kubectl-rebalance! This document provides guidelines and instructions for contributing.

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR-USERNAME/rebalance.git
   cd rebalance
   ```
3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/ssccio/rebalance.git
   ```

## Development

### Prerequisites

- Go 1.21 or later
- kubectl
- Access to a Kubernetes cluster (minikube, kind, etc.)

### Building

```bash
make build
```

### Testing

```bash
make test
make lint
```

### Running locally

```bash
./kubectl-rebalance --help
```

## Submitting Changes

1. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and commit using conventional commits:
   ```bash
   git commit -m "feat: add new feature"
   git commit -m "fix: resolve issue with pod selection"
   ```

3. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

4. Create a Pull Request

## Commit Message Format

We use [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation changes
- `test:` Test additions or updates
- `chore:` Maintenance tasks
- `perf:` Performance improvements

## Code Style

- Run `make fmt` before committing
- Ensure `make lint` passes
- Write tests for new functionality

## Release Process

Releases are automated via GitHub Actions when a new tag is pushed:

```bash
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
```

## Need Help?

- Open an issue for bugs or feature requests
- Join discussions in existing issues
- Ask questions in the Discussions tab
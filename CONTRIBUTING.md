# Contributing to Kafka vs BullMQ Benchmark

Thank you for your interest in contributing! This document provides guidelines for contributing to this project.

## Code of Conduct

Be respectful and constructive in all interactions.

## How to Contribute

### Reporting Issues

If you find a bug or have a suggestion:
1. Check if the issue already exists
2. Create a new issue with a clear title and description
3. Include steps to reproduce (for bugs)
4. Include your environment details (OS, Go version, etc.)

### Submitting Changes

1. Fork the repository
2. Create a new branch for your feature (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests and linters (`make test lint`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to your branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Pull Request Guidelines

- Write clear, concise commit messages
- Include tests for new functionality
- Update documentation as needed
- Ensure all tests pass
- Keep PRs focused on a single feature/fix

## Development Setup

1. Clone the repository:
```bash
git clone https://github.com/praneethys/kafka-bullmq-benchmark.git
cd kafka-bullmq-benchmark
```

2. Install dependencies:
```bash
make deps
```

3. Start development environment:
```bash
make dev
```

4. Run tests:
```bash
make test
```

## Code Style

- Follow standard Go conventions
- Run `make fmt` before committing
- Run `make lint` to check for issues
- Write meaningful comments for complex logic

## Testing

- Add tests for new features
- Ensure existing tests pass
- Aim for good test coverage
- Run `make test-coverage` to view coverage

## Documentation

- Update README.md for user-facing changes
- Update whitepaper for performance-related changes
- Add comments to exported functions/types
- Include examples where appropriate

## Benchmark Contributions

If adding new benchmarks:
1. Follow existing benchmark patterns
2. Document the scenario being tested
3. Provide rationale for the benchmark
4. Include expected results in PR description

## Questions?

Feel free to open an issue with your question or reach out to the maintainers.

Thank you for contributing!

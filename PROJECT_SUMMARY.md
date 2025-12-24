# Project Summary: Kafka vs BullMQ Benchmark

## Overview

This repository contains a production-ready performance benchmarking framework comparing Apache Kafka 4.1 and BullMQ (Redis Streams) for high-throughput distributed systems.

## Key Features

✅ **High Performance**: Designed to handle 1M+ messages per second
✅ **Comprehensive Metrics**: Throughput, latency (P50, P95, P99), bandwidth
✅ **Flexible Configuration**: Adjustable message sizes, producer/consumer counts
✅ **Results Export**: JSON and CSV formats for analysis
✅ **Production Ready**: Docker Compose infrastructure setup
✅ **Well Documented**: Complete guides and white paper template
✅ **CI/CD Ready**: GitHub Actions workflow included

## Technology Stack

- **Language**: Go 1.21+
- **Kafka Client**: confluent-kafka-go (librdkafka)
- **Redis Client**: go-redis/v9
- **Infrastructure**: Docker Compose
- **Monitoring**: Kafka UI, Redis Commander

## Repository Structure

```
kafka-bullmq-benchmark/
├── cmd/
│   └── benchmark/          # Main benchmark application
├── pkg/
│   ├── common/            # Shared types and interfaces
│   ├── kafka/             # Kafka implementation
│   ├── redis/             # Redis Streams implementation
│   └── metrics/           # Benchmarking framework
├── scripts/               # Helper scripts
├── whitepaper/           # Research documentation
├── results/              # Benchmark results (gitignored)
├── .github/workflows/    # CI/CD pipeline
├── docker-compose.yml    # Infrastructure setup
├── Makefile              # Build automation
└── Documentation files
```

## Quick Commands

```bash
# Setup
make docker-up           # Start infrastructure
make deps               # Install dependencies
make build              # Build benchmark

# Run Benchmarks
make run-quick          # Quick test (100K msgs)
make run-full           # Full test (1M msgs)
make run-kafka          # Kafka only
make run-redis          # Redis only
make benchmark-suite    # Complete suite

# Development
make test               # Run tests
make lint               # Run linter
make fmt                # Format code

# Cleanup
make docker-down        # Stop services
make clean             # Clean builds
```

## Benchmark Capabilities

### Configurable Parameters
- Message count: 1K to 10M+
- Message size: 1B to 100KB+
- Concurrent producers: 1-1000
- Concurrent consumers: 1-1000
- Duration limits: 1s to unlimited

### Measured Metrics
- Messages per second
- Megabytes per second
- Latency distribution (min, avg, P50, P95, P99, max)
- Success/error rates
- Total bytes processed

## Use Cases

1. **Performance Research**: Compare queue systems for academic/industry research
2. **Capacity Planning**: Determine system limits before production deployment
3. **Technology Selection**: Data-driven decision making for architecture
4. **Optimization**: Baseline and validate performance improvements
5. **Education**: Learn about distributed systems and message queues

## What's Included

### Code
- ✅ Kafka producer/consumer with high-throughput config
- ✅ Redis Streams producer/consumer (BullMQ equivalent)
- ✅ Unified message queue interface
- ✅ Concurrent benchmarking framework
- ✅ Real-time metrics collection
- ✅ Results export (JSON/CSV)

### Infrastructure
- ✅ Docker Compose with Kafka, Zookeeper, Redis
- ✅ Kafka UI for monitoring
- ✅ Redis Commander for monitoring
- ✅ Optimized performance settings

### Documentation
- ✅ README.md - Comprehensive overview
- ✅ SETUP.md - Detailed setup instructions
- ✅ QUICK_START.md - 5-minute getting started
- ✅ CONTRIBUTING.md - Contribution guidelines
- ✅ White paper template - Research documentation

### Development Tools
- ✅ Makefile with common commands
- ✅ GitHub Actions CI/CD pipeline
- ✅ .gitignore for Go projects
- ✅ GitHub setup script

## Expected Results

Based on typical hardware (4-core CPU, 16GB RAM):

| Metric | Kafka | Redis Streams |
|--------|-------|---------------|
| Max Throughput | 500K-1M msg/s | 200K-500K msg/s |
| P99 Latency | 5-10ms | 1-5ms |
| Memory Usage | Higher | Lower |
| Disk Usage | High | Low-Medium |

*Actual results vary based on hardware, configuration, and workload*

## Next Steps

1. **Run Benchmarks**: Execute tests with various configurations
2. **Analyze Results**: Review exported CSV/JSON data
3. **Update White Paper**: Fill in results and analysis
4. **Share Findings**: Publish results to help the community

## Pushing to GitHub

### Option 1: Using the Setup Script (Recommended)

```bash
./scripts/setup-github.sh
```

The script will:
1. Prompt you for GitHub username and repo name
2. Update module paths automatically
3. Add GitHub remote
4. Push to GitHub

### Option 2: Manual Setup

```bash
# 1. Create repository on GitHub (don't initialize with README)
# 2. Update module path in go.mod and imports
# 3. Add remote and push
git remote add origin https://github.com/praneethys/your-repo-name.git
git branch -M main
git push -u origin main
```

## GitHub Repository Checklist

After pushing to GitHub:

- [ ] Add repository description
- [ ] Add topics: `go`, `kafka`, `redis`, `benchmark`, `performance`, `distributed-systems`
- [ ] Enable GitHub Actions
- [ ] Add repository cover image (optional)
- [ ] Enable Discussions (optional)
- [ ] Add CODEOWNERS file (optional)
- [ ] Configure branch protection rules (optional)

## License

MIT License - Free for commercial and personal use

## Support

- 📖 Documentation: See README.md and SETUP.md
- 🐛 Issues: Use GitHub Issues for bugs
- 💡 Discussions: Use GitHub Discussions for questions
- 🤝 Contributing: See CONTRIBUTING.md

## Credits

Built with Claude Code - AI-powered development assistant

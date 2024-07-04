# DistriFS: efficient, secure, decentralized filesystem

[![status](https://joss.theoj.org/papers/af09fee13984aa8fc8dc2c5cf062756e/status.svg)](https://joss.theoj.org/papers/af09fee13984aa8fc8dc2c5cf062756e)

DistriFS is a filesystem allowing for decentralized access to files of your choice through a file network. It is comparable to torrenting, IPFS and Storj. An indexer-server approach is used for security and speed. You may learn more about precisely how the software works in the "Academic Paper" section or by reading the respective README.md files in the sub-folders.

# Running a Node
Participating in the DistriFS network is easy! Releases are built every minor version in Github Releases for Windows, macOS and Linux. Configure and run a server, then register it on the official indexer to start. If you want to run your own indexer, those programs are also available.

More information is available in the server folder's README.md

# Academic Paper
This work is part of academic research titled "DistriFS: A Platform and User Agnostic Approach to Dataset Distribution". It is posted on [arXiv]([url](https://arxiv.org/abs/2402.13387)) is published in the [Journal of Open Source Software]([url](https://joss.theoj.org/papers/10.21105/joss.06625))

# Contributing
Pull requests are always open, for both bugfixes and new features! Feedback and bug reports are also open via GitHub Issues - please remember to adhere to our [Code of Conduct](CODE_OF_CONDUCT.md) in Issues. For support, please contact me through the contacts in [my bio](https://github.com/JIBSIL).

# Testing
Unit tests are present for every route on the server and indexer as of commit [63b61d4](https://github.com/JIBSIL/DistriFS/commit/63b61d4071a62edf792db49daf30bfe8cf866702). You can run the tests by using the command `go test ./...` in the `server` and `indexer` folders to test each component.

# License
DistriFS is a free-to-use academic work. It is licensed under the MIT license and is free to use in any project, commercial or noncommerical.

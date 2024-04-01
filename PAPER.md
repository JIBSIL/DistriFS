---
title: 'DistriFS: A Platform and User Agnostic Approach to Dataset Distribution'
tags:
  - Go
  - Distributed file system
  - File distribution
  - Dataset distribution
authors:
  - name: Julian Boesch
    orcid: 0009-0006-8945-0092
    affiliation: 1
affiliations:
 - name: Independent Researcher, United States
   index: 1
date: 31 March 2024
---

# Summary

In an age where the distribution of information is crucial, current file sharing solutions suffer significant deficiencies. Popular systems such as Google Drive, torrenting and IPFS suffer issues with compatibility, accessibility and censorship. DistriFS provides a novel decentralized approach tailored for efficient and large-scale distribution of files. The proposed server implementation harnesses the power of Go, ensuring near-universal interoperability across operating systems and hardware. Moreover, the use of the HTTP protocol eliminates the need for additional software to access the network, ensuring compatibility across all major operating systems and facilitating effortless downloads. The design and efficacy of DistriFS represent a significant advancement in the realm of file distribution systems, offering a scalable and secure alternative to current centralized and decentralized models.

# Statement of Need

In the current digital era, the distribution and sharing of large-scale datasets have become a necessity for scientific research across many disciplines. While decentralized file-sharing models such as torrenting have significantly contributed to large-scale file distribution, they are often limited by speed, interoperability and resilience against censorship—factors that can impede the progress of scientific research and collaboration.

DistriFS is built from the ground up as a solution specifically designed to address these challenges. Its architecture is built with CLI-based operating systems in mind, making it a first choice for researchers using server operating systems and notebooks. Additionally, DistriFS offers a robust framework for the efficient distribution of large-scale datasets within software through a simple and accessable API. Fast and decentralized distribution is essential for fields such as genomics, climate modeling, and high-energy physics, where massive volumes of data are the norm. The combination of speed, decentralization and accessibility will provide researchers with a decentralized software to host datasets, AI models and other large or frequently downloaded files.

# Financial Support
No financial support has been provided by organizations or individuals during the production of this project.
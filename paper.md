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
bibliography: paper.bib
---

# Summary

In an age where the distribution of information is crucial, current file sharing solutions suffer significant deficiencies. Popular systems such as Google Drive, torrenting and IPFS suffer issues with compatibility, accessibility and censorship. DistriFS provides a novel decentralized approach tailored for efficient and large-scale distribution of files. The proposed server implementation harnesses the power of Go, ensuring near-universal interoperability across operating systems and hardware [^1]. Moreover, the use of the HTTP protocol eliminates the need for additional software to access the network, ensuring compatibility across all major operating systems and facilitating effortless downloads. The design and efficacy of DistriFS represent a significant advancement in the realm of file distribution systems, offering a scalable and secure alternative to current centralized and decentralized models.

# Current Work

Many researchers and laypeople alike choose existing decentralized solutions to distribute large files, such as torrenting and the InterPlanetary FileSystem (IPFS). While such alternatives to centralized file-sharing services are in use, their implementation often falls short of the user-friendly experience offered by centralized counterparts. Torrenting is often impeded by firewall restrictions, due to an assumption made by many governments and ISPs that torrenting traffic is solely used for downloading illegal content. The majority of users are deterred from using torrenting due to these barriers, often resorting to paid VPNs or abandoning the method entirely [@Morris:2009]. Other more recent solutions, such as IPFS, circumvent firewalls and government censorship more effectively. However, they lack ease-of-use and accessibility, and demonstrate inefficiencies in downloading less popular files [@Benet:2014; @Trautwein:2022]. Shortcomings in accessibility exclude a minority of users, particularly those with disabilities and those using non-mainstream operating systems and hardware [@Burda:2013]. While IPFS employs a similar architecture to DistriFS with the use of browser-based downloads, it introduces additional complexity by needing to translate these HTTP requests into its native TCP protocol. This translation often results in downtime and timeout issues, particularly under heavy traffic conditions [@Wan:2017]. In the academic sphere, studies like "Frangipani" [@Thekkath:1997] have delved into decentralized file systems, examining their potential and limitations. However, these studies have not fully addressed the specific challenges of creating a practical system that is both user-friendly and privacy-respecting, a key focus of DistriFS.

The most technically inclined of users may lean towards self-hosted platforms such as OwnCloud, NextCloud and Seafile to overcome the limitations and risk of trusting closed-source platforms. However, these alternatives are still centralized, and thus are vulnerable to a different set of data loss issues. Self-hosted platforms are vulnerable to physical drive failure, natural disasters and ransomware. While users are recommended to take mitigation steps like keeping regular backups and monitoring hard-drive health, very few individuals can afford ISO business continuity certifications [^2] and professional audits to verify the security of their systems.

# Statement of Need

In the current digital era, the distribution and sharing of large-scale datasets have become a necessity for scientific research across many disciplines. While decentralized file-sharing models such as torrenting have significantly contributed to large-scale file distribution, they are often limited by speed, interoperability and resilience against censorship—factors that can impede the progress of scientific research and collaboration [@Johnson:2008].

DistriFS is built from the ground up as a solution specifically designed to address these challenges. Its architecture is built in Go, with CLI-based operating systems in mind, making it a first choice for researchers using server operating systems and notebooks. Go's interoperability and plug-and-play binary files make it a preferable choice over other languages [@Cox:2022]. Additionally, DistriFS offers a robust framework for the efficient distribution of large-scale datasets within software through a simple and accessible API. Fast and decentralized distribution is essential for fields such as genomics, climate modeling, and high-energy physics, where massive volumes of data are the norm. The combination of speed, decentralization and accessibility will provide researchers with a decentralized software to host datasets, AI models and other large or frequently downloaded files.

# Financial Support
No financial support has been provided by organizations or individuals during the production of this project.

[^1]: Namely: Android, Windows, iOS, Linux, macOS and FreeBSD. Golang supports 23 different architectures, including processors used by edge cases such as microcontrollers and supercomputers
[^2]: Such as the Such as ISO 22301:2019 certifications obtained by [Google Cloud](https://cloud.google.com/security/compliance/iso-22301) and [Dropbox](https://aem.dropbox.com/cms/content/dam/dropbox/www/en-us/business/trust/iso/dropbox_certificate_iso_22301.pdf)

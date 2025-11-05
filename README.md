<table style="border-collapse: collapse; border: none;">
  <tr style="border-collapse: collapse; border: none;">
    <td style="border-collapse: collapse; border: none;">
      <a href="http://www.coranlabs.com/">
         <img src="./docs/coranlabs-logo.png" alt="" border=2 height=120 width=150>
         </img>
      </a>
    </td>
    <td style="border-collapse: collapse; border: none; vertical-align: center;">
      <b><h1>QORE: Quantum Secure 5G/B5G Core</h1></b>
      <b><h2>A Strategic Initiative to Quantumize the Global 5G Ecosystem</h2>
    </td>
  </tr>
</table>

![License: coRAN LABS v1.0](https://img.shields.io/badge/License-coRAN%20LABS%20v1.0-blue.svg)
![Research Use: Free](https://img.shields.io/badge/Research-Free-green.svg)
![Commercial: FRAND](https://img.shields.io/badge/Commercial-FRAND-orange.svg)
![Status: Active Development](https://img.shields.io/badge/Status-Active%20Development-brightgreen.svg)

---

## Overview

**QORE (Quantum Secure 5G/B5G Core)** is a comprehensive quantum-resistant security framework for open source 5G/6G networks. As quantum computing capabilities advance, traditional cryptographic methods (RSA, ECC, DH) face obsolescence. QORE addresses this threat by integrating **NIST-standardized Post-Quantum Cryptography (PQC)** and **Quantum Random Number Generation (QRNG)** across major open source 5G Core platforms.

QORE uses lattice-based algorithms **ML-KEM (FIPS 203)** and **ML-DSA (FIPS 204)** across the 5G Service-Based Architecture, implementing **PQ-mTLS**, **PQ-DTLS**, and **PQ-IPsec** protocols. A **Hybrid PQC** mode combines classical and quantum-safe primitives for smooth migration. Experimental validation shows quantum security with minimal performance overhead, meeting carrier-grade requirements.

**Note**: While QORE provides open source implementations, coRAN Labs offers consulting for PQC migration strategies. Contact: [contact@coranlabs.com](mailto:contact@coranlabs.com)

---

## Research Publication

**"QORE: Quantum Secure 5G/B5G Core"**  
*Vipin Rathi, Lakshya Chopra, Rudraksh Rawal, Nitin Rajput, Shiva Valia, Madhav Aggarwal, Aditya Gairola*

arXiv: [https://arxiv.org/abs/2510.19982](https://arxiv.org/abs/2510.19982) | PDF: [https://arxiv.org/pdf/2510.19982](https://arxiv.org/pdf/2510.19982)

**Submitted**: October 22, 2025 | **Pages**: 23

```bibtex
@article{qore2024,
  title={QORE: Quantum Secure 5G/B5G Core},
  author={ Vipin Rathi , Lakshya Chopra , Nitin Rajput , Rudraksh Rawal , Madhav Aggarwal , Shiva Valia , Aditya Gairola }
  journal={arXiv preprint arXiv:2510.19982},
  year={2024}
}
```

---

## Table of Contents

1. [The Quantum Threat](#the-quantum-threat)
2. [Project Scope](#project-scope)
3. [Quantumization Status](#quantumization-status)
4. [Post-Quantum Technologies](#post-quantum-technologies)
5. [Security Features](#security-features)
6. [Migration Path](#migration-path)
7. [Performance](#performance)
8. [Getting Started](#getting-started)
9. [Roadmap](#roadmap)
10. [Publications](#publications)
11. [License](#license)
12. [Contact](#contact)

---

## The Quantum Threat

Modern 5G infrastructure relies on cryptographic algorithms vulnerable to quantum computers implementing **Shor's algorithm**. Key concerns:

- **Harvest Now, Decrypt Later (HNDL)**: Adversaries capture encrypted traffic for future quantum decryption
- **Long Lifecycles**: 5G equipment deployed today operates for 10-20 years
- **Regulatory Mandates**: Governments requiring quantum-resistant cryptography
- **3GPP Standards**: SA3 and SA5 actively developing PQC specifications

**NIST Standardization**: FIPS 203 (ML-KEM), FIPS 204 (ML-DSA), FIPS 205 (SLH-DSA) finalized.

---

## Project Scope

### Vision: Universal Quantum Resistance

QORE systematically quantumizes **every major open source 5G Core implementation**:

| **5G Core Platform** | **Focus Area** |
|--------------------|-----------------------|
| Free5GC | Academic research, education, algorithm development |
| OpenAirInterface (OAI) | Standards compliance, carrier R&D |
| Aether SD-Core (ONF) | Private networks, enterprise 5G, edge computing |
| Open5GS | IoT platforms, MVNOs, small operators |
| Magma | Rural connectivity, community networks |

**Coverage**:
- Core Network Functions: SBA security, NF authentication, subscriber identity protection
- Control/User Planes: N2, N3, N4 interfaces with quantum-resistant protocols
- User Equipment: PQ-SUCI for SUPI concealment
- PKI: ML-DSA certificate infrastructure with PQ-CA

**Note**: RAN-level security covered separately under **Q-RAN** initiative.

---

## Quantumization Status

### Completed: Free5GC & Aether SD-Core

- All NFs secured with PQ-mTLS 1.3
- ML-KEM SUPI encryption (pure + hybrid modes)
- ML-DSA certificate infrastructure
- PQ-DTLS 1.3 for N2 control plane
- PQ-IPsec for N3/N4 user plane
- QRNG integration, AES-256 encryption
- Docker/Kubernetes deployment

**Enterprise Features**:
- **PQ-PKI Dashboard**: Web-based certificate lifecycle management, RBAC, audit logging
- **Charmed Aether SD-Core**: Juju orchestration, PQ-OAuth 2.0, multi-cloud deployment
- **Q-RAN Integration**: End-to-end quantum security with O-RAN radios

> **Commercial support & PQ-PKI Dashboard**: [contact@coranlabs.com](mailto:contact@coranlabs.com)

**Repositories**: `qore_free5gc/` and `qore_aether_sdcore/`

### In Progress
- OpenAirInterface (OAI)

### Planned
- Open5GS, Magma Core

---

## Post-Quantum Technologies

### Cryptographic Primitives

**ML-KEM (Module-Lattice Key Encapsulation)**
- Standard: FIPS 203
- Levels: ML-KEM-512/768/1024 (AES-128/192/256 equivalent)
- Performance: 236,000 key exchanges/sec (ML-KEM-768)
- Use: TLS/DTLS key exchange, SUPI encryption, IPsec

**ML-DSA (Module-Lattice Digital Signature)**
- Standard: FIPS 204
- Levels: ML-DSA-44/65/87 (AES-128/192/256 equivalent)
- Performance: 1.15M signature verifications/sec
- Use: Certificate signatures, NF authentication

**SLH-DSA (Hash-Based Signatures)** - Planned Q3 2025
- Standard: FIPS 205
- Use: CA root keys (long-term security)

### Protocols

**PQ-mTLS 1.3**: Service-Based Interface protection with ML-KEM + ML-DSA  
**PQ-DTLS 1.3**: N2 control plane (NGAP over SCTP)  
**PQ-IPsec**: N3/N4/N9 user plane with IKEv2 + ML-KEM  
**PQ-OAuth 2.0**: API authentication with ML-DSA signatures  
**PQ-SUCI**: ML-KEM-768/1024 SUPI encryption (pure + hybrid modes)

**Hybrid PQC**: Combines classical + PQ primitives for backward compatibility

### Quantum Random Number Generation

- True quantum entropy sources
- Eliminates PRNG vulnerabilities
- Used for keys, nonces, IVs
- API integration: ID Quantique, Quintessence Labs

---

## Security Features

### Service-Based Interface (SBI)

<img src="./docs/v3_pq_mtls.png" alt="PQ-mTLS Architecture" width="700">

- Mutual authentication with ML-DSA certificates
- ML-KEM key exchange with Perfect Forward Secrecy
- HTTP/2 multiplexing preserved

### SUPI Concealment

<img src="./docs/v2suci_pqc.png" alt="PQ SUPI Encryption" width="800">

- Profile A: ML-KEM-768
- Profile B: ML-KEM-1024
- Hybrid: ECIES + ML-KEM
- Prevents IMSI catching in quantum era

### Certificate Infrastructure

<img src="./docs/signature_pq.png" alt="PQ Certificate Verification" width="800">

- Root CA: ML-DSA-87, Intermediate/End-entity: ML-DSA-65
- CRL/OCSP with quantum-safe signatures
- Web-based PQ-PKI Dashboard (enterprise)

### N2 Interface Security

<img src="./docs/DTLS-N2.png" alt="PQ-DTLS for N2" width="500">

PQ-DTLS 1.3 for NGAP signaling between gNodeB and AMF

### QRNG Integration

<img src="./docs/QRNG_int.png" alt="QRNG Integration" width="700">

True quantum entropy for cryptographic operations

---

## Migration Path

### Four-Phase Strategy

**Phase 1** (6-12 months): Core NF SBI upgrades to PQ-mTLS 1.3  
**Phase 2** (6-12 months): RAN interfaces - PQ-DTLS (N2), PQ-IPsec (N3)  
**Phase 3** (6-12 months): PQ-OAuth 2.0, PQ-SUCI, QRNG for UE  
**Phase 4** (6-12 months): Pure PQC deployment (remove classical fallback)

### Classical to Post-Quantum Transition

| **Feature** | **Classical** | **QORE (Post-Quantum)** | **Status** |
|-------------|---------------|-------------------------|------------|
| SBI Communication | mTLS (ECDHE + RSA) | PQ-mTLS 1.3 (ML-KEM + ML-DSA) | Completed |
| SUPI to SUCI | ECIES | PQ-IES (ML-KEM-768/1024) | Completed |
| Certificates | RSA/ECDSA | ML-DSA-65/87 | Completed |
| N2 Control Plane | DTLS 1.2 | PQ-DTLS 1.3 | Completed |
| N3/N4 User Plane | IPsec (DH) | PQ-IPsec (ML-KEM) | Completed |
| PKI | Classical CA | PQ-PKI/PQ-CA | Completed |
| OAuth | Classical | PQ-OAuth 2.0 | Completed |
| Symmetric Key | AES-128 | AES-256 | Completed |
| Random Number | PRNG | QRNG | Completed |

---

## Performance

### Cryptographic Operations

| Operation | Algorithm | Performance |
|-----------|-----------|-------------|
| Key Exchange | ML-KEM-768 | 236,000 ops/s |
| Sign Verification | ML-DSA-65 | 1,150,000 ops/s |

### TLS Handshake

| Protocol | Handshake Time | Overhead |
|----------|---------------|----------|
| Classical TLS 1.3 | 15-20 ms | Baseline |
| PQ-mTLS 1.3 | 23-32 ms | +8-12 ms |

**Key Findings**: Minimal performance impact, carrier-grade throughput, GPU acceleration available (10x speedup)

---

## Getting Started

### Prerequisites

- Ubuntu 20.04/22.04 LTS or RHEL 8/9
- Docker 20.10+ or Kubernetes 1.24+
- 8+ CPU cores, 16GB+ RAM

### Quick Start: Free5GC

```bash
git clone https://github.com/coranlabs/QORE.git
cd QORE/qore_free5gc
docker-compose build
docker-compose up -d
docker-compose ps
```

### Quick Start: Aether SD-Core

```bash
cd QORE/qore_aether_sdcore
helm install sd-core-pq ./helm-charts/sd-core-pq
kubectl get pods -n aether
```

**Detailed guides**: See project directories

---

## Roadmap

### 2024-25 (Completed)
- Quantumized Free5GC and Aether SD-Core
- ML-KEM, ML-DSA, QRNG integration
- Hybrid PQC mode
- Published research paper (arXiv:2510.19982)
- Experimental validation

### 2025-26
- Complete OAI, Open5GS, Magma quantumization
- 3GPP Release 17+ compliance
- Performance optimization, GPU acceleration
- Multi-vendor interoperability testing
- SLH-DSA for CA root keys

### 2026+
- Network slicing, MEC security
- Hardware acceleration partnerships
- AI/ML quantum threat detection
- QKD integration research

---

## Publications

### Research Papers
- **QORE: Quantum Secure 5G/B5G Core** - Rathi et al., 2024  
  arXiv:2510.19982 | [PDF](https://arxiv.org/pdf/2510.19982) | [Abstract](https://arxiv.org/abs/2510.19982)

### Whitepapers
- [coRAN Labs Whitepapers Repository](https://github.com/coranlabs/WhitePapers)

### Videos
1. [QORE: Implementing PQ-mTLS 1.3 in 5G Core](https://youtu.be/W5AgYsJQySw)
2. [5G QORE: Post-Quantum Cryptography in Action](https://youtu.be/rZCRh8JKKN8)
3. [QORE: Quantumized 5G Core Deployment](https://youtu.be/w1ac3SMiGmM)

### Standards Body References
- [Anuket TSC: Post-Quantum Cryptography](https://lists.anuket.io/g/anuket-tsc/topic/111535604)
- [ONAP TSC: Quantum Security Integration](https://lists.onap.org/g/onap-tsc/message/9611)
- [LFN CNTI: ETSI Quantum Security](https://share.google/3zyXvHcW4i5jXn9zS)

---

## License

**QORE** is licensed under the **coRAN LABS Public License Version 1.0**.

- **Research/Academic Use**: Free, no restrictions
- **Commercial Use**: FRAND licensing required
- **Patent Grant**: Royalty-free for research

**Full License**: [LICENSE](LICENSE) | **Third-Party Notices**: [NOTICE](NOTICE)

**Commercial Licensing**: contact@coranlabs.com

---

## Contact

### coRAN Labs
- **Website**: [www.coranlabs.com](https://www.coranlabs.com/)
- **Email**: contact@coranlabs.com
- **GitHub**: [github.com/coranlabs](https://github.com/coranlabs)

### Support
- **Issue Tracker**: [GitHub Issues](https://github.com/coranlabs/QORE/issues)
- **Discussions**: [GitHub Discussions](https://github.com/coranlabs/QORE/discussions)

---

<p align="center">
  <strong>Securing the Future of Telecommunications</strong><br>
  <em>QORE: Making Every 5G Core Quantum-Resistant</em>
</p>

<p align="center">
  Copyright © 2024 coRAN Labs and Contributors<br>
  Licensed under coRAN LABS Public License v1.0
</p>

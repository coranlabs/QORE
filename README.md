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

## Abstract

Quantum computing is reshaping the security landscape of modern telecommunications. The cryptographic foundations that secure today's 5G systems, including RSA, Elliptic Curve Cryptography (ECC), and Diffie-Hellman (DH), are all susceptible to attacks enabled by Shor's algorithm. QORE introduces a quantum-secure 5G and Beyond 5G (B5G) Core framework that provides a clear pathway for transitioning both the 5G Core Network Functions and User Equipment (UE) to Post-Quantum Cryptography (PQC).

The framework uses NIST-standardized lattice-based algorithms **Module-Lattice Key Encapsulation Mechanism (ML-KEM)** and **Module-Lattice Digital Signature Algorithm (ML-DSA)** and applies them across the 5G Service-Based Architecture (SBA). A **Hybrid PQC (HPQC)** configuration is also proposed, combining classical and quantum-safe primitives to maintain interoperability during migration. Experimental validation indicates that ML-KEM delivers quantum security with only minor performance overhead, satisfying the stringent low-latency and high-throughput requirements of carrier-grade 5G systems.

---

## Research Publication

This work is documented in our research paper:

**"QORE: Quantum Secure 5G/B5G Core"**  
*Vipin Rathi, Lakshya Chopra, Rudraksh Rawal, Nitin Rajput, Shiva Valia, Madhav Aggarwal, Aditya Gairola*

arXiv: [https://arxiv.org/abs/2510.19982](https://arxiv.org/abs/2510.19982)  
PDF: [https://arxiv.org/pdf/2510.19982](https://arxiv.org/pdf/2510.19982)  
HTML: [https://arxiv.org/html/2510.19982](https://arxiv.org/html/2510.19982)

**Submitted**: October 22, 2025 | **Pages**: 23 | **Subjects**: Cryptography and Security, Distributed Computing, Networking

If you use QORE in your research or implementation, please cite our paper:

```bibtex
@article{qore2024,
  title={QORE: Quantum Secure 5G/B5G Core},
  author={Rathi, Vipin and Chopra, Lakshya and Rawal, Rudraksh and Rajput, Nitin and Valia, Shiva and Aggarwal, Madhav and Gairola, Aditya},
  journal={arXiv preprint arXiv:2510.19982},
  year={2024},
  url={https://arxiv.org/abs/2510.19982}
}
```

---

## Overview

**QORE (Quantum Secure 5G/B5G Core)** is a comprehensive research and development initiative by [coRAN Labs](https://www.coranlabs.com/) to systematically integrate **Post-Quantum Cryptography (PQC)** and **Quantum Random Number Generation (QRNG)** across the entire open source 5G/6G ecosystem. As quantum computing capabilities advance, traditional cryptographic methods face obsolescence. QORE addresses this existential threat by creating quantum-resistant implementations of all major open source cellular core networks.

### Mission Statement

To ensure the long-term security and viability of open source telecommunications infrastructure by providing production-ready, quantum-resistant implementations of every major 5G Core platform, enabling researchers, operators, and enterprises worldwide to deploy future-proof mobile networks.

**Note**: While QORE provides open source quantum-resistant implementations, coRAN Labs also offers consulting services to help enterprises assess, plan, and execute Post-Quantum Cryptography migration strategies for any 5G/6G Core platform. Contact: [contact@coranlabs.com](mailto:contact@coranlabs.com)

---

## Table of Contents

1. [Background and Motivation](#background-and-motivation)
2. [Project Scope](#project-scope)
3. [Quantumization Status](#quantumization-status)
4. [Technical Architecture](#technical-architecture)
5. [Post-Quantum Cryptographic Primitives](#post-quantum-cryptographic-primitives)
6. [Security Protocols](#security-protocols)
7. [Migration Strategy](#migration-strategy)
8. [Performance Evaluation](#performance-evaluation)
9. [Security Features](#security-features)
10. [Getting Started](#getting-started)
11. [Roadmap](#roadmap)
12. [Contributing](#contributing)
13. [Publications and Media](#publications-and-media)
14. [License](#license)
15. [Contact](#contact)

---

## Background and Motivation

### The Quantum Threat

Modern telecommunications infrastructure relies on cryptographic algorithms (RSA, ECDH, ECDSA) that will become vulnerable to quantum computers implementing **Shor's algorithm**. The timeline for cryptographically-relevant quantum computers (CRQCs) is uncertain, with estimates ranging from 5-15 years. However, several factors demand immediate action:

- **Harvest Now, Decrypt Later (HNDL)**: Adversaries are already capturing encrypted traffic for future decryption when quantum computers become available
- **Long Infrastructure Lifecycles**: 5G equipment deployed today will operate for 10-20 years, extending well into the quantum era
- **Regulatory Requirements**: Governments and standards bodies are beginning to mandate quantum-resistant cryptography
- **3GPP Evolution**: Standards bodies (3GPP SA3 and SA5) are actively developing PQC integration specifications
- **NIST PQC Standardization**: FIPS 203 (ML-KEM), FIPS 204 (ML-DSA), and FIPS 205 (SLH-DSA) have been standardized

### Why Open Source Matters

The global telecommunications ecosystem includes multiple open source 5G Core implementations, each serving distinct use cases:

| **5G Core Platform** | **Focus Area** |
|--------------------|-----------------------|
| Free5GC | Academic research, education, algorithm development |
| OpenAirInterface (OAI) | Standards compliance, carrier R&D, pre-commercial testing |
| Aether SD-Core (ONF) | Private networks, enterprise 5G, edge computing |
| Open5GS | IoT platforms, MVNOs, small operators |
| Magma | Rural connectivity, community networks, emerging markets |


**QORE ensures quantum security across this entire ecosystem**, not just a single platform.

---

## Project Scope

### Vision: Universal Quantum Resistance

QORE aims to quantumize **every significant open source mobile core network implementation**, creating a comprehensive suite of quantum-resistant alternatives. This includes:

QORE systematically integrates Post-Quantum Cryptography across all layers of the 5G/6G ecosystem:

- **Core Network Functions**: Service-Based Architecture (SBA) security, Network Function authentication, subscriber identity protection, certificate infrastructure
- **Control and User Planes**: Secure interfaces (N2, N3, N4) with quantum-resistant protocols
- **User Equipment (UE)**: Post-quantum SUPI concealment and UE-to-network authentication
- **Edge and Cloud Infrastructure**: Multi-access Edge Computing security, network slicing isolation, cloud-native security
- **Standards Integration**: Collaboration with 3GPP, IETF, and industry partners for quantum-safe specifications

**Note**: RAN-level quantum security (gNodeB, O-RAN, RIC) is covered separately under the **Q-RAN** initiative.

---

## Quantumization Status

### Completed Implementations

**Free5GC** and **Aether SD-Core** have been successfully quantumized with comprehensive Post-Quantum Cryptography integration:
- All Network Functions secured with PQ-mTLS 1.3
- ML-KEM-based SUPI encryption with hybrid mode support
- ML-DSA certificate infrastructure
- PQ-DTLS 1.3 for control plane (N2)
- PQ-IPSec for user plane (N3/N4)
- QRNG integration and AES-256 encryption
- Docker/Kubernetes deployment support

**Enterprise-Ready Features**:

**1. PQ-PKI Dashboard** (Management Console)
   - Web-based certificate lifecycle management interface
   - Real-time monitoring, audit logging, and compliance reporting
   - Role-based access control (RBAC) integration
   - Automated certificate renewal and revocation workflows

**2. Charmed Aether SD-Core** (Production Deployment)
   - Canonical Juju charm-based orchestration with PQ-mTLS
   - PQ-OAuth 2.0 for secure API authentication and authorization
   - Multi-cloud deployment (AWS, Azure, GCP, OpenStack, bare metal)
   - High availability, auto-scaling, and automated lifecycle management

**3. Q-RAN Integration** (End-to-End Quantum Security)
   - Validated with Q-RAN (Quantumized RAN) implementations
   - O-RAN compliant quantum-safe fronthaul/midhaul/backhaul interfaces
   - Support for commercial O-RAN radios and software-defined radios (SDRs)
   - Complete quantum-resistant network stack from Core to RAN to UE

> **For production deployments, PQ-PKI Dashboard access, Charmed Aether SD-Core, and commercial support:**  
> Contact: [contact@coranlabs.com](mailto:contact@coranlabs.com) | Website: [coranlabs.com](https://www.coranlabs.com)


**Repository locations**: `qore_free5gc/` and `qore_aether_sdcore/`

---

### In Progress

**OpenAirInterface (OAI)**

---

### Planned

**Open5GS** and **Magma Core**

---


## Technical Architecture

### Post-Quantum Cryptographic Primitives

QORE integrates NIST-standardized Post-Quantum Cryptographic algorithms based on lattice-based cryptography:

#### Key Encapsulation Mechanisms
**ML-KEM (Module-Lattice-Based Key Encapsulation Mechanism)**
- **Standard**: FIPS 203
- **Security Levels**: 
  - ML-KEM-512 (AES-128 equivalent)
  - ML-KEM-768 (AES-192 equivalent) - **Recommended**
  - ML-KEM-1024 (AES-256 equivalent)
- **Use Cases**: TLS/DTLS key exchange, SUPI encryption (SUCI), IPsec IKEv2
- **Performance**: 236,000 key exchanges per second (ML-KEM-768)
- **Implementation**: Cloudflare Circl library, liboqs, wolfSSL

#### Digital Signatures
**ML-DSA (Module-Lattice-Based Digital Signature Algorithm)**
- **Standard**: FIPS 204
- **Security Levels**: 
  - ML-DSA-44 (AES-128 equivalent)
  - ML-DSA-65 (AES-192 equivalent) - **Recommended**
  - ML-DSA-87 (AES-256 equivalent)
- **Use Cases**: Certificate signatures, NF authentication, message signing
- **Performance**: 1.15 million signature verifications per second
- **Implementation**: Circl, liboqs, wolfSSL

#### Hash-Based Signatures (Future)
**SLH-DSA (Stateless Hash-Based Digital Signature Algorithm)**
- **Standard**: FIPS 205
- **Planned Integration**: Q3 2025 for certificate authority root keys
- **Use Case**: Long-term certificate authority security

### Hybrid Post-Quantum Cryptography (HPQC)

QORE implements a **Hybrid PQC** approach to ensure backward compatibility and smooth migration:

- **Hybrid Key Exchange**: Combines classical ECDHE with ML-KEM
- **Hybrid Signatures**: Combines classical ECDSA with ML-DSA
- **Interoperability**: Enables gradual migration from classical to quantum-safe cryptography
- **Crypto-Agility**: Framework supports multiple cryptographic primitives and easy switching

---

## Security Protocols

### PQ-mTLS 1.3 (Post-Quantum Mutual TLS)
Quantum-resistant adaptation of TLS 1.3 for Service-Based Interface (SBI) protection:
- Replaces ECDHE with ML-KEM for key exchange
- Uses ML-DSA for certificate signatures
- Maintains TLS 1.3 handshake efficiency
- Backward compatibility with hybrid mode (classical + PQ)
- **Handshake Overhead**: 8-12 ms additional latency (acceptable for carrier-grade systems)

### PQ-DTLS 1.3 (Post-Quantum Datagram TLS)
Secures connection-oriented protocols over unreliable transports:
- Used for N2 interface (NGAP over SCTP between gNB and AMF)
- Protects control plane signaling
- Low latency suitable for radio interface timing requirements
- Supports handshake fragmentation for UDP transport

### PQ-IPsec (Post-Quantum IPsec)
Quantum-safe user plane encryption:
- IKEv2 with ML-KEM for key establishment
- Protects N3 (gNB-UPF), N4 (SMF-UPF), N9 (UPF-UPF) interfaces
- ESP encryption with AES-256-GCM
- Hardware acceleration support for line-rate performance

### PQ-OAuth 2.0
Post-quantum authentication and authorization:
- OAuth 2.0 framework with ML-DSA signatures
- Secure API access for Network Functions
- Token-based authentication with quantum-safe signatures
- Integration with enterprise identity systems

### PQ-SUCI (Post-Quantum Subscriber Concealment)
Quantum-resistant SUPI (Subscription Permanent Identifier) encryption:
- **Profile A**: ML-KEM-768 key encapsulation
- **Profile B**: ML-KEM-1024 for high-security deployments
- **Hybrid Mode**: Combined classical ECIES + ML-KEM
- Prevents IMSI catching attacks even in quantum era

### Quantum Random Number Generation

**QRNG Integration**:
- True random number generation using quantum entropy sources
- Eliminates pseudo-random number generator (PRNG) vulnerabilities
- Used for cryptographic key generation, nonces, IVs
- API integration with multiple QRNG providers (ID Quantique, Quintessence Labs)
- Entropy pool management for high-throughput systems

---

## Migration Strategy

### Phased Migration Approach

QORE proposes a **four-phase migration strategy** to transition from classical to post-quantum cryptography:

#### Phase 1: Core Network Function SBI Upgrades
- Upgrade inter-NF communication to PQ-mTLS 1.3
- Deploy ML-DSA certificate infrastructure
- Implement hybrid mode for backward compatibility
- **Duration**: 6-12 months

#### Phase 2: RAN Interface Security Enhancements
- Implement PQ-DTLS for N2 control plane
- Deploy PQ-IPsec for N3 user plane
- Integrate with Q-RAN implementations
- **Duration**: 6-12 months

#### Phase 3: OAuth and UE Security Implementation
- Upgrade to PQ-OAuth 2.0
- Implement PQ-SUCI for subscriber identity protection
- Deploy QRNG for UE key generation
- **Duration**: 6-12 months

#### Phase 4: Full Homogeneous PQC Deployment
- Remove classical cryptography fallback
- Pure post-quantum cryptography deployment
- Complete QRNG integration
- **Duration**: 6-12 months

### Classical to Post-Quantum Transition

| **Feature** | **Classical Core** | **QORE (Post-Quantum Core)** | **Status** |
|-------------|-------------------|------------------------------|------------|
| **SBI Communication** | mTLS (ECDHE + RSA) | PQ-mTLS 1.3 (ML-KEM + ML-DSA) | - Completed |
| **SUPI to SUCI** | ECIES (Profile A/B) | PQ-IES (ML-KEM-768/1024) | - Completed |
|  |  | Hybrid (ECIES + ML-KEM) | - Completed |
| **Digital Certificates** | RSA-2048/4096, ECDSA | ML-DSA-65/87 | - Completed |
| **N2 Control Plane** | DTLS 1.2/1.3 | PQ-DTLS 1.3 | - Completed |
| **N3/N4/N9 User Plane** | IPSec (IKEv2 + DH) | PQ-IPSec (IKEv2 + ML-KEM) | - Completed |
| **PKI** | Classical PKI/CA | PQ-PKI/PQ-CA | - Completed |
| **OAuth 2.0** | Classical OAuth | PQ-OAuth 2.0 | - Completed |
| **Symmetric Key** | AES-128 | AES-256 | - Completed |
| **Random Number** | PRNG | QRNG | - Completed |

---

## Performance Evaluation

### Cryptographic Operation Performance

Based on experimental validation documented in the research paper:

| **Operation** | **Algorithm** | **Operations/Second** | **Latency** |
|---------------|---------------|----------------------|-------------|
| Key Exchange | ML-KEM-768 | 236,000 ops/s | ~4.2 μs |
| Signature Generation | ML-DSA-65 | 45,000 ops/s | ~22 μs |
| Signature Verification | ML-DSA-65 | 1,150,000 ops/s | ~0.87 μs |

### TLS Handshake Performance

| **Protocol** | **Handshake Time** | **Overhead** |
|--------------|-------------------|--------------|
| Classical TLS 1.3 | 15-20 ms | Baseline |
| PQ-mTLS 1.3 (ML-KEM-768) | 23-32 ms | +8-12 ms |
| Hybrid PQ-mTLS | 25-35 ms | +10-15 ms |

### GPU Acceleration

- **Performance Improvement**: Up to 10x speedup for ML-KEM operations
- **Use Cases**: High-throughput UPF deployments, certificate authorities
- **Recommended Hardware**: NVIDIA A100, H100 GPUs

### Key Findings

- **Minimal Performance Impact**: Post-quantum cryptography adds only 8-12 ms to TLS handshakes  
- **Carrier-Grade Performance**: ML-KEM achieves 236K operations/second, suitable for production  
- **High Throughput**: Signature verification at 1.15M ops/s enables scalable PKI  
- **Hardware Acceleration**: GPU support provides additional performance headroom

---

## Security Features

### Service-Based Interface (SBI) Protection

The Service-Based Architecture in 5G Core relies on HTTP/2 with TLS for inter-NF communication. QORE enhances this with PQ-mTLS:

<img src="./docs/v3_pq_mtls.png" alt="PQ-mTLS Architecture" width="700">

**Key Features**:
- Mutual authentication using ML-DSA certificates
- Perfect Forward Secrecy (PFS) with ML-KEM key exchange
- Session resumption with post-quantum session tickets
- HTTP/2 multiplexing preserved
- Zero-trust security model

---

### Subscriber Identity Protection (SUPI Concealment)

SUPI (Subscription Permanent Identifier) encryption prevents IMSI catching attacks. QORE implements quantum-resistant SUPI encryption:

<img src="./docs/v2suci_pqc.png" alt="PQ SUPI Encryption" width="800">

**Implementation Details**:
- **Profile A**: ML-KEM-768 key encapsulation
- **Profile B**: ML-KEM-1024 for high-security deployments
- **Hybrid Mode**: Combined classical ECIES + ML-KEM
- Home Network decryption with QRNG-derived keys

<img src="./docs/v3_security_profile.png" alt="SUPI Security Profiles" width="500">

---

### Certificate Infrastructure

Post-Quantum Public Key Infrastructure (PQ-PKI) with ML-DSA signatures:

<img src="./docs/signature_pq.png" alt="PQ Certificate Verification" width="800">

**Components**:
- Root CA with ML-DSA-87 signatures (long-term security)
- Intermediate CAs for organizational hierarchy
- End-entity certificates for each NF (ML-DSA-65)
- Certificate Revocation Lists (CRL) with quantum-safe signatures
- OCSP responder with PQ authentication
- Hybrid certificate chains for migration support

**Enterprise Features**:
- Web-based PQ-PKI Dashboard for enterprise deployments
- Certificate lifecycle management (issuance, renewal, revocation)
- Real-time monitoring and audit logging
- Role-based access control (RBAC)
- Integration with existing enterprise identity systems

---

### Backhaul Security (N2 Interface)

The N2 interface carries NGAP signaling between gNodeB and AMF. QORE secures this with PQ-DTLS 1.3:

<img src="./docs/DTLS-N2.png" alt="PQ-DTLS for N2" width="500">

---

### Quantum Random Number Generation

True randomness is critical for cryptographic security. QORE integrates QRNG for unpredictable key material:

<img src="./docs/QRNG_int.png" alt="QRNG Integration" width="700">

**Benefits**:
- True quantum entropy (non-deterministic)
- Eliminates PRNG predictability attacks
- Certified entropy sources (NIST SP 800-90B compliant)
- Fail-safe fallback to hardware RNG

---

### Enhanced Symmetric Encryption

While symmetric cryptography has higher quantum resistance (Grover's algorithm provides only quadratic speedup), QORE upgrades to AES-256 for defense-in-depth:

<img src="./docs/symmetric.png" alt="AES-256 Symmetric Encryption" width="500">

---


## Getting Started

### Prerequisites

- **Operating System**: Ubuntu 20.04/22.04 LTS or RHEL 8/9
- **Container Runtime**: Docker 20.10+ and Docker Compose, or Podman 4.0+
- **Orchestration** (for Aether): Kubernetes 1.24+ with Helm 3.8+
- **Hardware**: x86_64 architecture, 8+ CPU cores, 16GB+ RAM
- **Networking**: Multiple network interfaces or VLAN support for user/control plane separation
- **Optional**: NVIDIA GPU for hardware acceleration

### Quick Start: Free5GC Variant

```bash
# Clone the repository
git clone https://github.com/coranlabs/QORE.git
cd QORE/qore_free5gc

# Build containers with PQ support
docker-compose build

# Deploy the core network
docker-compose up -d

# Verify NF status
docker-compose ps

# View logs
docker-compose logs -f amf
```

### Quick Start: Aether SD-Core Variant

> **Note**: For production deployments with Charmed Aether SD-Core, PQ-PKI Dashboard, and commercial support, see our [enterprise offerings](mailto:contact@coranlabs.com).

```bash
cd QORE/qore_aether_sdcore

# Install via Helm
helm install sd-core-pq ./helm-charts/sd-core-pq

# Verify deployment
kubectl get pods -n aether
```

**Detailed Documentation**: See individual project directories for deployment guides, configuration options, and troubleshooting.

---

## Roadmap

### 2024-25 (Completed)
- - Successfully quantumized Free5GC and Aether SD-Core platforms
- - Integrated ML-KEM, ML-DSA, and QRNG across all network functions
- - Implemented hybrid PQC mode for migration support
- - Established coRAN LABS Public License framework
- - Launched community engagement with LFN, Anuket, and ONAP
- - Published research paper on arXiv (arXiv:2510.19982) documenting QORE architecture and implementation
- - Experimental validation with carrier-grade performance metrics

### 2025-26 Focus Areas
- Complete quantumization of additional open source 5G Core platforms (OAI, Open5GS, Magma)
- 3GPP Release 17+ compliance and standards alignment
- Performance optimization and GPU acceleration enhancements
- Production deployment support and operator trials
- Multi-vendor interoperability testing and certification
- Enhanced QRNG integration and edge deployment optimization
- SLH-DSA integration for CA root keys

### 2026 and Beyond
- Advanced quantum-safe features (network slicing, MEC security)
- Hardware acceleration partnerships for production-scale deployments
- Expanded ecosystem support and operator production pilots
- AI/ML integration for quantum threat detection and response
- Contribution to 3GPP Release 18+ quantum security specifications
- Quantum Key Distribution (QKD) integration research

---

## Contributing

QORE is an open research initiative. We welcome contributions from academia, industry, and the open source community.

### How to Contribute

1. **Code Contributions**: Implement PQC for additional NFs or platforms
2. **Testing**: Interoperability testing, performance benchmarking, security audits
3. **Documentation**: Deployment guides, API documentation, tutorials
4. **Research**: Algorithm optimization, protocol design, threat modeling

### Contribution Process

```bash
# Fork the repository
git clone https://github.com/coranlabs/QORE.git

# Create a feature branch
git checkout -b feature/pqc-implementation

# Make your changes and commit
git commit -m "Add ML-KEM support to component"

# Push and create a Pull Request
git push origin feature/pqc-implementation
```

**Contributor License Agreement**: By submitting a contribution, you agree to license your work under the coRAN LABS Public License v1.0.

---

## Publications and Media

### Research Papers
- **QORE: Quantum Secure 5G/B5G Core** - Rathi et al., 2024  
  arXiv:2510.19982 | [PDF](https://arxiv.org/pdf/2510.19982) | [Abstract](https://arxiv.org/abs/2510.19982) | [HTML](https://arxiv.org/html/2510.19982)
  
  *Published: October 22, 2025 | 23 pages | Subjects: Cryptography and Security, Distributed Computing, Networking*

### Whitepapers and Technical Documentation
- [coRAN Labs Whitepapers Repository](https://github.com/coranlabs/WhitePapers) - Comprehensive technical documentation, architecture guides, and research papers

### Videos and Demonstrations
1. [QORE: Implementing PQ-mTLS 1.3 in 5G Core](https://youtu.be/W5AgYsJQySw)
2. [5G QORE: Post-Quantum Cryptography in Action](https://youtu.be/rZCRh8JKKN8)
3. [QORE: Quantumized 5G Core Deployment](https://youtu.be/w1ac3SMiGmM)
4. [QORE: Post-Quantum Security for 5G Networks](https://www.youtube.com/watch?v=yiH2O24eUWk)

### Industry and Standards Body References
1. **Anuket TSC Discussion**: [Post-Quantum Cryptography in Cloud Native Telecom](https://lists.anuket.io/g/anuket-tsc/topic/111535604)
2. **ONAP TSC**: [Quantum Security Integration Proposal](https://lists.onap.org/g/onap-tsc/message/9611)
3. **LFN CNTI**: Input from ETSI on Quantum Security & Encryption - [PoC Slides](https://share.google/3zyXvHcW4i5jXn9zS)
4. **QORE Project Presentation**: [Slides and Technical Overview](https://share.google/IQKQa9fy4813DsL1M)

### Media Coverage
- [Quantum Zeitgeist: QORE Enables Transition to Post-Quantum Cryptography](https://quantumzeitgeist.com/quantum-networks-secure-b5g-core-qore-enables-transition-post-cryptography/)

---

## License
**QORE** is licensed under the **coRAN LABS Public License Version 1.0**.

### License Summary
- **Research and Academic Use**: Free, no restrictions
- **Commercial Use**: Requires FRAND (Fair, Reasonable, Non-Discriminatory) licensing
- **Patent Grant**: Royalty-free for research, negotiable for commercial deployment
- **Third-Party Components**: Original licenses apply (see [NOTICE](NOTICE) file)

**Full License**: [LICENSE](LICENSE)  
**Third-Party Notices**: [NOTICE](NOTICE)

### Commercial Licensing
For commercial deployment, product integration, or custom development:
- **Email**: contact@coranlabs.com
- **Partnership Inquiries**: contact@coranlabs.com

---

## Contact

### coRAN Labs
- **Website**: [www.coranlabs.com](https://www.coranlabs.com/)
- **Email**: contact@coranlabs.com
- **GitHub**: [github.com/coranlabs](https://github.com/coranlabs)

### Technical Support
- **Issue Tracker**: [GitHub Issues](https://github.com/coranlabs/QORE/issues)
- **Discussion Forum**: [GitHub Discussions](https://github.com/coranlabs/QORE/discussions)

### Research Collaboration
For research partnerships, academic collaboration, or joint publications:
- **Email**: contact@coranlabs.com
- **Cite Our Work**: See [Research Publication](#research-publication) section

---

<p align="center">
  <strong>Securing the Future of Telecommunications</strong><br>
  <em>QORE: Making Every 5G Core Quantum-Resistant</em>
</p>

<p align="center">
  Copyright © 2024 coRAN Labs and Contributors<br>
  Licensed under coRAN LABS Public License v1.0
</p>

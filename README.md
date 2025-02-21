

<table style="border-collapse: collapse; border: none;">
  <tr style="border-collapse: collapse; border: none;">
    <td style="border-collapse: collapse; border: none;">
      <a href="http://www.coranlabs.com/">
         <img src="./docs/coranlabs-logo.png" alt="" border=2 height=120 width=150>
         </img>
      </a>
    </td>
    <td style="border-collapse: collapse; border: none; vertical-align: center;">
      <b><h1>QORE: Quantum Secure Core</h1></b>
      <b><h2>Beyond 5G Core integrated with Post Quantum Cryptography</h2></b>
    </td>
  </tr>
</table>


## TABLE OF CONTENTS
1. [Introduction](#introduction)
2. [Need for QORE](#need-for-qore)
   - [Problem Statement](#problem-statement)
   - [Solution](#solution)
3. [Migration to Post-Quantum Core using QORE](#migration-to-post-quantum-core-using-qore)
4. [Current Scenario of QORE](#current-scenario-of-qore)
   - [SBI Protection](#sbi-protection)
   - [Subscriber's identity concealment: SUPI to SUCI conversion](#subscribers-identity-concealment-supi-to-suci-conversion)
   - [Post-Quantum Signatures and Certificates](#post-quantum-signatures-and-certificates)
   - [BackHaul Security](#backhaul-security)
   - [QRNG integration in Core Network](#qrng-integration-in-core-network)
   - [Larger Symmetric Keys](#larger-symmetric-keys)
5. [QORE Video](#qore-video)



## Introduction

**QORE: Quantum Secure Core** integrates **Post-Quantum Cryptography** (and **Quantum Random Number Generator (QRNG)**) into the Core network. Developed by [CoRan Labs](https://www.coranlabs.com/), QORE represents a significant advancement in ensuring robust security for core network against the impending threat of **Quantum Attacks**. By migrating classical cryptographic techniques used in the Core to Post-Quantum Cryptographic techniques, QORE offers enhanced security and reliability.


## Need for QORE?

#### Problem Statement:

The current Core, as defined by the 3GPP standard, currently relies on **classical cryptographic** techniques. However, these traditional encryption methods are increasingly vulnerable to quantum threats. **With the rise of quantum computers, classical cryptography can be easily broken.** Quantum computers have the capability to solve complex problems exponentially faster, allowing them to break traditional cryptographic algorithms. This renders current encryption methods insecure, exposing the classical Core to significant security risks.

#### Solution:
To secure the classical Core against these quantum threats, it is necessary to migrate to a **Post-Quantum Core**. This migration involves utilizing post-quantum cryptographic algorithms that are designed to be secure against the capabilities of quantum computers. Additionally, the generation of truly random numbers is crucial to ensure that cryptographic keys remain safe from quantum attacks.
QORE addresses these needs by integrating the following post-quantum techniques:
* `ML-KEM`: Module-Lattice-Based Key-Encapsulation Mechanism, ensures secure key exchange and protection against quantum attacks, utilizing lattice-based cryptography for strong security foundations.

* `ML-DSA`: Module-Lattice-Based Digital Signature Algorithm, a lattice-based digital signature scheme offering strong security guarantees against quantum computing threats.

* `PQ-mTLS1.3`: Post-Quantum Mutual Transport Layer Security to secure communication channels.

* `PQ-IPSec`: Post-Quantum IPSec for securing Internet Protocol Security communications.

* `PQ-DTLS1.3`: Post-Quantum Datagram Transport Layer Security for securing datagram communications.

* `PQ-PKI`: Post-Quantum Public Key Infrastructure with Post Quantum Certificate Authority.

* `QRNG seeds`: Utilizes Quantum Random Number Generators to produce truly random seeds, enhancing key security.

* `AES256`: Advanced Encryption Standard with 256-bit keys to ensure robust encryption.



## Migration to Post-Quantum Core using QORE

| **Feature**               | **Classical Core**                                  | **Qore (Post-Quantum Core)**          | **Status**  |
|---------------------------|-----------------------------------------------------|---------------------------------------|-------------|
| **SBI Communication**     | mTLS                                                | PQ-mTLS1.3(mTLS1.3 with PQ)                              | ✅ Done     |
| **SUPI to SUCI**          | ECIES                                               | PQ-IES(ML-KEM)                        | ✅ Done     |
|                           |                                                     | PQ-IES(Hybrid ML-KEM)                 | ✅ Done     |  
| **Digital Certificates**  | Classical Certificates                   | ML-DSA                                | ✅ Done     |
| **N2 User Data**          | DTLS                                                | PQ-DTLS1.3(DTLS1.3 with PQ)                              | ✅ Done     |
| **N2 User Data**          | IPSec                                               | PQ-IPSec (IKEv2 with PQ)                             | ✅ Done   |
| **N3 User Data**          | IPSec                                               | PQ-IPSec (IKEv2 with PQ)                             | ✅ Done   |
| **N4 User Data**          | IPSec                                               | PQ-IPSec (IKEv2 with PQ)                              | ✅ Done   |
| **PKI**         | Classical PKI/Private CA                | PQ-PKI/Private PQ-CA                                | 🟡Ongoing     |
| **Symmetric Key**         | AES128                                              | AES256*                                | ✅ Done     |
| **Random Number**         | PRNG                | QRNG*                                 | ✅ Done     |



> **Note**: **"*"** represents "suggestions"(not mandate) to improve the security. Symmetric Cryptography seems not to be affected much by quantum attacks(although people have mixed opinion on this) still 3GPP is moving towards 256-bits symmetric based cryptography. Similarly QRNG is used to generate truly random seeds in order to generate truly random cryptographic keys. AES256 & QRNG will provide higher level of security(even if they seem irrelevant to quantum attacks).

## Current Scenario of QORE

### SBI Protection

Communication between NFs currently uses mTLS which is not secure against quantum attacks. To address this vulnerability, PQ-mTLS 1.3 can be implemented for secure and quantum-safe communication across the network.

<img src="./docs/v3_pq_mtls.png" alt="Architecture Diagram" style="width: 700px;">

### Subscriber's identity concealment: SUPI to SUCI conversion

Securing SUPI(IMSI) using Post-quantum cryptography & QRNG

<img src="./docs/v2suci_pqc.png" alt="Architecture Diagram" style="width: 800px;">


QORE supports multiple Encryption Profiles. The solution incorporates a QRNG for key generation, it then uses ML-KEM for Key exchange mechansim and additionally, AES-128 is replaced with AES-256, further strengthening encryption.


<img src="./docs/v3_security_profile.png" alt="Architecture Diagram" style="width: 500px;">


## Post-Quantum Signatures and Certificates

The process of certificate verification is crucial for ensuring secure and trustworthy communication between NFs. It involves validating the authenticity and validity of post-quantum certificates. By leveraging post-quantum cryptographic algorithms for certificate generation and verification, the network protects itself against quantum-capable adversaries, ensuring the security of all communications.

<img src="./docs/signature_pq.png" alt="Architecture Diagram" style="width: 800px;">

### BackHaul Security

PQ-DTLS1.3 is  designed to secure N2 interface against quantum attacks, utilizing post-quantum cryptographic algorithms to ensure that data transmitted over SCTP  remains confidential and tamper-resistant.

<img src="./docs/DTLS-N2.png" alt="Architecture Diagram" style="width: 500px;">


### QRNG integration in Core Network

QRNGs leverage quantum processes to produce truly random numbers, ensuring high unpredictability and entropy. This level of randomness is crucial for cryptographic key generation, as it significantly enhances security by making it nearly impossible for attackers to predict or reproduce keys. In a post-quantum world, QRNGs become essential for maintaining robust security in telecom networks.

<img src="./docs/QRNG_int.png" alt="Architecture Diagram" style="width: 700px;">


### Larger Symmetric Keys

Upgrading from AES-128 to AES-256 allows telecom networks to substantially enhance their symmetric encryption,

<img src="./docs/symmetric.png" alt="Architecture Diagram" style="width: 500px;">

## QORE Video

1. [QORE: Implementing PQ-mTLS 1.3 in 5G/B5G Core](https://youtu.be/W5AgYsJQySw?si=Wct89_Gb7YqbeiDv)
2. [5G QORE: 5G Core with Post-Quantum Cryptography](https://youtu.be/rZCRh8JKKN8?si=UzkhHaLbAOznBGr1)
3. [QORE: Quatumized 5G Core](https://youtu.be/w1ac3SMiGmM?si=yOtkogSB39Neu1eb)

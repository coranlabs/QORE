# PQC-QORE - Post-Quantum Cryptography 5G Core Network

A complete 5G Core Network implementation with Post-Quantum Cryptography support.

This repository contains the integrated codebase for all CORAN 5G Network Functions with PQ-Security enhancements.

## Network Functions (NFs)

### Core NFs
- **CORAN_AMF** - Access and Mobility Management Function
- **CORAN_AUSF** - Authentication Server Function  
- **CORAN_CHF** - Charging Function
- **CORAN_NRF** - Network Repository Function
- **CORAN_NSSF** - Network Slice Selection Function
- **CORAN_PCF** - Policy Control Function
- **CORAN_SMF** - Session Management Function
- **CORAN_UDM** - Unified Data Management
- **CORAN_UDR** - Unified Data Repository

### Additional NFs
- **CORAN_CONSOLE** - Web Console
- **CORAN_SCP** - Service Communication Proxy
- **CORAN_NWDAF** - Network Data Analytics Function
- **CORAN_UPF_eBPF** - User Plane Function (eBPF)
- **CORAN_UPF_up3c** - User Plane Function (UP3C)
- **HEXA_UPF** - HEXA User Plane Function

## Libraries

### CORAN Libraries
- CORAN_LIB_APER - ASN.1 PER encoding/decoding
- CORAN_LIB_NAS - NAS protocol implementation
- CORAN_LIB_NGAP - NGAP protocol implementation  
- CORAN_LIB_OPENAPI - OpenAPI client/server
- CORAN_LIB_SCTP - SCTP wrapper
- CORAN_LIB_UTIL - Common utilities
- CORAN_LIB_PFCP - PFCP protocol
- CORAN_LIB_TLV - TLV encoding/decoding
- CORAN_LIB_FSM - Finite State Machine
- CORAN_LIB_LOGGER_CONF - Logger configuration
- CORAN_LIB_LOGGER_UTIL - Logger utilities  
- CORAN_LIB_PATH_UTIL - Path utilities

### External Dependencies
- **go-post-quantum** - Post-Quantum Cryptography library
- **httpwrapper** - HTTP/HTTPS wrapper with PQ support
- **jwt** - JWT token handling
- **sctp** - SCTP library
- **util_3gpp** - 3GPP utilities
- **webconsole** - Web console frontend
- **qore-sba-k8s** - Kubernetes deployment configurations

## Features

- ✅ Full 5G Core Network Functions
- ✅ Post-Quantum Cryptography (PQ-mTLS)
- ✅ PQ-Security Phase 1 & 2 implementation
- ✅ Container-ready with Docker support
- ✅ Kubernetes deployment support

## License

Copyright © 2024 CoranLabs

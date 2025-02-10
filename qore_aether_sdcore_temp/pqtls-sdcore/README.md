# Post Quantum TLS in the 5G Core:

- 5G Core has a Service based architecture (SBA), where the Network Functions are treated as their own microservices. These NFs are usually configured as Nodes in a Kubernetes Cluster, where communication b/w them is achieved via the TCP/IP stack. 3gpp mandates that this communication should be protected via Transport Layer Security. However, TLS itself is not quantum secure yet (though there are drafts that implement PQC, there's still very little support to it). For this purpose, we have added a Post Quantum solution to 5G core, where the certificates are replaced with the hybrid PQ X.509 certs and Key exchange by hybrid PQC. The algorithms used are:
  
  - **Ed448-Dilithium3** for certificates.
  - **X25519Kyber768** for key exchange.
  - Note that Kyber and Dilithium have been standardized to ML_KEM & ML_DSA respectively in the recent NIST standardization project.

Here are a few figures illustrating TLS in 5G SBA:
1. NF registers with the NRF:
   ![image](https://github.com/user-attachments/assets/a973194d-2608-41b2-afa0-f6f5ddc52b01)
   ![image](https://github.com/user-attachments/assets/cf8708cd-3593-4502-8bb5-a51875c2cef1)
2. Hybrid PQ TLS rough overview:
   ![image](https://github.com/user-attachments/assets/36a13bf9-05bb-402f-920b-b4fcc8b44111)
3. [A concise overview of TLS key schedule and HKDF](https://www.figma.com/board/3MRPFPUQMAVZjIyLuE8TPq/TLS-Key-Schedule?node-id=0-1)
4. Generation of certificates and CA:
   ![image](https://github.com/user-attachments/assets/e7aad688-b310-4823-b23b-d6ce79ee4205)
5. 5G SBA with Hybrid PQ TLS:
   ![image](https://github.com/user-attachments/assets/9da9fa12-b390-4495-a1c6-bea109cc1290)
 
## Pre-requisites:
```
- Golang (v >= 1.21)
- K8S
- Aether-in-a-box
```

## Overview
In this project, I have added additional functionalities on top of the Network functions and the utilities provided by [OMEC-project](https://github.com/omec-project), which enable Hybrid Post Quantum TLS support. This is done by leveraging [Cloudflare's Golang with experimental patches](https://github.com/cloudflare/go), which uses Cloudflare's CIRCL library to add PQ functionality in the golang's `TLS & X509` standard libraries. For this, we are required to modify the Dockerfiles to change Golang's `environment variables` like `GOROOT`. Further, one also needs to copy the certificates from the docker images to the actual binary executable of the concerned Network Functions.
The certificates can be generated via golang's inbuilt utility: `crypto/tls/generate_cert.go`, pass the common names, CA or not flag, and key algorithms, along with the host, i.e. your IP (could be localhost for instance).

> guide:
>
> **GOROOT** -> places for go to look for the standard golang libraries. By default: `/usr/local/go`
>
> **GOBIN** -> Go executable path
> 
> **GOPATH** -> places where the non standard (installed) go packages can be found. By default: `~/go`
> 

1. In 5G Core, a Service Based Architecture (SBA) is followed, each NF is its own microservice, usually linked via NRF or a Service Communication Proxy (SCP).
2. At the time of initialization, every NF registers itself to the NRF, whose **uri** is already stored in NF's initial configuration. This is an http(s) request, utilised by Restful/Open APIs - **NF Register Request**. During this time, access tokens are provided to the NFs by the NRF for authorization & authentication later.
3. Whenever, an _NF_1_ wants to communicate via _NF_2_, it sends an **NF Discovery Request : NF_2** message to the NRF. The NRF is responsible for first verifying the authenticity of _NF_1_ and *NF_2*, which is done via the ***JWT*** access tokens of both the network functions. On successful verification, the request from NF_1 is routed to NF_2. The data exchange b/w the 2 NFs is then carried out via HTTPs.

#### This project supports both classical and post quantum methods, thus highly interoperable with the existing systems, making it an easier migration to quantum safety.

## Changes in Aether 5GC:
1. Modified ``sd-core-5g-values.yaml`` to add **HTTPS** to the SBI schemes of each NF. Also add the path to the certificate.
   FOR example:
   - ![image](https://github.com/user-attachments/assets/af8770c1-5e26-4448-885c-f86e8bb394e3)
   - [Config File that we've made use of](https://github.com/lakshya-chopra/PQ-TLS-IN-5GC/blob/main/sd-core-5g-values.yaml)
2. The helm charts have been rewritten to change `nrf URI` protocol from **http** to **https** in each _NF-NRF_ (read: NF to NRF and NRF to NF) communication API (for example: **NF register**).
   
   - ![image](https://github.com/user-attachments/assets/fdde25ee-82be-4697-8399-3497e49dad1c)
     
3. The Golang standard library has been changed to [Go with experimental patches](https://github.com/lakshya-chopra/go), a custom fork over cloudflare's go which enables PQ KEMs & Signature algorithms by default in the peer's supported groups and supported signature algorithms respectively (sent in ClientHello).
4. No changes to OpenAPIs have been made, because they are responsible for HTTP1.1/2 *GET* & *POST* requests, not for carrying the actual application data.
5. No changes have been made to SSL key logs, because they don't directly reflect what key exchanges have been used. Here's for example, the SSL key logs of TLS v/s PQ TLS:
   - TLS:
    ![image](https://github.com/user-attachments/assets/de9323ad-472b-41f1-b58f-e67b30c22748)
   - PQ TLS:
     ![image](https://github.com/user-attachments/assets/eeb36c7e-4037-4964-a1aa-bb1996b742b8)
6. Added utilities to read PQ & classical certificates and print some of their details, such as Public Key, Organization, Serial Number, Signature Algorithm, and the CA.

## Changes in OMEC's utilities:
- For creating a new HTTP(s) SBI server, all the Network functions make use of the `http2_util` - a small and concise utility to set up new servers and configure it for TLS. This utility has been modified to include PQ signature schemes, and log the connection details. It adds **X25519Kyber768Draft00** to the curve preferences as well. Note that this is only possible with the modified go standard library mentioned previously.
- Here's a small gist of the code (ignore the linter error, this is due to the use of non standard go lib):
  ![image](https://github.com/user-attachments/assets/9d5fcfa8-2c7b-4229-b15c-f90c6c5872ef)
- [Link to the full code](https://github.com/lakshya-chopra/http2_util)



## Logs showcasing PQ-TLS:
1. Certificates:
   
   ![image](https://github.com/user-attachments/assets/9938878a-0150-4bc3-9ac2-d05514951bea)
   ![image](https://github.com/user-attachments/assets/2cdb76fe-aae7-45de-9b71-1120ee941f38)
   ![image](https://github.com/user-attachments/assets/8e078c82-42b5-433f-a82e-1a8a03280bce)


3. NF register/connection request:
   
   ![image](https://github.com/user-attachments/assets/36726f00-2f2a-45d6-96b3-27e115fd4301)
   Here, the signature algorithms ID: **65121** & **65122** indicate *Ed448-Dilithium2 and Ed448-Dilithium3* respectively.
   Furthermore, the curve IDs: **25497** and **65074** indicate *X25519Kyber768/512*.

   In the newer version, I have made this even more verbose:
   ![image](https://github.com/user-attachments/assets/1a13ec3c-ecea-4ffa-a4d4-f58b97dbe9b9)

3. TCP dump which demonstrates the KEX as **x25519Kyber768** (+ P256_Kyber768) and **Ed448-Dilithium3** as the sig. algorithm:
   - KEX: ![image](https://github.com/user-attachments/assets/45334f2d-e032-45b9-bece-a387021b7c96)
    ![image](https://github.com/user-attachments/assets/83aa9543-d324-4f9e-b7d7-4d21df40f531)

   - Signature algorithm:
     
    ![image](https://github.com/user-attachments/assets/ced54653-e47d-4167-96db-38f5064b1c68)
   
   This pcap file can be generated via the command:
    ```
    sudo tcpdump -i any host <CLUSTER_IP_OF_NF> -w capture_nf.pcap
    ```
    View the NF's cluster IP via: `kubectl get svc -n aether` 

4. PDU session: UE and RAN stack made using the UERAN simulator were connected to the modified 5G core, with appropriate IMSI, PLMN and SUPI Protection schemes. A PDU session was successfully made:
   - AMF: ![image](https://github.com/user-attachments/assets/13150cb1-905c-4481-865d-d355a4cb09d5)
   - USIM: ![image](https://github.com/user-attachments/assets/0dc2956d-c34c-45ff-ad14-d87420e06089)
 




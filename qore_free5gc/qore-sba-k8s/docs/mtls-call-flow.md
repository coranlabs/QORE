## Steps to Enable PQ-mTLS Support in QORE (Kubernetes Environment)

1. **Generate PQ Certificates for Each Network Function (NF):** For every NF in the 5G Core Network, a post-quantum (PQ) certificate and its corresponding private key need to be generated using [custom-go lib](go-cert-gen.md).

2. **Organize Certificates and Keys:** Store the generated certificates and private keys under the `/cert/` directory, organized by the respective NF names for clarity and ease of access.

3. **Mount Certificates Using Persistent Volumes:** In the Kubernetes environment, mount the `/cert/` directory to a **Persistent Volume (PV)** to ensure certificates persist across pod restarts or rescheduling.

4. **Access via Persistent Volume Claims (PVCs):** Each NF should access its own PQ certificate, private key, and the **NRF’s PQ certificate** from the PV using a **Persistent Volume Claim (PVC)**.

> [!NOTE] 
> In this setup, the **NRF acts as the Certificate Authority (CA)**. Since an authorized external CA is not used, NRF issues and manages certificates. Therefore, the NRF’s certificate is made available to all NFs in the network for verification purposes.

5. **Establish Secure SBI Communication:** When one NF communicates with another over the **Service-Based Interface (SBI)**, it presents its own PQ certificate as part of the **PQ-mTLS handshake**.

6. **Certificate Verification:** The receiving NF validates the certificate presented by the peer NF using the **NRF’s public key**, which is embedded within the NRF’s certificate.



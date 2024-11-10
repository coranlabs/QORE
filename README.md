## A Simple Golang wrapper over OpenSSL's DTLS over SCTP/UDP
- Uses CGO
- Compatible with non-blocking sockets
- Can be used with both UDP & SCTP sockets
- Uses Memory BIOs for writing application data to the sockets - Encrypted data is written to the memory BIO & then sent to the peer, who then decrypts it using the SSL secrets established during the handshake.
  - Makes use of SSL's r/w BIOs.
  - Flow:
    ![image](https://github.com/user-attachments/assets/d0227bfa-f96f-40e9-bc7f-b4e3d8d6c3a1)
  

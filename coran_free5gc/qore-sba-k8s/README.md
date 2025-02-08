# QORE-SBA-K8s

Deployment of the QORE 5G Core in a Kubernetes (K8s) environment.

## Hardware Requirements

- **OS:** Ubuntu 22.04 LTS  
- **CPU:** >= 8 cores
- **RAM:** >= 8GB  
- **Network Interfaces:** 1 (minimum)


## QORE 5G Core Deployment

### Clone the Repository

```bash
git clone https://github.com/coranlabs/qore-sba-k8s.git
cd qore-sba-k8s
```

### Install Dependencies

```bash
sudo apt update
sudo apt install -y make net-tools
```

### Prepare Node and Install 5G Core

```bash
make node-prep
make install_5g_core
```

---

## UERANSIM Deployment

### Clone UERANSIM Repository

```bash
cd ~
git clone https://github.com/aligungr/UERANSIM
```

### Install Build Dependencies

```bash
sudo apt update
sudo apt install -y make gcc g++ libsctp-dev lksctp-tools iproute2 cmake
```

### Build UERANSIM

```bash
cd ~/UERANSIM
make
```

### Configuration Files

Add the following config files to `~/UERANSIM/config/`:

- [gNB Configuration](./ueransim-config/coran-gnb.yaml)
- [UE Configuration](./ueransim-config/coran-ue.yaml)

### Update IP Addresses

- In `coran-gnb.yaml`:
  - Update `linkIp`, `ngapIp`, and `gtpIp` with the IP address of UERANSIM gNB.
  - Update `amfConfigs` with the AMF IP address.

- In `coran-ue.yaml`:
  - Update `gnbSearchList` with the IP address of UERANSIM gNB.

---

## OAI Deployment:
OAI rel: [v2.2.0](https://gitlab.eurecom.fr/oai/openairinterface5g/-/releases/v2.2.0)

> [NOTE]: OAI v2.2.0 might give error regarding DRB and (or) QoS, since it only supports 1 DRB & 1 QoS per DRB. To bypass this, navigate to `openair2/COMMON/e1ap_messages_types.h` and change the value of the macro `E1AP_MAX_NUM_NGRAN_DRB` from 4 to 16. Then go to the file `openair2/LAYER2/nr_pdcp/cucp_cuup_handler.c` and remove the AssertFatal regarding NumDRBs.

> ![image](https://github.com/user-attachments/assets/c56628f3-0beb-4fd3-9318-a47ad2c5a02b)

Build
```sh
cd openairinterface5g/cmake_targets/
sudo ./build_oai -w USRP --ninja --nrUE --gNB --build-lib "nrscope" 
```

### Update AMF, NGU & NG_AMF IP addresses:
 - In `targets/PROJECTS/GENERIC-NR-5GC/CONF/gnb.sa.band78.fr1.106PRB.usrpb210.conf` file (we are using USRP B210):
   
   - Update `amf_ip_address` with the external IP of the AMF
   - Update `GNB_IPV4_ADDRESS_FOR_NG_AMF` & `GNB_IPV4_ADDRESS_FOR_NGU` to your RAN host's machine IP.
 
 - If testing via CU-DU split, update the corresponding `cu_gnb.conf` & `du_gnb.conf`

## Add Subscriber in Dashboard

1. Open the dashboard in your browser:

   ```
   http://<amf-external-ip>:30500
   ```

> [!NOTE]
> Get the AMF external IP using: `kubectl get svc -o wide`


2. In the **Subscribers** section, create a new user:
   - Update **IMSI** and **PLMN ID**.
   - Delete existing **DNN Configurations** and create a new one named `coran`.

---

## End-to-End Network Testing

### Terminal 1: View AMF Logs

```bash
kubectl logs <amf-pod> -n coran -f
```

### Testing via UERANSIM:
#### Terminal 2: Start gNB

```bash
cd ~/UERANSIM/config/
sudo ../build/nr-gnb -c coran-gnb.yaml
```

#### Terminal 3: Start UE

```bash
cd ~/UERANSIM/config/
sudo ../build/nr-ue -c coran-ue.yaml
```

### Testing via OAI:

#### Monolithic:
```sh
cd openairinterface5g/cmake_targets/ran_build/build
sudo ./nr-softmodem -O ../../../targets/PROJECTS/GENERIC-NR-5GC/CONF/gnb.sa.band78.fr1.106PRB.usrpb210.conf --sa -E --continuous-tx
```

#### CU-DU split:

CU:
```sh
sudo ./nr-softmodem -O ../../../targets/PROJECTS/GENERIC-NR-5GC/CONF/cu_gnb.conf --sa --continuous-tx --thread-pool 8,9,10
```
DU:
```sh
sudo ./nr-softmodem -O ../../../targets/PROJECTS/GENERIC-NR-5GC/CONF/du_gnb.conf --sa --continuous-tx --thread-pool 11,12,13,14,15
```
If the DU fails due to radio frequency mismatch errors, add an `-E` flag to the DU's cmd.


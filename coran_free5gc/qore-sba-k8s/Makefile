# Copyright 2018-present Open Networking Foundation
#
# SPDX-License-Identifier: Apache-2.0

SHELL		:= /bin/bash
MAKEDIR		:= $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
BUILD		?= $(MAKEDIR)/build
M           	?= $(BUILD)/milestones
SCRIPTDIR	:= $(MAKEDIR)/scripts
VENV		?= $(BUILD)/venv/aiab

KUBECTL_VERSION   ?= v1.25.15
LPP_VERSION       ?= v0.0.30
K8S_INSTALL ?= rke2
RKE2_K8S_VERSION  ?= v1.25.15+rke2r2

# K8S_VERSION       ?= v1.21.6
CTR_CMD     := sudo /var/lib/rancher/rke2/bin/ctr --address /run/k3s/containerd/containerd.sock --namespace k8s.io


# 5G CORE

NAMESPACE = coran
CORAN_CHART = $(MAKEDIR)/charts/coran-core
CORAN_GATEWAY = $(MAKEDIR)/charts/coran-gateway
CORAN_UPF = $(MAKEDIR)/charts/coran-upf
UERANSIM_CHART = $(MAKEDIR)/charts/ueransim
VALUES_FILE = $(MAKEDIR)/values.yaml
AMF_EXTERNAL_IP ?= $(NODE_IP)


NODE_IP ?= $(shell ip route get 8.8.8.8 | grep -oP 'src \K\S+')  # TODO: add warnings for each dynamic values
NODE_IP_SUBNET ?= $(shell echo $(NODE_IP) | cut -d. -f1-3) # TODO: add warnings for each dynamic values
DATA_IFACE ?= $(shell ifconfig | grep -B1 "${NODE_IP}" | head -n1 | awk '{print $$1}' | sed 's/://') # TODO: add warnings for each dynamic values
WEBUI_ADDR ?= $(shell echo $(NODE_IP) | xargs):30500
CERTS_PWD ?= $(shell echo $(PWD)/cert)

.PHONY: show
show:
	@echo "NODE_IP: $(NODE_IP)"
	@echo "NODE_IP_SUBNET: $(NODE_IP_SUBNET)"
	@echo "DATA_IFACE: $(DATA_IFACE)"
	@echo "WEBUI_ADDR: $(WEBUI_ADDR)"
	@echo "CERTS_PWD: $(CERTS_PWD)"

.PHONY: url
url:
	@echo "http://$(WEBUI_ADDR)/api/subscriber/imsi-$(IMSI)/$(PLMN)"


ifndef NODE_IP
$(error NODE_IP is not set)
endif

.PHONY: node-prep clean setup_env install_5g_core install_5g_gateway install_5g_upf \
	uninstall_5g_gateway uninstall_5g_upf uninstall_5g_core

$(M):
	mkdir -p $(M)

$(M)/setup: | $(M)
	@if [ "$(DATA_IFACE)" = "eth0" ]; then \
		echo "Data interface is not set. Using default data interface: $(DATA_IFACE)"; \
	else \
		echo "Data interface is set to: $(DATA_IFACE)"; \
	fi
	sudo apt update; sudo apt install -y software-properties-common python3 python3-pip python3-venv jq httpie ipvsadm net-tools
	touch $@

$(M)/k8s-ready: | $(M)/setup
	sudo mkdir -p /etc/rancher/rke2/
	echo "cni: multus,calico" >> config.yaml
	echo "cluster-cidr: 192.168.84.0/24" >> config.yaml
	echo "service-cidr: 192.168.85.0/24" >> config.yaml
	echo "kubelet-arg:" >> config.yaml
	echo "- --allowed-unsafe-sysctls="net.*"" >> config.yaml
	echo "- --node-ip="$(NODE_IP)"" >> config.yaml
	sudo mv config.yaml /etc/rancher/rke2/
	curl -sfL https://get.rke2.io | sudo INSTALL_RKE2_VERSION=$(RKE2_K8S_VERSION) sh -
	sudo systemctl enable rke2-server.service
	sudo systemctl start rke2-server.service
	sudo /var/lib/rancher/rke2/bin/kubectl --kubeconfig /etc/rancher/rke2/rke2.yaml wait nodes --for=condition=Ready --all --timeout=300s
	mkdir -p $(HOME)/.kube
	sudo cp /etc/rancher/rke2/rke2.yaml $(HOME)/.kube/config
	sudo chown -R $(shell id -u):$(shell id -g) $(HOME)/.kube
	sudo /var/lib/rancher/rke2/bin/kubectl --kubeconfig /etc/rancher/rke2/rke2.yaml wait deployment -n kube-system --for=condition=available --all --timeout=300s
	@$(eval STORAGE_CLASS := $(shell /var/lib/rancher/rke2/bin/kubectl --kubeconfig /etc/rancher/rke2/rke2.yaml get storageclass -o name))
	@echo "STORAGE_CLASS: ${STORAGE_CLASS}"
	if [ "$(STORAGE_CLASS)" == "" ]; then \
		sudo /var/lib/rancher/rke2/bin/kubectl --kubeconfig /etc/rancher/rke2/rke2.yaml apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/$(LPP_VERSION)/deploy/local-path-storage.yaml --wait=true; \
		sudo /var/lib/rancher/rke2/bin/kubectl --kubeconfig /etc/rancher/rke2/rke2.yaml patch storageclass local-path -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'; \
	fi
	curl -LO "https://dl.k8s.io/release/$(KUBECTL_VERSION)/bin/linux/amd64/kubectl"
	sudo chmod +x kubectl
	sudo mv kubectl /usr/local/bin/
	kubectl version --client
	mkdir -p $(HOME)/.kube
	sudo cp /etc/rancher/rke2/rke2.yaml $(HOME)/.kube/config
	sudo chown -R $(shell id -u):$(shell id -g) $(HOME)/.kube
	touch $@

$(M)/helm-ready: | $(M)/k8s-ready
	curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
	chmod 700 get_helm.sh
	sudo DESIRED_VERSION=$(HELM_VERSION) ./get_helm.sh
	touch $@

node-prep: | $(M)/helm-ready $(M)/k8s-ready


clean: |
	sudo /usr/local/bin/rke2-uninstall.sh || true
	sudo rm -rf /usr/local/bin/kubectl
	rm -rf $(M)
	touch $@

setup_env:
	@echo "Setting up the environment..."
	chmod +x ./setup_env.sh
	sudo ./setup_env.sh

# install_5g_core: install_5g_cp install_5g_gateway install_5g_upf  // TODO: remove 

install_5g_core: node-prep setup_env setup_network
	@echo "Installing Coran CP..."
	@if helm list -n $(NAMESPACE) | grep -q 'coran'; then \
		echo "coran CP is already installed."; \
	else \
		NODE_IP=${NODE_IP} DATA_IFACE=${DATA_IFACE} NODE_IP_SUBNET=${NODE_IP_SUBNET} AMF_EXTERNAL_IP=${AMF_EXTERNAL_IP} CERTS_PWD=${CERTS_PWD} envsubst < $(VALUES_FILE) > /tmp/values-substituted.yaml; \
		helm upgrade --create-namespace --install --wait \
			--namespace coran \
			-f /tmp/values-substituted.yaml \
			coran \
			$(CORAN_CHART); \
		$(MAKE) add_subscriber; \
	fi
	

# install_5g_gateway:   // TODO: remove 
# 	@echo "Installing Coran Gateway..."
# 		helm upgrade --create-namespace --install --wait \
# 			--namespace coran \
# 			gateway \
# 			$(CORAN_GATEWAY); \

setup_network: 
	@if ! ip link show access &>/dev/null; then \
		echo "Creating macvlan interface..."; \
		sudo ip link add access link $(DATA_IFACE) type macvlan mode bridge; \
	fi
	@if ! ip addr show access | grep -q '10.100.50.238'; then \
		echo "Assigning IP address to access interface..."; \
		sudo ip addr add 10.100.50.238/29 dev access; \
	fi
	sudo ip link set access up


# uninstall_5g_gateway: uninstall_5g_upf  // TODO: remove 
# 	@echo "Uninstalling Coran Gateway..."
# 	@if helm list -n coran | grep -q 'gateway'; then \
# 		helm uninstall gateway --namespace coran; \
# 	else \
# 		echo "Coran Gateway is not installed."; \
# 	fi

# uninstall_5g_upf:  // TODO: remove 
# 	@echo "Uninstalling Coran UPF..."
# 	@if helm list -n coran | grep -q 'upf'; then \
# 		helm uninstall upf --namespace coran; \
# 	else \
# 		echo "Coran UPF is not installed."; \
# 	fi

# uninstall_5g_core: uninstall_5g_upf uninstall_5g_gateway uninstall_5g_cp // TODO: remove 

uninstall_5g_core:
	helm delete -n coran $$(helm -n coran ls -qa) || true
	kubectl delete namespace coran || true
	@echo ""
	@echo "Wait for all pods to terminate..."
	kubectl wait -n coran --for=delete --all=true -l app!=ue pod --timeout=180s || true

reset_5g_core: uninstall_5g_core install_5g_core

install_ueransim:
	@echo "Installing UERANSIM..."
	@if helm list -n $(NAMESPACE) | grep -q 'ueransim'; then \
		echo "UERANSIM is already installed."; \
	else \
		helm install -n $(NAMESPACE) ueransim $(UERANSIM_CHART); \
	fi

# Target to uninstall UERANSIM
uninstall_ueransim:
	@echo "Uninstalling UERANSIM..."
	@if helm list -n $(NAMESPACE) | grep -q 'ueransim'; then \
		helm uninstall ueransim -n $(NAMESPACE); \
	else \
		echo "UERANSIM is not installed."; \
	fi

# Define the default IMSI, but allow it to be overridden when running the target
IMSI ?= 001010000000001
PLMN ?= 00101

# Define the curl command with the flags
# •	-X GET: This makes a GET request.
# •	-X POST: This makes a POST request.
# •	-s: This makes curl silent (no progress or error messages).
# •	-o /dev/null: This discards the body of the response.
# •	-w "IMSI $(IMSI) added successfully\n": This tells curl to display the custom success message only, after the command is executed.

# Command to get subscriber list
get_subscribers:
	@echo "Fetching the list of subscribers..."
	@curl -s -X GET http://$(WEBUI_ADDR)/api/subscriber -H "Token: admin" | jq

# Command to get a particular subscriber detail
get_subscriber:
	@if [ -z "$(IMSI)" ]; then \
		echo "Error: Please provide IMSI using 'IMSI=your-imsi-value'"; \
	else \
		echo "Fetching details for IMSI $(IMSI)..."; \
		curl -s -X GET http://$(WEBUI_ADDR)/api/subscriber/imsi-$(IMSI)/$(PLMN) -H "Token: admin" | jq; \
	fi

# Command to add a subscriber
add_subscriber:
	@if [ -z "$(IMSI)" ]; then \
		echo "Error: Please provide IMSI using 'IMSI=your-imsi-value'"; \
	else \
		echo "Adding subscriber with IMSI $(IMSI)..."; \
		curl -s -o /dev/null -w "IMSI $(IMSI) added successfully\n" -X POST http://$(WEBUI_ADDR)/api/subscriber/imsi-$(IMSI)/$(PLMN) -H "Token: admin" -d @subscriber.json; \
	fi

# Command to update a subscriber
update_subscriber:
	@if [ -z "$(IMSI)" ]; then \
		echo "Error: Please provide IMSI using 'IMSI=your-imsi-value'"; \
	else \
		echo "Updating subscriber with IMSI $(IMSI)..."; \
		curl -s -X PUT http://$(WEBUI_ADDR)/api/subscriber/imsi-$(IMSI)/$(PLMN) -H "Token: admin" -d @subscriber.json -o /dev/null -w "IMSI $(IMSI) updated successfully\n"; \
	fi

# Command to delete a subscriber
delete_subscriber:
	@if [ -z "$(IMSI)" ]; then \
		echo "Error: Please provide IMSI using 'IMSI=your-imsi-value'"; \
	else \
		echo "Deleting subscriber with IMSI $(IMSI)..."; \
		curl -s -X DELETE http://$(WEBUI_ADDR)/api/subscriber/imsi-$(IMSI)/$(PLMN) -H "Token: admin" -o /dev/null -w "IMSI $(IMSI) deleted successfully\n"; \
	fi

# Command to check PDP active subscriber
check_pdp_active_subscriber:
	@echo "Fetching PDP active subscribers..."
	@curl -s -X GET http://$(WEBUI_ADDR)/api/registered-ue-context -H "Token: admin" | jq
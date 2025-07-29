#!/bin/bash

# Check if the argument is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <t2_value>"
    exit 1
fi

test=$1

# Build the Docker image
docker build -t amf:v1.1.0-$test .

# Tag the Docker image
docker tag amf:v1.1.0-$test khushichhillar/shadow_amf:v1.1.0-$test

# Push the Docker image
docker push khushichhillar/shadow_amf:v1.1.0-$test
cd ~/aether-in-a-box
#sudo sed -i "s|khushichhillar/shadow_amf:v1.1.0-t1|khushichhillar/shadow_amf:v1.1.0-$test|g" sd-core-5g-values.yaml
sudo make omec-clean
sudo ENABLE_GNBSIM=false DATA_IFACE=eth0 CHARTS=local make 5g-core -B

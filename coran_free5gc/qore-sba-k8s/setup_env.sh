#!/bin/bash

# Function to install environment and set up components
#!/bin/bash

# Function to install environment and set up components
setup_environment() {
    # Install necessary packages
    echo "Installing necessary packages (make, gcc)..."
    sudo apt update
    sudo apt install -y make gcc

    # Clone and build the up3c module
    echo "Cloning and building up3c module..."
    if [ ! -d "up3c" ]; then
        git clone https://github.com/coranlabs/up3c.git
    fi

    cd up3c || exit
    git checkout v0.8.5
    make clean
    make
    sudo make install
    
    if [ -d "up3c" ]; then
        echo "Removing up3c directory..."
        sudo rm -rf up3c
    fi
    # Go back to the parent directory
    cd ..
}

# Function to check if up3c module is already installed
check_up3c_installed() {
    if lsmod | grep -q 'up3c'; then
        echo "up3c module is already installed."
        return 0
    else
        echo "up3c module is not installed. Proceeding with installation..."
        return 1
    fi
}

# Function to clean up the environment
clean_setup_environment() {
    echo "Cleaning up the setup environment..."

    # Remove the up3c directory
    if [ -d "up3c" ]; then
        echo "Uninstalling up3c directory..."
        cd up3c || exit
        sudo make clean
        sudo make uninstall
    else
        git clone https://github.com/coranlabs/up3c.git
        echo "Uninstalling up3c directory..."
        cd up3c || exit
        sudo make clean
        sudo make uninstall
    fi
}

# Main logic
if [ "$1" == "clean" ]; then
    clean_setup_environment
else
    # Check if up3c module is already installed
    if ! check_up3c_installed; then
        setup_environment
    fi
fi


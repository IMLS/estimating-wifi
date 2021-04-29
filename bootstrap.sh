#!/bin/bash

# TESTING ENV VARS
# NOKEYREAD - set this to 1 to prevent the key from being read in.
# NOLOCKDOWN - prevents the pi from hardening and locking down. For testing.
# NOREBOOT - prevents reboot at end of bootstrap.sh
# 
# Usage:
# NOKEYREAD=1 NOLOCKDOWN=1 bash <(curl -s ...)

# We assume this script will be curl'd in.
# We assume it will be curl'd in and sudo will prompt for a password.

# It is for configuring a Raspberry Pi to be part of a pilot data collection effort.
# That pilot is being run in partnership between 10x/18F/IMLS.
# If you are not one of the people taking part in that pilot, then this
# software will *not* be useful to you.
# It will do things to your Raspberry Pi.
# Things you might not want.
# You have been warned.
# Here be krackens.

# CRITICAL GLOBALS
REPOS_ROOT="https://github.com/jadudm"
PLAYBOOK_REPOS="imls-client-pi-playbook"
PLAYBOOK_URL="${REPOS_ROOT}/${PLAYBOOK_REPOS}"
PLAYBOOK_WORKING_DIR="/opt/imls"
INITIAL_CONFIGURATION_BINARY_URL="https://github.com/jadudm/input-initial-configuration/releases/download/v0.0.3/input-initial-configuration"
SESSION_COUNTER_CONFIG_DIR="/etc/session-counter"
RALINK_DIR="/tmp/ralink"
RALNK_BINARY="https://github.com/jadudm/find-ralink/releases/download/v0.0.7/find-ralink"

# A GLOBAL CATCH
# If something goes wrong, set this to 1.
# If the _err function is ever used, it sets this automatically.
SOMETHING_WENT_WRONG=0

# PURPOSE
# Creates a temporary logfile in a way that lets the OS
# decide where it should go.
create_logfile () {
    export SETUP_LOGFILE=$(mktemp -t "setup-log-XXX")
}

mangle_console () {
    # https://serverfault.com/questions/103501/how-can-i-fully-log-all-bash-scripts-actions
    # Save all the pipes.
    # 3 is Stdout. 4 is stderr.
    exec 3>&1 4>&2
    # Restore some.
    trap 'exec 2>&4 1>&3' 0 1 2 3
    exec 1>> /dev/null 2>&1
}

# PURPOSE
# Sets up redirects so that STDOUT and STDERR make their way to
# a temporary logfile.
setup_logging () {
    mangle_console
    # Redirect stdout/stderr to a logfile.
    exec 1>> "${SETUP_LOGFILE}" 2>&1
    _status "Logfile started. It can be accessed for debugging purposes."
    _variable "SETUP_LOGFILE"
}


# COLORS!
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
PURPLE='\033[0;35m'
# No color
NC='\033[0m'

_msg () {
    TAG="$1"
    COLOR="$2"
    MSG="$3"
    printf "[${TAG}] ${MSG}\n" >&1
    printf "[${COLOR}${TAG}${NC}] ${MSG}\n" >&3
}

_status () {
    MSG="$1"
    _msg "STATUS" "${GREEN}" "${MSG}"
}

_debug () {
    MSG="$1"
    _msg "DEBUG" "${YELLOW}" "${MSG}"
}

_err () {
    SOMETHING_WENT_WRONG=1
    MSG="$1"
    _msg "ERROR" "${RED}" "${MSG}"
}

_variable () {
    VAR="$1"
    _msg "$VAR" "${PURPLE}" "${!VAR}"
}

####################################
# CHECKS
# These are helper functions for checking if things exist,
# etc. Used a lot, clarifies the code.

# https://stackoverflow.com/questions/592620/how-can-i-check-if-a-program-exists-from-a-bash-script
command_exists () {
    type "$1" &> /dev/null ;
}

command_does_not_exist () {
    if command_exists "$1"; then
        return 1
    else
        return 0
    fi
}

# PURPOSE
# Restores the file descriptors after capturing/redirecting.
restore_console () {
    # https://stackoverflow.com/questions/21106465/restoring-stdout-and-stderr-to-default-value
    # Reconnect stdout and close the third filedescriptor.
    exec 1>&4 4>&-
    # Reconnect stderr
    exec 1>&3 3>&-
}

check_for_usb_wifi () {
    mangle_console
    # This will need lshw
    sudo apt update
    sudo apt install -y lshw

    rm -rf ${RALINK_DIR}
    mkdir -p ${RALINK_DIR}
    pushd ${RALINK_DIR}
        rm -f find-ralink
        curl -L -s -o find-ralink ${RALNK_BINARY}
        chmod 755 find-ralink
        
        if [[ "$(./find-ralink --exists)" =~ "false" ]]; then
            restore_console
            echo "********************* PANIC OH NOES! *********************"
            echo "We think you did not plug in the USB wifi adapter!"
            echo "Please do the following:"
            echo ""
            echo " 1. Plug in the USB wifi adapter."
            echo " 2. Push the up arrow. (This brings back the bash command.)"
            echo " 3. Press enter."
            echo ""
            echo "This will start the setup process again."
            echo "********************* PANIC OH NOES! *********************"
            echo ""
            exit
        fi
    popd
    restore_console
}

read_initial_configuration () {
    # Create a place for it to go
    sudo mkdir -p $SESSION_COUNTER_CONFIG_DIR
    # Fetch the binary.
    pushd /tmp
        # 20210427 MCJ Again, in dev/testing conditions, wipe things out.
        rm -f iic
        curl -L -s -o iic ${INITIAL_CONFIGURATION_BINARY_URL}
        chmod 755 iic
        sudo ./iic --fcfs-seq --tag --word-pairs --write
    popd
}

bootstrap_ansible () {
    _status "Bootstrapping Ansible"
    _status "Updating sources."
    echo "deb http://ppa.launchpad.net/ansible/ansible/ubuntu trusty main" | sudo tee -a /etc/apt/sources.list
    _status "Installing dirmngr."
    sudo apt-get install dirmngr -y
    _status "Adding local keys."
    sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 93C4A3FD7BB9C367
    _status "Doing an update. This may take a moment. Be patient."
    sudo apt-get update
    _status "Installing the most recent ansible."
    sudo apt-get install -y ansible
}

install_prerequisites () {
    sudo apt-get install -y git
}

# PURPOSE
# This clones and runs the playbook for configuring the
# RPi for the IMLS/10x/18F data collection pilot.
ansible_pull_playbook () {
    _status "Installing hardening playbook."
    ansible-galaxy collection install devsec.hardening

    # 20210427 MCJ Make sure the working dir is always fresh
    sudo rm -rf $PLAYBOOK_WORKING_DIR
    sudo mkdir -p $PLAYBOOK_WORKING_DIR
    sudo chown -R pi:pi $PLAYBOOK_WORKING_DIR

    pushd $PLAYBOOK_WORKING_DIR
        _status "Cloning the playbook: ${PLAYBOOK_URL}"
        # 20210427 MCJ Adding a shallow clone.
        git clone --depth=1 $PLAYBOOK_URL
        pushd $PLAYBOOK_REPOS
            _status "Running the playbook. This will take a while."
            # For testing/dev purposes, we might not want to lock things down
            # when we're done. The lockdown flag is required to run the 
            # hardening and lockdown roles.

            # -z checks if the var is UNSET.
            if [[ -z "${NOLOCKDOWN}" ]]; then            
                ansible-playbook -i inventory.yaml playbook.yaml --extra-vars "lockdown=yes"
            else
                _status "Running playbook WITHOUT lockdown"
                ansible-playbook -i inventory.yaml playbook.yaml
            fi
            ANSIBLE_EXIT_STATUS=$?
        popd
    popd
    _status "Done running playbook."
    if [ $ANSIBLE_EXIT_STATUS -ne 0 ]; then
        _err "Ansible playbook failed."
        _err "Exit code: ${ANSIBLE_EXIT_STATUS}"
        _err "Check the log: ${SETUP_LOGFILE}"
    fi
}

disable_interactive_login () {
    # https://www.raspberrypi.org/forums/viewtopic.php?t=21632
    # Disables console and desktop login using the builtin script.
    # This tells the pi to boot to the console login, but not to auto-login `pi`
    # https://github.com/RPi-Distro/raspi-config/blob/master/raspi-config#L1308
    sudo /usr/bin/raspi-config nonint do_boot_behaviour B1
}

main () {
    check_for_usb_wifi
    if [[ -z "${NOKEYREAD}" ]]; then 
        # If NOREAD is undefined, we should read in the config.
        read_initial_configuration
    else
        echo " -- SKIPPING CONFIG ENTRY FOR TESTING PURPOSES --"
    fi
    create_logfile
    setup_logging
    bootstrap_ansible
    install_prerequisites
    ansible_pull_playbook
    disable_interactive_login
    if [ $SOMETHING_WENT_WRONG -ne 0 ]; then
        _err "Things finished with errors."
        _err "We may need to see the logs: ${SETUP_LOGFILE}"
    else
        _status "All done!"
        _status "We're rebooting in one minute!"

        # If the NOREBOOT flag is NOT set, then reboot.
        if [[ -z "${NOREBOOT}" ]]; then 
            sleep 60
            sudo reboot
        else             
            _status "Reboot prevented by env flag."
        fi
    fi
}

main

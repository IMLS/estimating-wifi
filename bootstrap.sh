#!/bin/bash

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
INITIAL_CONFIGURATION_BINARY_URL="https://github.com/jadudm/input-initial-configuration/releases/download/v0.0.2/input-initial-configuration"

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

# PURPOSE
# Sets up redirects so that STDOUT and STDERR make their way to
# a temporary logfile.
setup_logging () {
    # https://serverfault.com/questions/103501/how-can-i-fully-log-all-bash-scripts-actions
    # Save all the pipes.
    # 3 is Stdout. 4 is stderr.
    exec 3>&1 4>&2
    # Restore some.
    trap 'exec 2>&4 1>&3' 0 1 2 3
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

read_initial_configuration () {
    # Fetch the binary.
    pushd /tmp
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
            ansible-playbook -i inventory.yaml playbook.yaml --extra-vars "lockdown=yes"
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

main () {
    if [[ -z "${NOREAD}" ]]; then 
        # If NOREAD is undefined, we should read in the config.
        read_initial_configuration
    else
        _debug "Skipping config entry for dev purposes."
    fi
    create_logfile
    setup_logging
    bootstrap_ansible
    install_prerequisites
    ansible_pull_playbook
    if [ $SOMETHING_WENT_WRONG -ne 0 ]; then
        _err "Things finished with errors."
        _err "We may need to see the logs: ${SETUP_LOGFILE}"
    else
        _status "All done!"
        _status "We're rebooting in one minute!"
        sleep 60
        sudo reboot
    fi
}

main

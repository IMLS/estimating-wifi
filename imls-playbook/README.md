# imls-client-pi-playbook

This playbook is pulled repeatedly on the production machines, and is our primary mechanism for updating them "live."

"main" is "live." Therefore, all work/advancement must happen in branches, and then we can merge it back into "main" when we want it to be "live."

## hardening

Testing lockdown is tough, because... it prevents access to the machine when it is done.

Therefore, lockdown must be enabled on the command line. The `bootstrap.sh` includes lockdown, but day-to-day testing of the playbook can ignore it.

```
ansible-playbook -i inventory.yaml playbook.yaml --extra-vars "lockdown=yes"
```

Run with `preserve_ssh` and `preserve_users` so that users and SSH access are preserved while testing lockdown:

```
ansible-playbook -i inventory.yaml playbook.yaml --extra-vars "lockdown=yes, preserve_ssh=yes, preserve_users=yes"
```

## dependencies

This repository depends on the ansible [`devsec.hardening`](https://github.com/dev-sec/ansible-collection-hardening) repository.

This dependency is installed when starting up via [bootstrap.sh](./bootstrap.sh) _and_ in the [IMLS systemd service](./roles/unattended-upgrades/files/run-imls-client-pi-playbook.service).

Further dependencies should also be installed in both places.

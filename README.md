# imls-client-pi-playbook

This playbook is pulled repeatedly on the production machines, and is our primary mechanism for updating them "live."

"main" is "live." Therefore, all work/advancement must happen in branches, and then we can merge it back into "main" when we want it to be "

## hardening

Testing lockdown is tough, because... it prevents access to the machine when it is done.
Therefore, it must be enabled on the command line. The `bootstrap.sh` can include lockdown, but day-to-day testing of the playbook can ignore it.

```
ansible-playbook -i inventory.yaml playbook.yaml --extra-vars "lockdown=yes"
```

Run with `preserve_ssh` so that SSH access is preserved while testing lockdown

```
ansible-playbook -i inventory.yaml playbook.yaml --extra-vars "lockdown=yes, preserve_ssh=yes"
```

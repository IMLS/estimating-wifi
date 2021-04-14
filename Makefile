.PHONY: test

test:
	ansible-playbook -i inventory.yaml playbook.yaml --extra-vars "lockdown=yes, preserve_ssh=yes"


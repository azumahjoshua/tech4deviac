#!/bin/bash
source $VIRTUAL_ENV/bin/activate
ansible-playbook playbook.yml \
  --vault-password-file $VIRTUAL_ENV/.vault_pass \
  -i inventory
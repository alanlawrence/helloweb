#!/bin/bash
echo "Configuring git ..."
set -x # prints commands
git config --global user.email "alan.gill.lawrence@gmail.com"
git config --global user.name "alanlawrence"
git config --global core.editor vim
set +x # turns off printing commands

# For setting the Git access token, if not already done:
#  * see drive file git-access-token.txt

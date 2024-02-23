#!/bin/bash

set -e

GOOS=linux GOARCH=amd64 go build .
rsync -avz -e "ssh -i ~/.ssh/snake_team7V2" starter-snake-go root@165.232.75.224:~/
ssh -i ~/.ssh/snake_team7V2 root@165.232.75.224 "systemctl restart snake.service"
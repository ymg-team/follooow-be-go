#!/bin/bash
git pull origin main && go build main.go && pm2 restart follooow-be
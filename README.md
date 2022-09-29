# ethereum-monitor

docker run --rm --detach -p 8200:8200 -e 'VAULT_DEV_ROOT_TOKEN_ID=dev-only-token' vault:1.11.0

ocker stop "${container_id}" > /dev/null
#!/bin/sh
COMMAND="laqpay --data-dir $DATA_DIR --wallet-dir $WALLET_DIR $@"

adduser -D -u 10000 laqpay

if [[ \! -d $DATA_DIR ]]; then
    mkdir -p $DATA_DIR
fi
if [[ \! -d $WALLET_DIR ]]; then
    mkdir -p $WALLET_DIR
fi

chown -R laqpay:laqpay $( realpath $DATA_DIR )
chown -R laqpay:laqpay $( realpath $WALLET_DIR )

su laqpay -c "$COMMAND"

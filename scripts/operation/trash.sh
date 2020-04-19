#!/bin/sh

function generate_additional_key() {
    if test "$1" = "" ; then
        echo argument1 is required as account
    fi
    if test "$2" = "" ; then
        echo argument2 is required as key number
    fi

    file_name=$(coldwallet1 -k -m 30 -a "$1")

    if [ $1 != client ] && [ $1 != authorization ]; then
    fi

    if [ $1 != authorization ]; then
    fi
}

function quoine_to_payment() {
    #run after 6confirmation, so monitoring is required
    #check_confirmation payment
    while [ $(check_confirmation payment) -eq 0 ];do
        echo 'waiting payment for confirmation until 6...' && sleep 60;
    done
}

function check_confirmation() {
    json_data=$(bitcoin-cli -rpcconnect=111.111.111.111 -rpcport=18332 -rpcuser=hiromaily -rpcpassword=hiromaily listunspent 6)
    len=$(echo $json_data | jq length)
    for i in $( seq 0 $(($len - 1)) ); do
        row=$(echo $json_data | jq .[$i])
        account=$(echo $row | jq '.account')
        if [ -z "$account" ]; then
            account=$(echo $row | jq '.label')
        fi
        if [ `echo ${account} | grep ${1}` ] ; then
            conf=$(echo $row | jq '.confirmations')
            if [ $conf -ge 6 ]; then
                ret=1
            fi
        fi
    done

    echo $ret
}

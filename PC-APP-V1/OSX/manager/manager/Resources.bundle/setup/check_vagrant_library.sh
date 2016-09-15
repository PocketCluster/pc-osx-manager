#!/bin/bash

VAG_HOME="$(type -p vagrant)"
VIR_HOME="$(type -p virtualbox)"

if [[ ! -n $VIR_HOME ]] || [[ ! -x $VIR_HOME ]]; then
    exit 27
fi

if [[ ! -n $VAG_HOME ]] || [[ ! -x $VAG_HOME ]]; then
	exit 28
fi

exit 0
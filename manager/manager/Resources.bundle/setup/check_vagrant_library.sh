#!/bin/bash

VAG_HOME="$(type -p vagrant)"
VIR_HOME="$(type -p virtualbox)"

if [[ ! -n $VAG_HOME ]] || [[ ! -x $VAG_HOME ]]; then
	exit 404
fi

if [[ ! -n $VIR_HOME ]] || [[ ! -x $VIR_HOME ]]; then
	exit 406
fi

exit 0
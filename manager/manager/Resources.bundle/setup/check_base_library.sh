#!/bin/bash

JHOME="$(type -p java)"
BHOME="$(type -p brew)"

if [[ ! -n $JHOME ]] || [[ ! -x $JHOME ]]; then
	exit 25
fi

if [[ ! -n $BHOME ]] || [[ ! -x $BHOME ]]; then
	exit 26
fi

exit 0
#!/bin/bash

JHOME="$(type -p java)"
BHOME="$(type -p brew)"

if [[ ! -n $JHOME ]] || [[ ! -x $JHOME ]]; then
	exit 400
fi

if [[ ! -n $BHOME ]] || [[ ! -x $BHOME ]]; then
	exit 402
fi

exit 0
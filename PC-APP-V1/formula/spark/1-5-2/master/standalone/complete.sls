2.4.0/namenode1:
  cmd.script:
    - source: salt://spark/1-5-2/master/standalone/complete.sh
    - env:
        - NUM_NODES: '{{ salt['pillar.get']('numnodes') }}'
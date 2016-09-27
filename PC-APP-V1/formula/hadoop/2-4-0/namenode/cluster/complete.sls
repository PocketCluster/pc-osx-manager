2.4.0/namenode1:
  cmd.script:
    - source: salt://hadoop/2-4-0/namenode/cluster/complete.sh
    - env:
        - NUM_NODES: '{{ salt['pillar.get']('numnodes') }}'
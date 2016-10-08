r-base.repo:
  pkgrepo.managed:
    - humanname: R Base
    - name: ppa:marutter/rrutter
    - dist: trusty
    - file: /etc/apt/sources.list.d/cran-r.list
    - keyid: E084DAB9
    - keyserver: keyserver.ubuntu.com
    - refresh_db: True
    - require_in:
      - pkg: r-base

r-base:
  pkg:
    - installed
#    - fromrepo: ppa:marutter/rrutter
#    - name: r-base
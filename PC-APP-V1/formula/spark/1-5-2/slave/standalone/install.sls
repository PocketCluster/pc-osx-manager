1.5.2/slave:
  cmd.script:
    - source: salt://spark/1-5-2/slave/standalone/install.sh

conf.bashrc:
  file:
    - managed
    - name: /pocket/conf/spark/1.5.2/standalone/conf.bashrc
    - source: salt://spark/1-5-2/slave/standalone/conf.bashrc
    - user: pocket
    - group: pocket
    - mode: 644
    - template: jinja

spark-env.sh:
  file:
    - managed
    - name: /pocket/conf/spark/1.5.2/standalone/spark-env.sh
    - source: salt://spark/1-5-2/slave/standalone/spark-env.sh
    - user: pocket
    - group: pocket
    - mode: 644
    - template: jinja

spark-defaults.conf:
  file:
    - managed
    - name: /pocket/conf/spark/1.5.2/standalone/spark-defaults.conf
    - source: salt://spark/1-5-2/slave/standalone/spark-defaults.conf
    - user: pocket
    - group: pocket
    - mode: 644
    - template: jinja

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
      - pkg: r-base-core

r-base-core:
  pkg:
    - installed

install_scala:
  pkg.installed:
    - sources:
      - scala: http://downloads.typesafe.com/scala/2.11.7/scala-2.11.7.deb

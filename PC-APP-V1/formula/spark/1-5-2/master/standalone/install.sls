1.5.2/master:
  cmd.script:
    - source: salt://spark/1-5-2/master/standalone/install.sh

conf.bashrc:
  file:
    - managed
    - name: /pocket/conf/spark/1.5.2/standalone/conf.bashrc
    - source: salt://spark/1-5-2/master/standalone/conf.bashrc
    - group: staff
    - mode: 644
    - template: jinja

slaves.sh:
  file:
    - managed
    - name: /bigpkg/spark-1.5.2-bin-hadoop2.4/sbin/slaves.sh
    - source: salt://spark/1-5-2/master/standalone/slaves.sh
    - group: staff
    - mode: 755
    - template: jinja

spark-env.sh:
  file:
    - managed
    - name: /pocket/conf/spark/1.5.2/standalone/spark-env.sh
    - source: salt://spark/1-5-2/master/standalone/spark-env.sh
    - group: staff
    - mode: 644
    - template: jinja

spark-defaults.conf:
  file:
    - managed
    - name: /pocket/conf/spark/1.5.2/standalone/spark-defaults.conf
    - source: salt://spark/1-5-2/master/standalone/spark-defaults.conf
    - group: staff
    - mode: 644
    - template: jinja

spark-start.sh:
  file:
    - managed
    - name: /bigpkg/spark-1.5.2-bin-hadoop2.4/bin/spark-start.sh
    - source: salt://spark/1-5-2/master/standalone/spark-start.sh
    - group: staff
    - mode: 755
    - template: jinja

spark-stop.sh:
  file:
    - managed
    - name: /bigpkg/spark-1.5.2-bin-hadoop2.4/bin/spark-stop.sh
    - source: salt://spark/1-5-2/master/standalone/spark-stop.sh
    - group: staff
    - mode: 755
    - template: jinja

metastore_db:
  file.directory:
    - name: /bigpkg/spark-1.5.2-bin-hadoop2.4/metastore_db
    - group: staff
    - dir_mode: 755

work:
  file.directory:
    - name: /bigpkg/spark-1.5.2-bin-hadoop2.4/work
    - group: staff
    - dir_mode: 755

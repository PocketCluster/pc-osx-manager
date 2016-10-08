2.4.0/namenode1:
  cmd.script:
    - source: salt://hadoop/2-4-0/namenode/cluster/install.sh

capacity-scheduler.xml:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/capacity-scheduler.xml
    - source: salt://hadoop/2-4-0/namenode/cluster/capacity-scheduler.xml
    - group: staff
    - mode: 644
    - template: jinja

conf.bashrc:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/conf.bashrc
    - source: salt://hadoop/2-4-0/namenode/cluster/conf.bashrc
    - group: staff
    - mode: 644
    - template: jinja

configuration.xsl:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/configuration.xsl
    - source: salt://hadoop/2-4-0/namenode/cluster/configuration.xsl
    - group: staff
    - mode: 644
    - template: jinja

core-site.xml:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/core-site.xml
    - source: salt://hadoop/2-4-0/namenode/cluster/core-site.xml
    - group: staff
    - mode: 644
    - template: jinja

hadoop-env.sh:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/hadoop-env.sh
    - source: salt://hadoop/2-4-0/namenode/cluster/hadoop-env.sh
    - group: staff
    - mode: 644
    - template: jinja

hadoop-metrics.properties:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/hadoop-metrics.properties
    - source: salt://hadoop/2-4-0/namenode/cluster/hadoop-metrics.properties
    - group: staff
    - mode: 644
    - template: jinja

hadoop-policy.xml:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/hadoop-policy.xml
    - source: salt://hadoop/2-4-0/namenode/cluster/hadoop-policy.xml
    - group: staff
    - mode: 644
    - template: jinja

hdfs-site.xml:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/hdfs-site.xml
    - source: salt://hadoop/2-4-0/namenode/cluster/hdfs-site.xml
    - group: staff
    - mode: 644
    - template: jinja

hadoop-start.sh:
  file:
    - managed
    - name: /bigpkg/hadoop-2.4.0/bin/hadoop-start.sh
    - source: salt://hadoop/2-4-0/namenode/cluster/hadoop-start.sh
    - group: staff
    - mode: 644
    - template: jinja

hadoop-stop.sh:
  file:
    - managed
    - name: /bigpkg/hadoop-2.4.0/bin/hadoop-stop.sh
    - source: salt://hadoop/2-4-0/namenode/cluster/hadoop-stop.sh
    - group: staff
    - mode: 644
    - template: jinja

httpfs-log4j.properties:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/httpfs-log4j.properties
    - source: salt://hadoop/2-4-0/namenode/cluster/httpfs-log4j.properties
    - group: staff
    - mode: 644
    - template: jinja

log4j.properties:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/log4j.properties
    - source: salt://hadoop/2-4-0/namenode/cluster/log4j.properties
    - group: staff
    - mode: 644
    - template: jinja

masters:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/masters
    - source: salt://hadoop/2-4-0/namenode/cluster/masters
    - group: staff
    - mode: 644
    - template: jinja

mapred-env.sh:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/mapred-env.sh
    - source: salt://hadoop/2-4-0/namenode/cluster/mapred-env.sh
    - group: staff
    - mode: 644
    - template: jinja

mapred-site.xml:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/mapred-site.xml
    - source: salt://hadoop/2-4-0/namenode/cluster/mapred-site.xml
    - group: staff
    - mode: 644
    - template: jinja

slaves.sh:
  file:
    - managed
    - name: /bigpkg/hadoop-2.4.0/sbin/slaves.sh
    - source: salt://hadoop/2-4-0/namenode/cluster/slaves.sh
    - group: staff
    - mode: 755
    - template: jinja

task-controller.cfg:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/task-controller.cfg
    - source: salt://hadoop/2-4-0/namenode/cluster/task-controller.cfg
    - group: staff
    - mode: 644
    - template: jinja

yarn-env.sh:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/yarn-env.sh
    - source: salt://hadoop/2-4-0/namenode/cluster/yarn-env.sh
    - group: staff
    - mode: 644
    - template: jinja

yarn-site.xml:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/yarn-site.xml
    - source: salt://hadoop/2-4-0/namenode/cluster/yarn-site.xml
    - group: staff
    - mode: 644
    - template: jinja

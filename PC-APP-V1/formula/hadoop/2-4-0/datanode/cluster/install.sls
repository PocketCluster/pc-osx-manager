2.4.0/datanode:
  cmd.script:
    - source: salt://hadoop/2-4-0/datanode/cluster/install.sh
    - cwd: /home/pocket
    - user: pocket

conf.bashrc:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/conf.bashrc
    - source: salt://hadoop/2-4-0/datanode/cluster/conf.bashrc
    - user: pocket
    - group: pocket
    - mode: 644
    - template: jinja

configuration.xsl:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/configuration.xsl
    - source: salt://hadoop/2-4-0/datanode/cluster/configuration.xsl
    - user: pocket
    - group: pocket
    - mode: 644
    - template: jinja

core-site.xml:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/core-site.xml
    - source: salt://hadoop/2-4-0/datanode/cluster/core-site.xml
    - user: pocket
    - group: pocket
    - mode: 644
    - template: jinja

hadoop-env.sh:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/hadoop-env.sh
    - source: salt://hadoop/2-4-0/datanode/cluster/hadoop-env.sh
    - user: pocket
    - group: pocket
    - mode: 644
    - template: jinja

hadoop-metrics.properties:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/hadoop-metrics.properties
    - source: salt://hadoop/2-4-0/datanode/cluster/hadoop-metrics.properties
    - user: pocket
    - group: pocket
    - mode: 644
    - template: jinja

hadoop-policy.xml:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/hadoop-policy.xml
    - source: salt://hadoop/2-4-0/datanode/cluster/hadoop-policy.xml
    - user: pocket
    - group: pocket
    - mode: 644
    - template: jinja

hdfs-site.xml:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/hdfs-site.xml
    - source: salt://hadoop/2-4-0/datanode/cluster/hdfs-site.xml
    - user: pocket
    - group: pocket
    - mode: 644
    - template: jinja

httpfs-log4j.properties:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/httpfs-log4j.properties
    - source: salt://hadoop/2-4-0/datanode/cluster/httpfs-log4j.properties
    - user: pocket
    - group: pocket
    - mode: 644
    - template: jinja

log4j.properties:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/log4j.properties
    - source: salt://hadoop/2-4-0/datanode/cluster/log4j.properties
    - user: pocket
    - group: pocket
    - mode: 644
    - template: jinja

mapred-env.sh:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/mapred-env.sh
    - source: salt://hadoop/2-4-0/datanode/cluster/mapred-env.sh
    - user: pocket
    - group: pocket
    - mode: 644
    - template: jinja

mapred-site.xml:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/mapred-site.xml
    - source: salt://hadoop/2-4-0/datanode/cluster/mapred-site.xml
    - user: pocket
    - group: pocket
    - mode: 644
    - template: jinja

yarn-env.sh:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/yarn-env.sh
    - source: salt://hadoop/2-4-0/datanode/cluster/yarn-env.sh
    - user: pocket
    - group: pocket
    - mode: 644
    - template: jinja

yarn-site.xml:
  file:
    - managed
    - name: /pocket/conf/hadoop/2.4.0/cluster/yarn-site.xml
    - source: salt://hadoop/2-4-0/datanode/cluster/yarn-site.xml
    - user: pocket
    - group: pocket
    - mode: 644
    - template: jinja

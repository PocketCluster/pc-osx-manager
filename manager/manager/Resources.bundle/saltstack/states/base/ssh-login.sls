/home/pocket/.ssh:
  file.directory:
    - user: pocket
    - group: pocket
    - mode: 700
    - makedirs: True

config:
  file:
    - managed
    - name: /home/pocket/.ssh/config
    - source: salt://base/ssh/config
    - user: pocket
    - group: pocket
    - mode: 600
    - template: jinja

authorized_keys:
  file.append:
    - name: /home/pocket/.ssh/authorized_keys
    - source: salt://base/ssh/authorized_keys

id_rsa:
  file:
    - managed
    - name: /home/pocket/.ssh/id_rsa
    - source: salt://base/ssh/id_rsa
    - user: pocket
    - group: pocket
    - mode: 600
    - template: jinja

id_rsa.pub:
  file:
    - managed
    - name: /home/pocket/.ssh/id_rsa.pub
    - source: salt://base/ssh/id_rsa.pub
    - user: pocket
    - group: pocket
    - mode: 600
    - template: jinja

known_hosts:
  file:
    - managed
    - name: /home/pocket/.ssh/known_hosts
    - source: salt://base/ssh/known_hosts
    - user: pocket
    - group: pocket
    - mode: 600
    - template: jinja

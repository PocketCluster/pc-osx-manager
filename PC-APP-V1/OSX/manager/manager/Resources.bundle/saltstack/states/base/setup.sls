/etc/default:
  file.directory:
    - name: /etc/default
    - makedirs: True

locale:
  file:
    - managed
    - name: /etc/default/locale
    - source: salt://base/locale
    - mode: 644
    - template: jinja
import ConfigParser

if __name__ == '__main__':

    config = ConfigParser.RawConfigParser()

    # When adding sections or items, add them in the reverse order of
    # how you want them to be displayed in the actual file.
    # In addition, please note that using RawConfigParser's and the raw
    # mode of ConfigParser's respective set functions, you can assign
    # non-string values to keys internally, but will receive an error
    # when attempting to write to a file or when you get it in non-raw
    # mode. SafeConfigParser does not allow such assignments to take place.
    config.add_section('global')
    config.set('global', 'version', '1.0.0')

    config.add_section('agent')
    config.set('agent', 'ip4', '')
    config.set('agent', 'cert', '')

    config.add_section('node')
    config.set('node', 'name', '')
    config.set('node', 'ip4', '')
    config.set('node', 'gateway', '')
    config.set('node', 'netmask', '')


    # Writing our configuration file to 'example.cfg'
    with open('/etc/pocket/config.ini', 'wb') as configfile:
        config.write(configfile)

    exit(0)
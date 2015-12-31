__author__ = 'stkim1'

from pocketd.const import *

def redef_timezone(timezone="Etc/UTC"):
    with open(HOST_TIMEZONE, "w") as timezone_file:
        timezone_file.write(timezone)

def get_timezone():
    with open(HOST_TIMEZONE, "r") as timezone_file:
        return timezone_file.read()


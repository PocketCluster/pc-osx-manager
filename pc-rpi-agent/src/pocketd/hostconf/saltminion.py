__author__ = 'stkim1'

from pocketd.const import *

def redef_salt_minion_id(minion_id="pc-node"):
    with open(SALT_MINION_FILE, "w") as minionfile:
        minionfile.write(minion_id)

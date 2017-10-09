#include "sharedfolder.h"

void
DisplayVBoxSharedFolders(VBoxSharedFolder** sfolders, int len) {
    for (int i = 0; i < len; i++) {
        printf("C-Side >> Folder Name : %s | Path : %s\n", sfolders[i]->Name, sfolders[i]->Path);
    }
}
#include <stdio.h>

typedef struct VBoxSharedFolder {
    char*    Name;
    char*    Path;
} VBoxSharedFolder;

extern void
DisplayVBoxSharedFolders(VBoxSharedFolder** sfolders, int len);
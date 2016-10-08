//
//  libvboxcom.h
//  PC-MASTER-VM
//
//  Created by Almighty Kim on 7/12/16.
//  Copyright © 2016 io.pocketcluster. All rights reserved.
//


#ifndef __LIBVBOXCOM_H__
#define __LIBVBOXCOM_H__

#include <Block.h>

typedef void** VOID_DPTR;

typedef enum VBRESULT {
    GOOD = 0,
    INFO,
    FATAL
} VBRESULT;

extern void VboxUtf8Free(char* cstring);

#pragma mark - APP & API VERSION
unsigned int vbox_app_version();


#pragma mark - GET MACHINE ID
VBRESULT vbox_machine_getid(VOID_DPTR vbox_machine, char** machine_id, char *error_message);

int vbox_machine_is_setting_changed(VOID_DPTR vbox_machine, VOID_DPTR vbox_session, char *error_message);

#pragma mark - MACHINE STATUS
VBRESULT vbox_machine_setting_path(VOID_DPTR vbox_machine, char** base_folder, char* error_message);


#pragma mark - INIT & CLOSE
VBRESULT vbox_init(char* error_message);

void vbox_term();


#pragma mark - SESSION INIT
VBRESULT vbox_session_init(VOID_DPTR vbox_client, VOID_DPTR vbox_session, VOID_DPTR virtualbox, char* error_message);

VBRESULT vbox_session_close(VOID_DPTR vbox_client, VOID_DPTR vbox_session, VOID_DPTR virtualbox, char* error_message);


#pragma mark - FIND, BUILD & DESTROY MACHINE
VBRESULT vbox_machine_find(VOID_DPTR virtualbox, VOID_DPTR vbox_machine , const char* machine_name, char* error_message);

VBRESULT vbox_machine_create(VOID_DPTR virtualbox, VOID_DPTR vbox_machine, const char* machine_name, char** base_folder, char* error_message);

VBRESULT vbox_machine_release(VOID_DPTR vbox_machine, char* base_folder, char* error_message);


#pragma mark - BUILD MACHINE BASE
VBRESULT vbox_machine_build(VOID_DPTR virtualbox, VOID_DPTR vbox_machine, int cpu_count, int memory_size, char* error_message);

VBRESULT vbox_machine_add_bridged_network(VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* host_interface, char* error_message);

VBRESULT vbox_machine_add_shared_folder(VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* shared_name, const char *host_folder, char* error_message);


VBRESULT vbox_machine_add_storage_controller(VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* storage_controller_name, char* error_message);

VBRESULT vbox_machine_add_boot_image(VOID_DPTR virtualbox, VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* storage_controller_name, const char *boot_image_path, char *error_message);

VBRESULT vbox_machine_add_hard_disk(VOID_DPTR virtualbox, VOID_DPTR vbox_machine, VOID_DPTR vbox_session, const char* storage_controller_name, const char *hdd_medium_path, int disk_size, void(^build_progress)(int progress, int done), char *error_message);


#pragma mark - DELETE & RELEASE MACHINE
VBRESULT vbox_machine_destroy(VOID_DPTR vbox_machine, char* base_folder, const char* storage_controller_name, int remove_dvd, void(^build_progress)(int progress, int done), char *error_message);

#endif /* __LIBVBOXCOM_H__ */

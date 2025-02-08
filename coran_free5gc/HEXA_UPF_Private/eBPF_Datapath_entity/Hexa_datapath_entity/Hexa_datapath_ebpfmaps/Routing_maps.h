#pragma once
#ifndef __INTERFACES_H__
#define __INTERFACES_H__

#include <linux/types.h>
#include <stdint.h>
#include <bpf/bpf_helpers.h>
#include <linux/bpf.h>

typedef enum {
  N3_INTERFACE,
  N6_INTERFACE,
  N4_INTERFACE,
  N9_INTERFACE,
  N19_INTERFACE
} UPF_INTERFACES;

#define MAX_ARP_ENTRIES 12


struct s_interface {
  __uint32_t ipv4_address;
  __uint32_t port;
  const char* if_name;
};

struct Mac_address {
  uint8_t mac_address[6];
};


// struct bpf_map_def SEC("maps") Arp_map = {
//     .type        = BPF_MAP_TYPE_HASH,
//     .key_size    = sizeof(UPF_INTERFACES),                   // IPv4 address
//     .value_size  = sizeof(struct Mac_address),  // <IP Address, MAC address>
//     .max_entries = MAX_ARP_ENTRIES,          // 2,
// };
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, MAX_ARP_ENTRIES);        // Set max number of entries for the map
    __type(key,  UPF_INTERFACES);   // Define key type (e.g., ARP table key)
    __type(value, struct Mac_address); // Define value type (e.g., ARP table value)
} Arp_map SEC(".maps");



#endif  // __INTERFACES_H__

#ifndef __ARP_TABLE_MAP_H__
#define __ARP_TABLE_MAP_H__




#endif  // __ARP_TABLE_MAP_H__
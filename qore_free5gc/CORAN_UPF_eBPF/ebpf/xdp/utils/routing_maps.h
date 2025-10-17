#ifndef __INTERFACES_H__
#define __INTERFACES_H__

// #include "types.h"
// #include <stdint.h>
#include <bpf/bpf_helpers.h>
#include <linux/bpf.h>

typedef enum {
  N3_INTERFACE,
  N6_INTERFACE,
  N4_INTERFACE,
  N9_INTERFACE,
  N19_INTERFACE
} e_reference_point;

#define ARP_ENTRIES_MAX_SIZE 12


struct s_interface {
  __uint32_t ipv4_address;
  __uint32_t port;
  const char* if_name;
};

struct s_arp_mapping {
  uint8_t mac_address[6];
};


// struct bpf_map_def SEC("maps") m_arp_table = {
//     .type        = BPF_MAP_TYPE_HASH,
//     .key_size    = sizeof(e_reference_point),                   // IPv4 address
//     .value_size  = sizeof(struct s_arp_mapping),  // <IP Address, MAC address>
//     .max_entries = ARP_ENTRIES_MAX_SIZE,          // 2,
// };
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, ARP_ENTRIES_MAX_SIZE);        // Set max number of entries for the map
    __type(key,  e_reference_point);   // Define key type (e.g., ARP table key)
    __type(value, struct s_arp_mapping); // Define value type (e.g., ARP table value)
} m_arp_table SEC(".maps");



#endif  // __INTERFACES_H__

#ifndef __ARP_TABLE_MAP_H__
#define __ARP_TABLE_MAP_H__




#endif  // __ARP_TABLE_MAP_H__
#pragma once

#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/in.h>
#include <linux/ip.h>
#include <linux/types.h>
#include <sys/socket.h>
#include "../Hexa_datapath_ebpfmaps/Routing_maps.h"




#define ARP_ENTRIES_MAX_SIZE 12

// static __always_inline void print_mac_address(const unsigned char *mac) {
//     if (!mac) {
//         bpf_printk("Invalid MAC address.");
//         return;
//     }

//     bpf_printk("MAC Address Part 1: %02X:%02X:%02X",
//                mac[0],
//                mac[1],
//                mac[2]);

//     bpf_printk("MAC Address Part 2: %02X:%02X:%02X",
//                mac[3],
//                mac[4],
//                mac[5]);
// }


// static __always_inline enum xdp_action update_dst_mac_address(UPF_INTERFACES itf, struct ethhdr* eth) {
 
// struct Mac_address* map_table;
//   __builtin_memset(&map_table, 0, sizeof(struct Mac_address));
//   map_table = bpf_map_lookup_elem(&Arp_map, &itf);

//   if (!map_table) {
//     //print_mac_address(map_table);
//     return XDP_ABORTED;
//   }

// if (!eth->h_dest) {
//     //print_mac_address(map_table);
//     return XDP_ABORTED;
//   }

// bpf_printk("MAC Address was before dest mac:");  
// eth->h_dest = map_table->mac_address; // Assuming you have a valid pointer to the memory location where you want to store the new value
// uint8_t* newDest =  // Replace ... with the actual memory location
// __builtin_memcpy(eth->h_dest, newDest, ETH_ALEN);
//  // __builtin_memcpy(eth->h_dest, &map_table->mac_address, ETH_ALEN);
  

//   return XDP_TX;
// }





static __always_inline enum xdp_action Route_uplink_ipv4_packet(struct Packet_content *ctx) {

    bpf_printk("uplink updating mac addr ");
    struct Mac_address* map_table;
  __builtin_memset(&map_table, 0, sizeof(struct Mac_address));
  UPF_INTERFACES itf = N6_INTERFACE ;
  map_table = bpf_map_lookup_elem(&Arp_map, &itf);
  if (!map_table) {
    //print_mac_address(map_table);
    return XDP_ABORTED;
  }
  struct ethhdr *eth = (struct ethhdr *)ctx->packet_context->data;
  if ((const char *)(eth + 1) > (const char *)ctx->packet_context->data_end){
        return XDP_PASS;}
if (!ctx->eth->h_dest) {
    //print_mac_address(map_table);
    return XDP_ABORTED;
  }

  //ctx->eth->h_dest= &map_table->mac_address;
  bpf_printk("eth hdest %p",&eth->h_dest);
  
  __builtin_memcpy(eth->h_dest, &map_table->mac_address , ETH_ALEN);



    return XDP_TX;
}

static __always_inline enum xdp_action Route_downlink_ipv4_packet(struct Packet_content *ctx) {

//struct ethhdr *eth = (struct ethhdr *)(void *)(uintptr_t)ctx->data;
    bpf_printk("downlink updating mac addr ");
    
    struct Mac_address* map_table;
  __builtin_memset(&map_table, 0, sizeof(struct Mac_address));
  UPF_INTERFACES itf = N3_INTERFACE ;
  map_table = bpf_map_lookup_elem(&Arp_map, &itf);
  if (!map_table) {
    //print_mac_address(map_table);
    return XDP_ABORTED;
  }
  struct ethhdr *eth = (struct ethhdr *)ctx->packet_context->data;
  if ((const char *)(eth + 1) > (const char *)ctx->packet_context->data_end){
        return XDP_PASS;}
if (!ctx->eth->h_dest) {
    //print_mac_address(map_table);
    return XDP_ABORTED;
  }
  bpf_printk("eth hdest %p",&eth->h_dest);
  
  __builtin_memcpy(eth->h_dest, &map_table->mac_address , ETH_ALEN);
return XDP_TX;
        //return update_dst_mac_address(N3_INTERFACE, ctx->eth); 

}


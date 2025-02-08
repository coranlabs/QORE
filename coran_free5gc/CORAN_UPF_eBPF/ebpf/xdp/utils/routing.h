/**
 * Copyright 2023 Edgecom LLC
 * 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * 
 *     http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#pragma once

#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/in.h>
#include <linux/ip.h>
#include <linux/types.h>
#include <sys/socket.h>
#include "routing_maps.h"


#include "xdp/utils/trace.h"

struct route_stat {
    __u64 fib_lookup_ip4_cache;
    __u64 fib_lookup_ip4_ok;
    __u64 fib_lookup_ip4_error_drop;
    __u64 fib_lookup_ip4_error_pass;
    __u64 fib_lookup_ip6_cache;
    __u64 fib_lookup_ip6_ok;
    __u64 fib_lookup_ip6_error_drop;
    __u64 fib_lookup_ip6_error_pass;
};
// typedef enum {
//   N3_INTERFACE,
//   N6_INTERFACE,
//   N4_INTERFACE,
//   N9_INTERFACE,
//   N19_INTERFACE
// } e_reference_point;

#define ARP_ENTRIES_MAX_SIZE 12
struct
{
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __type(key, __u32);
    __type(value, struct route_stat);
    __uint(max_entries, 1);
} upf_route_stat SEC(".maps");

// struct s_arp_mapping {
//   uint8_t mac_address[6];
// };

// struct {
//     __uint(type, BPF_MAP_TYPE_HASH);
//     __uint(max_entries, ARP_ENTRIES_MAX_SIZE);        // Set max number of entries for the map
//     __type(key,  e_reference_point);   // Define key type (e.g., ARP table key)
//     __type(value, struct s_arp_mapping); // Define value type (e.g., ARP table value)
// } m_arp_tablea SEC("maps");

// struct {
//     __uint(type, BPF_MAP_TYPE_HASH);
//     __uint(max_entries, 12);        // Set max number of entries for the map
//     __type(key,  e_reference_point);   // Define key type (e.g., ARP table key)
//     __type(value, struct s_arp_mapping); // Define value type (e.g., ARP table value)
// } m_arp_table SEC("maps");


#ifdef ENABLE_ROUTE_CACHE

#warning "Routing cache enabled"

#define ROUTE_CACHE_IPV4_SIZE 256
#define ROUTE_CACHE_IPV6_SIZE 256

struct route_record {
    int ifindex;
    __u8 smac[6];
    __u8 dmac[6];
};

/* ipv4 -> fib cached result */
struct
{
    __uint(type, BPF_MAP_TYPE_LRU_HASH);
    __type(key, __u32);
    __type(value, struct route_record);
    __uint(max_entries, ROUTE_CACHE_IPV4_SIZE);
} upf_route_cache_ip4 SEC(".maps");

/* ipv6 -> fib cached result */
struct
{
    __uint(type, BPF_MAP_TYPE_LRU_HASH);
    __type(key, struct in6_addr);
    __type(value, struct route_record);
    __uint(max_entries, ROUTE_CACHE_IPV6_SIZE);
} upf_route_cache_ip6 SEC(".maps");

static __always_inline void update_route_cache_ipv4(const struct bpf_fib_lookup *fib_params, __u32 daddr) {
    struct route_record route = {
        .ifindex = fib_params->ifindex,
    };
    __builtin_memcpy(route.smac, fib_params->smac, ETH_ALEN);
    __builtin_memcpy(route.dmac, fib_params->dmac, ETH_ALEN);
    bpf_map_update_elem(&upf_route_cache_ip4, &daddr, &route, BPF_ANY);
}
#endif

static __always_inline enum xdp_action update_dst_mac_address(e_reference_point itf, struct ethhdr* eth) {
 
  struct s_arp_mapping* map_table;
  __builtin_memset(&map_table, 0, sizeof(struct s_arp_mapping));

  map_table = bpf_map_lookup_elem(&m_arp_table, &itf);

  if (!map_table) {
    bpf_printk("MAC Address NOT Found for IP addr: 0x%x", map_table->mac_address);
    return XDP_PASS;
  }
  bpf_printk("MAC Address was Found for IP addr: 0x%x", map_table->mac_address);
    bpf_printk("MAC Address was before dest mac: 0x%x", eth->h_dest);
  bpf_printk("MAC Address was after source mac: 0x%x", eth->h_source);
  __builtin_memcpy(eth->h_dest, &map_table->mac_address, ETH_ALEN);
//   memcpy(p_eth->h_dest, map_table->mac_address, sizeof(p_eth->h_dest));
  bpf_printk("MAC Address was updated dest mac: 0x%x", eth->h_dest);
  bpf_printk("MAC Address was updated source mac: 0x%x", eth->h_source);

  return XDP_TX;
}


//superman


static __always_inline enum xdp_action uplink_route_ipv4(struct xdp_md *ctx, struct ethhdr *eth, const struct iphdr *ip4) {
    const __u32 key = 0;
    struct route_stat *statistic = bpf_map_lookup_elem(&upf_route_stat, &key);
    if (!statistic) {
        return XDP_ABORTED;
    }
    //update_dst_mac_address(N6_INTERFACE, eth );
    bpf_printk("updated mac addr ");

    int rc = 0;
    switch (rc) {
        case BPF_FIB_LKUP_RET_SUCCESS:
            // upf_printk("upf: bpf_fib_lookup %pI4 -> %pI4: nexthop: %pI4", &ip4->saddr, &ip4->daddr, &fib_params.ipv4_dst);
            // statistic->fib_lookup_ip4_ok += 1;

            return update_dst_mac_address(N6_INTERFACE, eth );

        default:
            // upf_printk("upf: bpf_fib_lookup %pI4 -> %pI4: %d", &ip4->saddr, &ip4->daddr, rc);
            // statistic->fib_lookup_ip4_error_pass += 1;
            return XDP_PASS; /* Let's kernel takes care */
    }
}

static __always_inline enum xdp_action downlink_route_ipv4(struct xdp_md *ctx, struct ethhdr *eth, const struct iphdr *ip4) {
    const __u32 key = 0;
    struct route_stat *statistic = bpf_map_lookup_elem(&upf_route_stat, &key);
    if (!statistic) {
        return XDP_ABORTED;
    }
    

    int rc = 0;
    switch (rc) {
        case BPF_FIB_LKUP_RET_SUCCESS:
            // upf_printk("upf: bpf_fib_lookup %pI4 -> %pI4: nexthop: %pI4", &ip4->saddr, &ip4->daddr, &fib_params.ipv4_dst);
            // statistic->fib_lookup_ip4_ok += 1;

            return update_dst_mac_address(N3_INTERFACE, eth );

        default:
            // upf_printk("upf: bpf_fib_lookup %pI4 -> %pI4: %d", &ip4->saddr, &ip4->daddr, rc);
            // statistic->fib_lookup_ip4_error_pass += 1;
            return XDP_PASS; /* Let's kernel takes care */
    }
}

static __always_inline enum xdp_action route_ipv6(struct xdp_md *ctx, struct ethhdr *eth, const struct ipv6hdr *ip6) {
    const __u32 key = 0;
    struct route_stat *statistic = bpf_map_lookup_elem(&upf_route_stat, &key);
    if (!statistic) {
        return XDP_ABORTED;
    }

    struct bpf_fib_lookup fib_params = {};
    fib_params.family = AF_INET;
    // fib_params.tos = ip6->flow_lbl;
    fib_params.l4_protocol = ip6->nexthdr;
    fib_params.sport = 0;
    fib_params.dport = 0;
    fib_params.tot_len = bpf_ntohs(ip6->payload_len);
    __builtin_memcpy(fib_params.ipv6_src, &ip6->saddr, sizeof(ip6->saddr));
    __builtin_memcpy(fib_params.ipv6_dst, &ip6->daddr, sizeof(ip6->daddr));
    fib_params.ifindex = ctx->ingress_ifindex;

    int rc = bpf_fib_lookup(ctx, &fib_params, sizeof(fib_params), 0 /*BPF_FIB_LOOKUP_OUTPUT*/);
    switch (rc) {
        case BPF_FIB_LKUP_RET_SUCCESS:
            upf_printk("upf: bpf_fib_lookup %pI6c -> %pI6c: nexthop: %pI4", &ip6->saddr, &ip6->daddr, &fib_params.ipv4_dst);
            statistic->fib_lookup_ip6_ok += 1;
            //_decr_ttl(ether_proto, l3hdr);
            __builtin_memcpy(eth->h_dest, fib_params.dmac, ETH_ALEN);
            __builtin_memcpy(eth->h_source, fib_params.smac, ETH_ALEN);
            upf_printk("upf: bpf_redirect: if=%d %lu -> %lu", fib_params.ifindex, fib_params.smac, fib_params.dmac);

            if (fib_params.ifindex == ctx->ingress_ifindex)
                return XDP_TX;

            return bpf_redirect(fib_params.ifindex, 0);
        case BPF_FIB_LKUP_RET_BLACKHOLE:
        case BPF_FIB_LKUP_RET_UNREACHABLE:
        case BPF_FIB_LKUP_RET_PROHIBIT:
            upf_printk("upf: bpf_fib_lookup %pI6c -> %pI6c: %d", &ip6->saddr, &ip6->daddr, rc);
            statistic->fib_lookup_ip6_error_drop += 1;
            return XDP_DROP;
        case BPF_FIB_LKUP_RET_NOT_FWDED:
        case BPF_FIB_LKUP_RET_FWD_DISABLED:
        case BPF_FIB_LKUP_RET_UNSUPP_LWT:
        case BPF_FIB_LKUP_RET_NO_NEIGH:
        case BPF_FIB_LKUP_RET_FRAG_NEEDED:
        default:
            upf_printk("upf: bpf_fib_lookup %pI6c -> %pI6c: %d", &ip6->saddr, &ip6->daddr, rc);
            statistic->fib_lookup_ip6_error_pass += 1;
            return XDP_PASS; /* Let's kernel takes care */
    }
}
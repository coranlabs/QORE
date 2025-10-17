#pragma once

#include <bpf/bpf_helpers.h>
#include <linux/bpf.h>
#include <linux/ipv6.h>


#define PDR_MAP_UPLINK_SIZE 4024
#define PDR_MAP_DOWNLINK_IPV4_SIZE 4024
#define PDR_MAP_DOWNLINK_IPV6_SIZE 1024
#define FAR_MAP_SIZE 4024


enum outer_header_removal_values {
    OHR_GTP_U_UDP_IPv4 = 0,
    OHR_GTP_U_UDP_IPv6 = 1,
    OHR_UDP_IPv4 = 2,
    OHR_UDP_IPv6 = 3,
    OHR_IPv4 = 4,
    OHR_IPv6 = 5,
    OHR_GTP_U_UDP_IP = 6,
    OHR_VLAN_S_TAG = 7,
    OHR_S_TAG_C_TAG = 8,
};


struct PDR {
    __u32 Far_id;
    __u32 Qer_id;
    __u8 OHR;
};

/* ipv4 -> PDR */ 
struct
{
    __uint(type, BPF_MAP_TYPE_HASH);
    __type(key, __u32);
    __type(value, struct PDR);
    __uint(max_entries, PDR_MAP_DOWNLINK_IPV4_SIZE);
} PDR_downlink_map SEC(".maps");




/* TEID -> PDR */
struct
{
    __uint(type, BPF_MAP_TYPE_HASH);
    __type(key, __u32);
    __type(value, struct PDR);
    __uint(max_entries, PDR_MAP_UPLINK_SIZE);
} PDR_uplink_map SEC(".maps");

enum far_action_mask {
    FAR_DROP = 0x01,
    FAR_FORW = 0x02,
    FAR_BUFF = 0x04,
    FAR_NOCP = 0x08,
    FAR_DUPL = 0x10,
    FAR_IPMA = 0x20,
    FAR_IPMD = 0x40,
    FAR_DFRT = 0x80,
};

enum outer_header_creation_values {
    OHC_GTP_U_UDP_IPv4 = 0x01,
    OHC_GTP_U_UDP_IPv6 = 0x02,
    OHC_UDP_IPv4 = 0x04,
    OHC_UDP_IPv6 = 0x08,
};

struct FAR {
    __u8 Action;
    __u8 OHC;
    __u32 TEID;
    __u32 Remoteip;
    __u32 Localip;
    /* first octet DSCP value in the Type-of-Service, second octet shall contain the ToS/Traffic Class mask field, which shall be set to "0xFC". */
    __u16 TLM; // transport level marking
};

/* FAR ID -> FAR */
struct
{
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __type(key, __u32);
    __type(value, struct FAR);
    __uint(max_entries, FAR_MAP_SIZE);
} FAR_map SEC(".maps");


enum gate_status {
    GATE_OPEN = 0,
    GATE_CLOSED = 1,
    GATE_RESERVED1 = 2,
    GATE_RESERVED2 = 3,
};

struct QER {
    __u8 ULGate_status;
    __u8 DLGate_status;
    __u8 Qfi;
    __u32 ULMax_bitrate;
    __u32 DLMax_bitrate;
    __u64 ul_start;
    __u64 dl_start;
};

#define QER_MAP_SIZE 1024

/* QER ID -> QER */
struct
{
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __type(key, __u32);
    __type(value, struct QER);
    __uint(max_entries, QER_MAP_SIZE);
} QER_map SEC(".maps");
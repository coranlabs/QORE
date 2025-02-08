#pragma once

/* Protocol Headers */
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/ipv6.h>
#include <linux/udp.h>
#include <linux/tcp.h>

/* Linux Types and Utilities */
#include <linux/types.h>
#include <stdint.h>



/* GTP UDP Port Number */
#define GTP_UDP_PORT 2152u

/* GTP Flags */
#define GTP_FLAGS 0x30  /* Version: GTPv1, Protocol Type: GTP, Others: 0 */

/* GTP Message Types */
#define GTPU_ECHO_REQUEST                    1
#define GTPU_ECHO_RESPONSE                   2
#define GTPU_ERROR_INDICATION                26
#define GTPU_SUPPORTED_EXTENSION_HEADERS_NOTIFICATION 31
#define GTPU_END_MARKER                      254
#define GTPU_G_PDU                           255
#define GTP_ERROR        -1

/* Custom Headers */

/*3GPP TS 29.281 version 17.4.0 Release 17  5.1*/
struct GTPU_header {
#if __BYTE_ORDER == __LITTLE_ENDIAN
    unsigned int PN    : 1;  /* Sequence Number Flag */
    unsigned int S     : 1;  /* Extension Header Flag */
    unsigned int E     : 1;  /* TEID Flag */
    unsigned int Spare : 1;
    unsigned int PT    : 1;  /* Protocol Type */
    unsigned int Version : 3;
#elif __BYTE_ORDER == __BIG_ENDIAN
    unsigned int Version : 3;
    unsigned int PT   : 1;  /* Protocol Type */
    unsigned int Spare : 1;
    unsigned int E    : 1;  /* Extension Header Flag */
    unsigned int S     : 1;  /* Sequence Number Flag */
    unsigned int PN    : 1;  /* N-PDU Number Flag */
#else
#error "Please fix <bits/endian.h>"
#endif
    __u8 Message_type;       /* GTP Message Type */
    __u16 Message_length;    /* Length of GTP Message */
    __u32 TEID;               /* Tunnel Endpoint Identifier */
    __u16 Sequence;
    __u8 N_PDU_number;
    __u8 Next_extension_header_type;
} __attribute__((packed));


struct GTPU_Base_header {
#if __BYTE_ORDER == __LITTLE_ENDIAN
    unsigned int PN    : 1;  /* Sequence Number Flag */
    unsigned int S     : 1;  /* Extension Header Flag */
    unsigned int E     : 1;  /* TEID Flag */
    unsigned int Spare : 1;
    unsigned int PT    : 1;  /* Protocol Type */
    unsigned int Version : 3;
#elif __BYTE_ORDER == __BIG_ENDIAN
    unsigned int Version : 3;
    unsigned int PT   : 1;  /* Protocol Type */
    unsigned int Spare : 1;
    unsigned int E    : 1;  /* Extension Header Flag */
    unsigned int S     : 1;  /* Sequence Number Flag */
    unsigned int PN    : 1;  /* N-PDU Number Flag */
#else
#error "Please fix <bits/endian.h>"
#endif
    __u8 Message_type;       /* GTP Message Type */
    __u16 Message_length;    /* Length of GTP Message */
    __u32 TEID;               /* Tunnel Endpoint Identifier */

} __attribute__((packed));

struct GTPU_conditional_header{

    __u16 Sequence;
    __u8 N_PDU_number;
    __u8 Next_extension_header_type;
} __attribute__((packed));

// /*3GPP TS 29.281 version 17.4.0 Release 17  5.2*/
// struct gtp_hdr_ext {
//     __u16 sqn;   /* Sequence Number */
//     __u8 npdu;        /* N-PDU Number */
//     __u8 next_ext; /* Next Extension Header Type */
// } __attribute__((packed));

/**
 * @struct Packet_content
 * @brief Structure to track packet parsing state and protocol headers
 * 
 * This structure maintains the current parsing position and pointers
 * to various protocol headers found in the packet. It's designed for
 * use in XDP (eXpress Data Path) packet processing.
 */

struct packet_context{

    char *data;          /* XDP program context */
    const char *data_end;

};

struct Packet_content {
    /* Buffer Management */

    struct packet_context volatile_packet_ctx;     /* Packet context */
    /* XDP Context */
    struct xdp_md *packet_context;          /* XDP program context */
    
    
    /* Protocol Headers */
    struct ethhdr *eth;              /* Ethernet header */
    struct iphdr *ip4;               /* IPv4 header */
    struct ipv6hdr *ip6;             /* IPv6 header */
    struct udphdr *udp;              /* UDP header */
    struct tcphdr *tcp;              /* TCP header */
    struct GTPU_header *gtp;             /* GTP-U header */
};


// GTP-U Extension Header Field Values
// 3gpp 29.281 v17.4.0 Release 17
#define GTPU_EXT_TYPE_NO_MORE_HEADERS         0x00  // 0000 0000: No more extension headers
#define GTPU_EXT_TYPE_RESERVED_CP1            0x01  // 0000 0001: Reserved - Control Plane only
#define GTPU_EXT_TYPE_RESERVED_CP2            0x02  // 0000 0010: Reserved - Control Plane only
#define GTPU_EXT_TYPE_LONG_PDCP_PDU_NUMBER1   0x03  // 0000 0011: Long PDCP PDU Number
#define GTPU_EXT_TYPE_SERVICE_CLASS_INDICATOR 0x20  // 0010 0000: Service Class Indicator
#define GTPU_EXT_TYPE_UDP_PORT                0x40  // 0100 0000: UDP Port
#define GTPU_EXT_TYPE_RAN_CONTAINER           0x81  // 1000 0001: RAN Container
#define GTPU_EXT_TYPE_LONG_PDCP_PDU_NUMBER2   0x82  // 1000 0010: Long PDCP PDU Number
#define GTPU_EXT_TYPE_XW_RAN_CONTAINER        0x83  // 1000 0011: Xw RAN Container
#define GTPU_EXT_TYPE_NR_RAN_CONTAINER        0x84  // 1000 0100: NR RAN Container
#define GTPU_EXT_TYPE_PDU_SESSION_CONTAINER   0x85  // 1000 0101: PDU Session Container
#define GTPU_EXT_TYPE_PDCP_PDU_NUMBER         0xC0  // 1100 0000: PDCP PDU Number
#define GTPU_EXT_TYPE_RESERVED_CP3            0xC1  // 1100 0001: Reserved - Control Plane only
#define GTPU_EXT_TYPE_RESERVED_CP4            0xC2  // 1100 0010: Reserved - Control Plane only
#define GTP_PSC_PDU_type_DL_PDU_SESSION_INFORMATION            0  
#define GTP_PSC_PDU_type_UL_PDU_SESSION_INFORMATION            1  



/**
 * @struct GTP_ext_type_PDU_session_container
 * @brief GTP PDU Session Container extension header
 */

//3gpp 38.415 v17.1.0 Release 17
struct GTP_ext_type_PDU_session_container {

    __u8 Extension_Header_Length;          
    __u8 Pdu_type : 4;   // Matches last 4 bits of the second byte
} __attribute__((packed));

// my design from latest
struct GTP_PDU_session_type_UL_PDU_SESSION_INFORMATION {
    __u8 Extension_Header_Length;         
    __u8 PDU_type : 4;   // Matches last 4 bits of the second byte
    __u8 QMP : 1;     //QOS monitoring policy
    __u8 DL_delay_ind : 1;     //
    __u8 UL_delay_ind : 1;    // 
    __u8 SNP : 1;     //Sequence number presence 
    __u8 N3_N9_Delay_indcator : 1;   // Matches last 1 bits of the second byte
    __u8 New_IE_Flag  : 1;     // Matches last 1 bit of the third byte
    __u8 QoS_Flow_Identifier : 6;        // Matches first 6 bits of the third byte
    __u8 Next_Extension_Header_Type;       // Matches the fourth byte
} __attribute__((packed));

//design for DL pdu session protocol
struct GTP_PDU_session_container_type_DL_PDU_SESSION_INFORMATION {
    __u8 Extension_Header_Length;         
    __u8 PDU_type : 4;   // Matches last 4 bits of the second byte
    __u8 QMP : 1;     //QOS monitoring policy
    __u8 SNP : 1;     //Sequence number presence 
    __u8 MSNP : 1;     //
    __u8 Spare : 1;    // 
    __u8 PPP : 1;   // Matches last 1 bits of the second byte
    __u8 RQI  : 1;     // Matches last 1 bit of the third byte
    __u8 QoS_Flow_Identifier : 6;        // Matches first 6 bits of the third byte
    __u8 Next_Extension_Header_Type;       // Matches the fourth byte
} __attribute__((packed));

// Define the PDU Session Container extension header structure
typedef struct {
    uint8_t pdu_type : 4;               // PDU Type (4 bits)
    uint8_t spare : 4;                  // Spare (4 bits, set to 0)
    uint8_t paging_policy_presence;      // Paging Policy Presence (1 bit)
    uint8_t reflective_qos_indicator;    // Reflective QoS Indicator (1 bit)
    uint8_t qos_flow_identifier;         // QoS Flow Identifier (6 bits)
    uint8_t next_extension_header_type;  // Next Extension Header Type (8 bits)
} pdu_session_container_header_t;


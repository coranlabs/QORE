#pragma once

#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/in.h>
#include <linux/ip.h>
#include <linux/types.h>
#include <linux/udp.h>
#include <linux/icmp.h>

#include "Packet_header_parser.h"
#include "Checksum.h"

/* 
 * GTP (GPRS Tunneling Protocol) Header Parsing and Manipulation Utilities
 * 
 * This header provides low-level functions for:
 * - Parsing GTP headers
 * - Handling GTP Echo Requests
 * - Adding/Removing GTP encapsulation
 * - Filling various network protocol headers
 *
 * Primarily used in XDP (eXpress Data Path) networking context
 */


/* 
 * Parse GTP (GPRS Tunneling Protocol) header
 * 
 * @param ctx Packet context containing packet data
 * @return Message type of GTP header or -1 on error
 * 
 * Handles extended GTP headers by skipping additional bytes if 
 * extension flags (e, s, pn) are set
 */



static __always_inline __u32 parse_gtp_extension_header_PDU_SESSION_CONTAINER( struct Packet_content *ctx ) {




    struct GTP_ext_type_PDU_session_container *gtp_ext_pdu_session_container_type = (struct GTP_ext_type_PDU_session_container *)ctx->volatile_packet_ctx.data;
    if (ctx->volatile_packet_ctx.data + sizeof(struct GTP_ext_type_PDU_session_container) > ctx->packet_context->data_end) {
    return -1;
}else{
gtp_ext_pdu_session_container_type = (struct GTP_ext_type_PDU_session_container *)ctx->volatile_packet_ctx.data;
    bpf_printk("Hexa: GTP pdutype: %d", gtp_ext_pdu_session_container_type->Pdu_type);
    }

    switch (gtp_ext_pdu_session_container_type->Pdu_type ){
        case GTP_PSC_PDU_type_DL_PDU_SESSION_INFORMATION:{
            struct GTP_PDU_session_container_type_DL_PDU_SESSION_INFORMATION *DL_PDU_session_information = (struct GTP_PDU_session_type_DL_PDU_SESSION_INFORMATION *)ctx->volatile_packet_ctx.data;
            long pointer_displacement_result = Bpf_packet_context_adjust(&ctx->volatile_packet_ctx, sizeof(*DL_PDU_session_information));
            if (pointer_displacement_result != 0){
                return -1;}
            return sizeof(*DL_PDU_session_information);}
        case GTP_PSC_PDU_type_UL_PDU_SESSION_INFORMATION:{
            struct GTP_PDU_session_type_UL_PDU_SESSION_INFORMATION *UL_PDU_session_information = (struct GTP_PDU_session_type_UL_PDU_SESSION_INFORMATION *)ctx->volatile_packet_ctx.data;
            long pointer_displacement_result = Bpf_packet_context_adjust(&ctx->volatile_packet_ctx, sizeof(*UL_PDU_session_information));
            if (pointer_displacement_result != 0){
                return -1;}
            return sizeof(*UL_PDU_session_information);}
        default:
            bpf_printk("invalid pdu type of PDU_SESSION_CONTAINER");
            return -1;
    }

    
}
// static __always_inline __u32 Parse_gtp_header(struct Packet_content *ctx) {
//     struct GTPU_header *gtp = (struct GTPU_header *)ctx->data;
//     if ((const char *)(gtp + 1) > ctx->data_end)
//         return -1;

//     long pointer_displacement_result = bpf_xdp_adjust_head(ctx->data, sizeof(*gtp));
//     if (pointer_displacement_result != 0){       
//         return -1;}
//     // Skip extended headers if extension flags are set
//     ctx->gtp = gtp;
//     return ctx->gtp->Message_type;
// }

/* 
 * Guess Ethernet protocol based on IP version
 * 
 * @param data Pointer to packet data
 * @return Ethernet protocol type (ETH_P_IP or ETH_P_IPV6) or -1 for non-IP
 */
static __always_inline int verify_eth_protocol(const char *data) {
    const __u8 ip_version = (*(const __u8 *)data) >> 4;
    switch (ip_version) {
        case 6: {
             ETH_P_IPV6_BE;  // Note: Missing return (likely a typo in original code)
        }
        case 4: {
            return ETH_P_IP_BE;
        }
        default:
            // Log non-IP packets
            bpf_printk("hexa: can't process non-IP packet: %d", ip_version);
            return -1;
    }
}

/* 
 * Remove GTP header from packet
 * 
 * @param ctx Packet context
 * @return 0 on success, negative on error
 * 
 * Removes GTP encapsulation by adjusting packet head and 
 * reinitializing packet context
 */
static __always_inline long Remove_gtp_header(struct Packet_content *ctx) {
    // Validate GTP packet
    if (!ctx->gtp) {
        bpf_printk("Hexa: Unable to remove gtp header, gtp is null");
        return -1;
    }

    // Calculate extended GTP header size
    size_t Gtp_extension_header_size = 0;
    struct GTPU_header *gtp = ctx->gtp;
    bpf_printk("Hexa: GTP header type: %d", gtp->E);
    //return 1;
    if (gtp->E){
    bpf_printk("Hexa: GTP sjdc njd type: ");
        //return -1;
        switch (gtp->Next_extension_header_type){
            case GTPU_EXT_TYPE_PDU_SESSION_CONTAINER:{
                Gtp_extension_header_size = parse_gtp_extension_header_PDU_SESSION_CONTAINER( ctx);
                bpf_printk("Hexa: GTPextension header value: %d", Gtp_extension_header_size);
                //return -1;
                break;}
            default:{
                bpf_printk("Hexa: Unsupported extended GTP header type: %d", gtp->Next_extension_header_type);
                return -1;}
        }
       }


    // Total size of GTP encapsulation
const size_t total_gtp_size = sizeof(struct iphdr) + sizeof(struct udphdr) + sizeof(struct GTPU_header) + Gtp_extension_header_size;

    // Packet data pointers
    // void *data = (void *)(long)ctx->packet_context->data;
    // const void *data_end = (const void *)(long)ctx->packet_context->data_end;
    
    // // Validate Ethernet header
    // struct ethhdr *eth = (struct ethhdr *)data;
    // if ((const void *)(eth + 1) > data_end) {
    //     bpf_printk("Hexa: ctx went out bounds ");
    //     return -1;
    // }

    // // Adjust data pointer and validate new Ethernet header
    // long pointer_displacement_result = bpf_xdp_adjust_head(ctx->volatile_xdp_ctx, sizeof(*eth));
    // if (pointer_displacement_result != 0){
    //     return -1;}
    // struct ethhdr *New_eth_hdr = (struct ethhdr *)data;
    // if ((const void *)(New_eth_hdr + 1) > data_end) {
    //     bpf_printk("Hexa: ctx went out bounds ");
    //     return -1;
    // }

    // // Guess Ethernet protocol
    // pointer_displacement_result = bpf_xdp_adjust_head(ctx->data, sizeof(*New_eth_hdr));
    // if (pointer_displacement_result != 0){
    //     return -1;}
    // const int eth_proto = Verify_eth_protocol(data);

    // if (eth_proto == -1){
    //     bpf_printk("packet eth prot issue");
    //     return -1;
    // }


    //return 1;
    if (ctx->volatile_packet_ctx.data + sizeof(struct GTP_ext_type_PDU_session_container) > ctx->packet_context->data_end) {
    return -1;}

    long pointer_displacement_result = Bpf_packet_context_adjust(&ctx->volatile_packet_ctx, (int)-sizeof(struct ethhdr));

    if (pointer_displacement_result != 0){
        return -1;}
    
    struct ethhdr *New_eth_hdr = (struct ethhdr *)ctx->volatile_packet_ctx.data;
    // Copy original Ethernet header
    //return 1;
    __builtin_memcpy(New_eth_hdr, ctx->eth, sizeof(*New_eth_hdr));
    pointer_displacement_result = bpf_xdp_adjust_head(ctx->packet_context, total_gtp_size);
    if (pointer_displacement_result != 0){
        return -1;
    }
    // ctx->packet_context->data = (__u32)(uintptr_t)ctx->volatile_packet_ctx.data;
    // pointer_displacement_result = bpf_xdp_adjust_head(ctx->packet_context, total_gtp_size);
    // bpf_printk("the result of bpfxdpadjust head is:%d ", pointer_displacement_result);
    // if (pointer_displacement_result != 0){
    //     return -1;}        


    return 1;
}

/* Remaining functions (Build_ip_header, Build_udp_header, Build_gtp_header, 
   fill_gtp_ext_header, Build_gtpu_extension_type_PSC_header, add_gtp_over_ip4_headers, 
   update_gtp_tunnel) have similar detailed comments */

static __always_inline void Build_ip_header(struct iphdr *ip, int saddr, int daddr, int tot_len) {
    ip->version = 4;
    ip->ihl = 5; /* No options */
    ip->tos = 96;
    ip->tot_len = bpf_htons(tot_len);
    ip->id = 0;            /* No fragmentation */
    ip->frag_off = 0x0040; /* Don't fragment; Fragment offset = 0 */
    ip->ttl = 64;
    ip->protocol = IPPROTO_UDP;
    ip->check = 0;
    ip->saddr = saddr;
    ip->daddr = daddr;
}

static __always_inline void Build_udp_header(struct udphdr *udp, int port, int len) {
    udp->source = bpf_htons(port);
    udp->dest = udp->source;
    udp->len = bpf_htons(len);
    udp->check = 0;
}

static __always_inline void Build_gtp_base_header(struct GTPU_header *gtp, int TEID, int len, int Extension_header_type) {
    *(__u8 *)gtp = GTP_FLAGS;
    gtp->E = 1;
    gtp->Message_type = GTPU_G_PDU;
    gtp->Message_length = bpf_htons(len);
    gtp->TEID = bpf_htonl(TEID);
    gtp->S = 0;
    gtp->N_PDU_number = 0;
    gtp->Next_extension_header_type = Extension_header_type;
}


static __always_inline void Build_PDU_session_container_DL_PDU_SESSION_INFORMATION(struct GTP_PDU_session_container_type_DL_PDU_SESSION_INFORMATION *gtp_DL_pdu_session_information, int QFI, int pdu_type) {
    gtp_DL_pdu_session_information->Extension_Header_Length = 1;
    gtp_DL_pdu_session_information->PDU_type = pdu_type;
    gtp_DL_pdu_session_information->QMP = 0;
    gtp_DL_pdu_session_information->SNP = 0;
    gtp_DL_pdu_session_information->MSNP = 0;
    gtp_DL_pdu_session_information->Spare = 0;
    gtp_DL_pdu_session_information->PPP = 0;
    gtp_DL_pdu_session_information->RQI = 0;
    gtp_DL_pdu_session_information->QoS_Flow_Identifier = QFI;
    gtp_DL_pdu_session_information->Next_Extension_Header_Type = 0;
}

static __always_inline __u32 GTP_encapsulator(struct Packet_content *ctx, int saddr, int daddr, __u8 qfi, int tos,int TEID) {
    static const size_t gtp_full_hdr_size = sizeof(struct GTPU_header) + sizeof(struct GTP_PDU_session_container_type_DL_PDU_SESSION_INFORMATION);
    static const size_t gtp_size = sizeof(struct GTP_PDU_session_container_type_DL_PDU_SESSION_INFORMATION) + sizeof(struct GTP_PDU_session_container_type_DL_PDU_SESSION_INFORMATION);
    static const size_t total_gtp_size = sizeof(struct iphdr) + sizeof(struct udphdr) + gtp_full_hdr_size;
    static const size_t GTP_base_header_size = sizeof(struct GTPU_Base_header);

    // int internal_ip_packet_len = (ctx->packet_context->data_end - ctx->packet_context->data) - sizeof(*eth);
    int internal_ip_packet_len = bpf_ntohs(ctx->ip4->tot_len);

    //void *data = (void *)(long)ctx->packet_context->data;
    
    //const void *data_end = (const void *)(long)ctx->packet_context->data_end;
    // if (data + 1 > data_end){
    //     bpf_printk("packet size issue");

    //     return -1;
    // }
    
    int pointer_displacement_result = bpf_xdp_adjust_head(ctx->packet_context, (__s32)-total_gtp_size);
    if (pointer_displacement_result != 0){
        return -1;}
    ctx->volatile_packet_ctx.data = (char *)(uintptr_t)ctx->packet_context->data;

    if (ctx->volatile_packet_ctx.data + 1 > (const char *)ctx->packet_context->data_end){
        return -1;
    }
    struct ethhdr *orig_eth = (struct ethhdr *)(ctx->volatile_packet_ctx.data + total_gtp_size);
    if ((const char *)(orig_eth + 1) > (const char *)ctx->packet_context->data_end){
        return -1;}
    if (!orig_eth){
        return -1;
    }

    //data = (void *)(long)ctx->packet_context->data;
    
    struct ethhdr *eth = (struct ethhdr *)ctx->packet_context->data;
    if ((const char *)(eth + 1) > (const char *)ctx->packet_context->data_end){
        return -1;}


    __builtin_memcpy(eth, orig_eth, sizeof(struct ethhdr));
    //return 1;
    // eth->h_proto = bpf_htons(ETH_P_IP);

    struct iphdr *ip = (struct iphdr *)(eth + 1);
    if ((const void *)(ip + 1) > (const void *)ctx->packet_context->data_end){
        return -1;}

    /* Add the outer IP header */
    //return -1;
    Build_ip_header(ip, saddr, daddr, internal_ip_packet_len + total_gtp_size);

    /* Add the UDP header */
    struct udphdr *udp = (struct udphdr *)(ip + 1);
    if ((const char *)(udp + 1) > (const void *)ctx->packet_context->data_end)
        return -1;

    Build_udp_header(udp, GTP_UDP_PORT, internal_ip_packet_len + sizeof(*udp) + gtp_full_hdr_size);

    /* Add the GTP header */
    struct GTPU_header *gtp = (struct GTPU_header *)(udp + 1);
    if ((const char *)(gtp + 1) > (const void *)ctx->packet_context->data_end)
        return -1;

    Build_gtp_base_header(gtp, TEID, gtp_size + internal_ip_packet_len, GTPU_EXT_TYPE_PDU_SESSION_CONTAINER);

    /* Add the GTP ext header */
    struct GTP_PDU_session_container_type_DL_PDU_SESSION_INFORMATION *gtp_DL_pdu_session_information = (struct GTP_PDU_session_container_type_DL_PDU_SESSION_INFORMATION *)(gtp + 1);
    if ((const char *)(gtp_DL_pdu_session_information + 1) > (const void *)ctx->packet_context->data_end)
        return -1;

    Build_PDU_session_container_DL_PDU_SESSION_INFORMATION(gtp_DL_pdu_session_information,qfi, GTP_PSC_PDU_type_DL_PDU_SESSION_INFORMATION);



    ip->check = Generate_ipv4_checksum(ip, sizeof(*ip));

    /* TODO: implement UDP Checksum_difference which pass ebpf verifier checks successfully */
    // cs = 0;
    // const void* udp_start = (void*)udp;
    // const __u16 udp_len = bpf_htons(udp->len);
    // ipv4_l4_Checksum_difference(udp, udp_len, &cs, ip);
    // udp->check = cs;

    /* Update packet pointers */
    // ctx->packet_context->data = data;
    // ctx->packet_context->data_end = data_end;
    return 0;
}

static __always_inline void update_gtp_tunnel(struct Packet_content *ctx, int srcip, int dstip, __u8 tos, int TEID) {

    ctx->gtp->TEID = bpf_htonl(TEID);
    ctx->ip4->saddr = srcip;
    ctx->ip4->daddr = dstip;
    ctx->ip4->check = 0;
    ctx->ip4->check = Generate_ipv4_checksum(ctx->ip4, sizeof(*ctx->ip4));
}

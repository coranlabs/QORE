
#pragma once
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/types.h>
#include <linux/udp.h>
#include <linux/tcp.h>

#include "../Hexa_datapath_ebpfmaps/Rules_maps.h"

#include "Packet_headers.h"

/* Constants */
#define ETH_P_IPV6_BE   0xDD86
#define ETH_P_IP_BE     0x0008
#define IP_FRAGMENTED   65343  /* 0x3FFF mask to check for fragment offset field */

/* Error codes */
#define PARSER_OK        0
#define PARSER_ERROR    -1




typedef enum {
  UPLINK_FLOW,
  DOWNLINK_FLOW,
  NON_FLOW,
  BAD_FLOW,
  UNHANDLED_FLOW,
} FLOW_STATES;

typedef struct {
    FLOW_STATES status;
    struct PDR *value;
} FLOW_CONTAINER;

/* parsing functions*/



static __always_inline int Bpf_packet_context_adjust(struct packet_context *ctx, signed int delta) {
    ctx->data += delta;
    if (ctx->data_end < ctx->data) {
        return -1;
    }
    return 0;
}



static __always_inline  FLOW_CONTAINER Parse_packet_flow(struct Packet_content *ctx) {

    FLOW_CONTAINER flow_values;
    struct ethhdr *eth = (struct ethhdr *)ctx->volatile_packet_ctx.data;
    if ((const char *)(eth + 1) > ctx->volatile_packet_ctx.data_end) {
                        flow_values.status = BAD_FLOW;
                        return flow_values;}
    __u16 eth_proto = bpf_ntohs(eth->h_proto);
    ctx->eth = eth;
    int ethernet_header_size = sizeof(*eth);
    long pointer_displacement_result = Bpf_packet_context_adjust(&ctx->volatile_packet_ctx, ethernet_header_size);
    if (pointer_displacement_result != 0){
        flow_values.status = BAD_FLOW;
        return flow_values;}
    bpf_printk("ethernet protocol is %d", ctx->eth->h_proto);
    switch (ctx->eth->h_proto) {
        case ETH_P_IPV6_BE:{
            flow_values.status = BAD_FLOW;
            return flow_values;}
        case ETH_P_IP_BE:
        {
            struct iphdr *ip4 = (struct iphdr *)ctx->volatile_packet_ctx.data;
            ctx->ip4 = ip4;
            long pointer_displacement_result = Bpf_packet_context_adjust(&ctx->volatile_packet_ctx, sizeof(*ip4));
            if (pointer_displacement_result != 0){
                bpf_printk("Hexa: BAD_FLOW at ip4pointer_displacement_result is %d", pointer_displacement_result);

                flow_values.status = BAD_FLOW;
                return flow_values;}
                

            bpf_printk("ip4 protocol is %d", ctx->ip4->protocol);
            bpf_printk("ip4 protocol is %d", ctx->ip4->daddr);
            switch (ctx->ip4->protocol) {
                case IPPROTO_UDP:{
                    struct udphdr *udp = (struct udphdr *)ctx->volatile_packet_ctx.data;
                    if ((const char *)(udp + 1) > ctx->volatile_packet_ctx.data_end){
                        flow_values.status = BAD_FLOW;
                        bpf_printk("Hexa: BAD_FLOW at udp");
                        return flow_values;}

                    long pointer_displacement_result = Bpf_packet_context_adjust(&ctx->volatile_packet_ctx, sizeof(*udp));
                    if (pointer_displacement_result != 0){
                        bpf_printk("Hexa: BAD_FLOW pointer_displacement_result is %d", pointer_displacement_result);
                        flow_values.status = BAD_FLOW;
                        return flow_values;} 
                    int udp_header_size = sizeof(*udp);
                                         
                    ctx->udp = udp;
                    bpf_printk("Hexa: BAD_FLOW at udp %d ", bpf_ntohs(ctx->udp->dest));
                    
                    switch (bpf_ntohs(ctx->udp->dest)) {
                        case GTP_UDP_PORT:{
                            
                            struct GTPU_header *gtp = (struct GTPU_header *)ctx->volatile_packet_ctx.data;
                            long pointer_displacement_result = Bpf_packet_context_adjust(&ctx->volatile_packet_ctx, sizeof(*gtp));
                            if (pointer_displacement_result != 0){
                                bpf_printk("Hexa: BAD_FLOW  at udp pointer_displacement_result is %d", pointer_displacement_result);
                                flow_values.status = BAD_FLOW;
                                return flow_values;}
                            
                            

                            ctx->gtp = gtp;

                            bpf_printk("GTPU message type is %d", gtp->Message_type);
                            switch (gtp->Message_type) {
                                case GTPU_G_PDU:{
                                    flow_values.status = UPLINK_FLOW;
                                    __u32 TEID = bpf_htonl(ctx->gtp->TEID);
                                    struct PDR *PDR_content = bpf_map_lookup_elem(&PDR_uplink_map, &TEID);
                                    if (!PDR_content) {
                                        bpf_printk("Hexa: Unknown TEID:%d No session found in PDR", TEID);
                                        flow_values.status = UNHANDLED_FLOW;
                                        return flow_values;
                                    }
                                        flow_values.value = PDR_content;
                                        return flow_values;
                                     
                            //         flow_values.status = BAD_FLOW;
                            // bpf_printk("Hexa: BAD_FLOW at gtpport");
                            // return flow_values;                               
                                }
                                default:{
                                    flow_values.status = NON_FLOW;
                                    return flow_values;
                                }
                            }

                        }
                        default:{
                            // flow_values.status = BAD_FLOW;
                            // bpf_printk("Hexa: BAD_FLOW at default udpport");
                            // bpf_printk("Hexa: NO downlink session found for IP:%d", ctx->ip4->daddr);

                            // return flow_values;
                            struct PDR *PDR_content = bpf_map_lookup_elem(&PDR_downlink_map, &ctx->ip4->daddr);
                            
                            if (!PDR_content) {
                                bpf_printk("Hexa: NO downlink session found for IP:%d", ctx->ip4->daddr);
                                flow_values.status = NON_FLOW;
                                return flow_values;
                            }else{
                                flow_values.status = DOWNLINK_FLOW;
                                flow_values.value = PDR_content;
                                return flow_values;
                            }
                            
                        }
                    }
                }
                    
                case IPPROTO_TCP:{
                    struct PDR *PDR_content = bpf_map_lookup_elem(&PDR_downlink_map, &ctx->ip4->daddr);
                            
                            if (!PDR_content) {
                                bpf_printk("Hexa: NO downlink session found for IP:%d", ctx->ip4->daddr);
                                flow_values.status = NON_FLOW;
                                return flow_values;
                            }else{
                                flow_values.status = DOWNLINK_FLOW;
                                flow_values.value = PDR_content;
                                return flow_values;
                            }
                }
                default:{
                    //bpf_printk("Hexa: Bdefaut UFLOW");
                    struct PDR *PDR_content = bpf_map_lookup_elem(&PDR_downlink_map, &ctx->ip4->daddr);
                            
                            if (!PDR_content) {
                                bpf_printk("Hexa: NO downlink session found for IP:%d", ctx->ip4->daddr);
                                flow_values.status = UNHANDLED_FLOW;
                                return flow_values;
                            }else{
                                flow_values.status = DOWNLINK_FLOW;
                                flow_values.value = PDR_content;
                                return flow_values;
                            }
                    
                }
            }
        }

        default:{
            bpf_printk("Hexa: defaut UFLOW");
            flow_values.status = UNHANDLED_FLOW;
            return flow_values;
        }
    }

}
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>
#include <linux/in.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/ipv6.h>
#include <linux/tcp.h>
#include <linux/udp.h>
#include <sys/socket.h>

#include "Hexa_datapath_utilities/Gtpu_header_operations.h"
#include "Hexa_datapath_utilities/Router.h"





static __always_inline enum xdp_action Handle_buffer_action(struct Packet_content *ctx, struct FAR *Far_content, struct PDR *PDR_content) {
    return XDP_PASS;
}
static __always_inline enum xdp_action Send_downlink_ipv4_packet(struct Packet_content *ctx, int Srcip, int Destip, __u8 tos, __u8 Qfi, int TEID) {
    //return XDP_PASS;
    if (-1 == GTP_encapsulator(ctx, Srcip, Destip ,Qfi, tos,TEID)){
        return XDP_ABORTED;}
    bpf_printk("Hexa: send downlink gtp pdu %pI4 -> %pI4", &ctx->ip4->saddr, &ctx->ip4->daddr);
    return Route_downlink_ipv4_packet(ctx);
}
static __always_inline enum xdp_action Handle_vlan_tagged_packet(struct Packet_content *ctx) {
    return XDP_PASS;
}
static __always_inline __u16 Handle_downlink_packet(struct Packet_content *ctx, FLOW_CONTAINER *Packetflow_direction) {


    bpf_printk("Hexa: Handle_downlink_packet for ");
    //return XDP_PASS;
    struct PDR *PDR_content = Packetflow_direction->value;
if (!PDR_content){
        return XDP_PASS;
    }
    
    struct FAR *Far_content = bpf_map_lookup_elem(&FAR_map, &PDR_content->Far_id);
    if (!Far_content){
        return XDP_PASS;
    }


    switch (Far_content->Action) {
        case FAR_FORW:{
            struct QER *Qer_content = bpf_map_lookup_elem(&QER_map, &PDR_content->Qer_id);
            if (!Qer_content) {
                bpf_printk("Hexa: QER not found for TEID:%d Qer_ID:%d",Far_content->TEID, PDR_content->Qer_id);
                return XDP_DROP;
            }
            if (Qer_content->ULGate_status != GATE_OPEN)
                return XDP_DROP;
            if (Far_content->OHC == 1) {
                __u8 tos = Far_content->TLM >> 8;

                return Send_downlink_ipv4_packet(ctx, Far_content->Localip, Far_content->Remoteip, Qer_content->Qfi, tos,Far_content->TEID);            }

        }
        case FAR_DROP:
            return XDP_DROP;
        case FAR_BUFF:
            return Handle_buffer_action(ctx, Far_content, PDR_content);
        default:
            return XDP_DROP;
    }
     
}

static __always_inline enum xdp_action IPV4_PDU_handler(struct Packet_content *ctx ,FLOW_CONTAINER *Packetflow_direction) {

    __u32 TEID = ctx->gtp->TEID;
    bpf_printk("Hexa: Handle_uplink ip_packet for ");
    struct PDR *PDR_content = Packetflow_direction->value;
//return XDP_PASS;    


    
    struct FAR *Far_content = bpf_map_lookup_elem(&FAR_map, &PDR_content->Far_id);
    if (!Far_content) {
        bpf_printk("Hexa: FAR not found for TEID:%d Far_ID:%d", TEID, PDR_content->Far_id);
        return XDP_DROP;
    }
    bpf_printk("Hexa: FAR found for &ID:%d ", &PDR_content->Far_id);
    bpf_printk("Hexa: FAR found for ID:%d ", PDR_content->Far_id);
    bpf_printk("Hexa: FAR found for TEID:%d ", Far_content->TEID);
    bpf_printk("Hexa:  found for TEID:%d ", TEID);
    bpf_printk("Hexa: FAR found for &TEID:%d ", Far_content->TEID);
    bpf_printk("Hexa: FAR found for Action:%d ", Far_content->Action);
    bpf_printk("Hexa: FAR found for &Action:%d ", &Far_content->Action);
    bpf_printk("Hexa: FAR found for localip :%d ", Far_content->Localip);
    //return XDP_PASS;
    switch (Far_content->Action) {
        case (__u8)FAR_FORW:{
            bpf_printk("Hexa: FAR_FORW ");
            struct QER *Qer_content = bpf_map_lookup_elem(&QER_map, &PDR_content->Qer_id);
            if (!Qer_content) {
                bpf_printk("Hexa: QER not found for TEID:%d Qer_ID:%d", TEID, PDR_content->Qer_id);
                return XDP_DROP;
            }
            if (Qer_content->ULGate_status != (__u8)GATE_OPEN){
                return XDP_DROP;}
            //return XDP_PASS;
            if (PDR_content->OHR == (__u8)0) {
                //return XDP_PASS;
                int r = Remove_gtp_header(ctx);
                if (r != 1){
                    return XDP_ABORTED;}
                bpf_printk("Hexa: Remove gtp header outcome %d", r);
            }
            //return XDP_PASS;
            return Route_uplink_ipv4_packet(ctx);
            }
        case (__u8)FAR_DROP:{
            bpf_printk("Hexa: FAR_drop ");
            return XDP_PASS;
            return XDP_DROP;}
        case (__u8)FAR_BUFF:{
            bpf_printk("Hexa: FAR_buff ");
            return XDP_PASS;
            return Handle_buffer_action(ctx, Far_content, PDR_content);}
        default:{
            bpf_printk("Hexa: FAR_default ");
            return XDP_PASS;
            return XDP_DROP;}
    }
      
}

static __always_inline enum xdp_action IPV6_PDU_handler(struct Packet_content *ctx) {
    //TODO
    return XDP_PASS;
}
static __always_inline enum xdp_action Ethernet_PDU_handler(struct Packet_content *ctx) {
    //TODO
    return XDP_PASS;
}


static __always_inline enum xdp_action Handle_uplink_gtpu_packet(struct Packet_content *ctx, FLOW_CONTAINER *Packetflow_direction) {
   int SMF_pdutype = ctx->eth->h_proto;
       bpf_printk("Hexa: Handle_uplink sel_packet for 0x%x", SMF_pdutype);
   switch (bpf_ntohs(ctx->eth->h_proto)) {
       case ETH_P_IP:{
       bpf_printk("Hexa: Handle_uplink ip4pdu_packet for %d", SMF_pdutype);
           return IPV4_PDU_handler(ctx, Packetflow_direction);}
       case ETH_P_IPV6:{
       bpf_printk("Hexa: Handle_uplink ip6pdu_packet for %d", SMF_pdutype);
           return IPV6_PDU_handler(ctx);}
       case ETH_P_8021Q:{
       bpf_printk("Hexa: Handle_uplink vlanpdu_packet for %d", SMF_pdutype);
           return Ethernet_PDU_handler(ctx);}
       default:{
       bpf_printk("Hexa: Handle_uplink defaultpdu_packet for %d", SMF_pdutype);
           return XDP_DROP;}
   }
   
}





char _license[] SEC("license") = "GPL";
SEC("xdp/Hexa_datapath_entrypoint")
int Hexa_datapath_entrypoint(struct xdp_md *ctx) {

    struct Packet_content Packet_content = {
        .volatile_packet_ctx.data = ctx->data,
        .volatile_packet_ctx.data_end = ctx->data_end,
        .packet_context = ctx,};


    FLOW_CONTAINER Packetflow_direction =  Parse_packet_flow(&Packet_content);
    switch (Packetflow_direction.status) {
        case BAD_FLOW:{
            bpf_printk("Hexa: BAD_FLOW");
            return XDP_DROP;}
        case DOWNLINK_FLOW:{
            bpf_printk("Hexa: DOWNLINK_FLOW");
            //return XDP_PASS;
            return Handle_downlink_packet(&Packet_content, &Packetflow_direction);
            }
        case UPLINK_FLOW:{
            bpf_printk("Hexa: UPLINK_FLOW");
            return Handle_uplink_gtpu_packet(&Packet_content, &Packetflow_direction);}
        case UNHANDLED_FLOW:{
            bpf_printk("Hexa: UNHANDLED_FLOW");
            return XDP_PASS;}
        case NON_FLOW:
            bpf_printk("Hexa: NON_FLOW");
            return XDP_PASS;
    }
    
}
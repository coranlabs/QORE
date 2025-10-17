#pragma once

#include <bpf/bpf_helpers.h>
#include <linux/bpf.h>
#include <linux/types.h>



struct Checksum{
    __u64 Checksum_difference;
    __u16 Compressed_checksum;
};



static __always_inline __u16 Generate_ipv4_checksum(void *data_start, __u32 data_size) {
    __u64 Checksum_difference = bpf_csum_diff(0, 0, data_start, data_size, 0);
    struct Checksum csum = {
        .Checksum_difference = Checksum_difference
    };

    csum.Checksum_difference = (csum.Checksum_difference & 0xffffffff) + (csum.Checksum_difference >> 32);
    csum.Checksum_difference = (csum.Checksum_difference & 0xffff) + (csum.Checksum_difference >> 16);
    csum.Checksum_difference = (csum.Checksum_difference & 0xffff) + (csum.Checksum_difference >> 16);


    return ~csum.Checksum_difference;}



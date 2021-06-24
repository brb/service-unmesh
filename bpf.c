#include <string.h>
#include <linux/tcp.h>
#include <linux/bpf.h>
#include <netinet/in.h>
#include <bpf/bpf_helpers.h>

char _license[] SEC("license") = "GPL";

#ifndef SOL_TCP
#define SOL_TCP IPPROTO_TCP
#endif

#define SO_ORIGINAL_DST 80

SEC("cgroup/getsockopt")
int _getsockopt(struct bpf_sockopt *ctx)
{
	__u8 *optval_end = ctx->optval_end;
	__u8 *optval = ctx->optval;

	if (ctx->level == SOL_IP && ctx->optname == SO_ORIGINAL_DST) {
		if (optval + 10 > optval_end)
			return 0; /* EPERM, bounds check */
		optval[4] = 2;
		return 1;
	}

	return 1;
}

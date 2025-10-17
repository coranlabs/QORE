## Changing UPF Mode

This Kubernetes setup supports the User Plane Function (UPF) in two distinct modes:

- **`normal`** – Classical UPF mode  
- **`eupf`** – eBPF-based enhanced UPF mode

### How to Change the UPF Mode

To switch between UPF modes in your deployment:

1. Open the [`values.yaml`](../values.yaml) configuration file.
2. Locate the UPF configuration section.
3. Set the `mode:` field to either `normal` or `eupf`, depending on your desired mode.

```yaml
mode: eupf  # or normal
```

> [!CAUTION]
> Ensure that the `n6IPs` value is set to an IP address that is **not currently assigned** to any device in your network.  
> 
> For example:
> ```yaml
> n6IPs: "${NODE_IP_SUBNET}.7/24"
> ```
> In this case, the `.7` IP should be available. If `.7` is already in use, modify the last octet to another unassigned address.
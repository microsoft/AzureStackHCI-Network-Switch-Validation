# Switch Validation Troubleshooting Guide

### "Error opening interface, no suce device exists"

- Make sure the `hostInterfaceIP` value in `input.ini` file is being updated correctly.
- Make sure the host does configure the right IP on the connected NIC.
- Make sure the subnet mask is configured properly on both host and switch.

### "VLAN not match"

- Make sure `vlanIDs` value in `input.ini` file is being updated correctly.
- Make sure the switch configured VLAN match the input.
- Make sure the trunk interface connected with host configured properly.

### "Incorrect Maximum Frame Size"

- Make sure `mtuSize` value in `input.ini` file is being updated correctly.
- Make sure the switch configured the interface mtu match the input.

### "Incorrect ETS Maximum Number of Traffic Classes"

- ETS configuration not being detected via LLDP packet
- Make sure it is being configured in switch
- Make sure LLDP is enabled.

### "Incorrect ETS Class Bandwidth Configured"

- Make sure the switch ETS configuration match default validation ETS value: `0:48,1:0,2:0,3:50,4:0,5:2,6:0,7:0`
- [Not Recommend] Otherwise, uncomment the `ets` session in `input.ini`, and update customized value accordingly.

### "Incorrect PFC Maximum Number of Traffic Classes"

- PFC configuration not being detected via LLDP packet
- Make sure it is being configured in switch
- Make sure LLDP is enabled.

### "Incorrect PFC Priority Class Enabled"

- Make sure the switch PFC configuration match default validation PFC value: `0:0,1:0,2:0,3:1,4:0,5:0,6:0,7:0`
- [Not Recommend] Otherwise, uncomment the `pfc` session in `input.ini`, and update customized value accordingly.

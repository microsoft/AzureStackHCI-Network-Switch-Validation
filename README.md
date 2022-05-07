# Network Device Validation for Azure Stack

## Background

### Challenge

Physical network switches can pose a challenge to Azure Stack customers for several reasons including:

- No unified management plane due to 3P device.
- Inconsistent deployment methods and types based on customers' needs.

Thus, a standard/unified work stream will address physical network ethernet switch validation to:

- Enable informed purchasing (by customers) of ethernet switches.
- Ensure partners are selling devices that can function with the current Azure Stack requirements.
- Create leverage in the switch ecosystem, driving feature development by network vendors for future Azure Stack requirements.

## How to use the tool

The validation tool will collect network traffic and decode packages to validate protocol value required.

### Preparation

Here is the reference lab setup, but can be modified accordingly based on needs.
![Reference Lab Setup](./images/switchValidationLab01.png)

#### Set up Host Connection

- Prepare a host which has at least two NICs.
- The host can be a Server/Laptop with Windows/Linux system.
- Connect two NICs with the network switch interfaces.

#### Configure Network Switch

Here is the reference switch configuration based on [DellOS10](./switchReferenceConfig/Dell_OS10.conf)

Notice:

- Spanning Tree mode must be PVST for tool to capture all VLANID.
- LLDP must be enabled.

#### Download Validation Tool

- Download the right version based on the host OS.
  - [Windows Version](./switchValidationTool.exe)
  - [Linux Version](./switchValidationTool)
- Download [input.ini](input.ini) file.
- Update input variables accordingly.

```
C:\>switchValidationTool.exe -h
Usage of switchValidationTool.exe:
  -iniFilePath string
        Please input INI file path. (default "./input.ini")
```

### Execution

**Tool must be run with Administrator/Sudo privilege**

```
C:\>switchValidationTool.exe
2022/05/07 10:49:48 main.go:90: ./input.ini founded.
{10.10.10.11/24 [710 711 712] 9214 8 0:48,1:0,2:0,3:50,4:0,5:2,6:0,7:0 8 0:0,1:0,2:0,3:1,4:0,5:0,6:0,7:0}
2022/05/07 10:49:48 main.go:121: Found matched host interface by IP: 10.10.10.11/24 - \Device\NPF_{96CB802D-E41B-477E-BC46-B37A001AD1EF}
Processing, please wait up to ~2 mins, otherwise please double check if the interface has live traffic.
Collecting Network Packages: [1 / 300 (Max)]
Collecting Network Packages: [2 / 300 (Max)]
Collecting Network Packages: [3 / 300 (Max)]
...
Collecting Network Packages: [261 / 300 (Max)]
2022/05/07 10:51:21 packetCollect.go:61: Reach preset max session time 1m30s, close live collection.
2022/05/07 10:51:21 main.go:90: ./result.pcap founded.
Result PDF File Generated

### Result Summary ###

BGP - PASS
DHCPRelay - PASS
LLDP - PASS
VLAN - PASS
```

- To avoid endless running, the tool has preset maximum timeout condition, and will stop collecting whenever hit first.

  - 90 seconds
  - 300 network packets

- Please double check the interface connection and configuration if no network packet being collected.

### Post Execution

- Please check the result and re-test if need.
- Please upload these files to MSFT for further validation.
  - result.pdf
  - result.pcap
  - result.log

## Sample Result

Result report will be PDF, and check sample results under `sampleResult` folder.

## What will be validated in current version

### BGP

- TCP destination port 179

### DHCP Relay

- UDP destination port 67

### LLDP

- Subtype 1 (Native VLAN)
- Subtype 3 (All VLANs)
- Subtype 4 (MTU)
- Subtype 9 (ETS Configuration)
- Subtype B (PFC)
- Chassis sub type: MAC Address

### VLAN

- VLAN IDs allowed in the trunk

## Common Questions

### What should do if met errors while running the tool?

Please check [Troubleshooting_Manual](./Troubleshooting_Manual.md) to find matched error. If error not existing, please file issues to the repository.

### Host not able to run the tool or `alert security scan required`

Current version is still beta version, so hasn't signed, so that cause the alert, but it will be passed if running with `administrator` level.

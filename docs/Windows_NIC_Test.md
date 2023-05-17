# Windows NIC Test

## Background

The tool is written by Golang, and using [gopacket](https://pkg.go.dev/github.com/google/gopacket/pcap), but the network adapter/interface id/name is difference between Linux and Windows, so has to modified to fit the gopacket.

### Reference Links
[Golang and Windows Network Interfaces](https://haydz.github.io/2020/07/06/Go-Windows-NIC.html)
[net: get npcap usable windows network device names](https://github.com/golang/go/issues/35095#issuecomment-545528366%3E)

### Example
#### Linux
```golang
handle, err := pcap.OpenLive("eth0", 1600, true, pcap.BlockForever)
```

#### Windows
```golang
handle, err := pcap.OpenLive("\Device\NPF_{89A8C9DB-D0A4-4E79-8D7D-2FB8A578A221}", 1600, true, pcap.BlockForever)
```

The Interface GUID can be found by using Powershell
```powershell
PS C:\Test> Get-NetAdapter | Select-Object InterfaceAlias,InterfaceIndex,InterfaceGuid,DeviceName
    InterfaceAlias InterfaceIndex InterfaceGuid                          DeviceName
    -------------- -------------- -------------                          ----------
    NIC1                      18 {89A8C9DB-D0A4-4E79-8D7D-2FB8A578A221} \Device\{89A8C9DB-D0A4-4E79-8D7D-2FB8A578A221}

# Powershell interfaceGUID: "{89A8C9DB-D0A4-4E79-8D7D-2FB8A578A221}"
# Gopacket interface format: "\Device\NPF_{89A8C9DB-D0A4-4E79-8D7D-2FB8A578A221}"
```

## Test with Compiled file directly
Download the validation tool and open windows folder

```powershell
PS C:\windows> .\SwitchValidationTool.exe -h
Usage of C:\Users\liunick\Downloads\SwitchValidationTool_1.2305.77\windows\SwitchValidationTool.exe:
  -allVlanIDs string
        vlan list string separate with comma. Minimum 10 vlans required. (default "710,711,712,713,714,715,716,717,718,719,720")
  -etsBWbyPG string
        bandwidth for PGID in ETS configuration (default "0:48,1:0,2:0,3:50,4:0,5:2,6:0,7:0")
  -etsMaxClass int
        maximum number of traffic classes in ETS configuration (default 8)
  -interfaceAlias string
        Powershell: Get-NetAdapter | Select-Object InterfaceAlias,InterfaceGuid
  -interfaceGUID string
        Powershell: Get-NetAdapter | Select-Object InterfaceAlias,InterfaceGuid
  -mtu int
        mtu value configured on the switch interface (default 9214)
  -nativeVlanID int
        native vlan id (default 710)
  -pfcMaxClass int
        maximum PFC enabled traffic classes in PFC configuration (default 8)
  -pfcPriorityEnabled string
        PFC for priority in PFC configuration (default "0:0,1:0,2:0,3:1,4:0,5:0,6:0,7:0")

PS C:\windows> Get-NetAdapter | Select-Object InterfaceAlias,InterfaceIndex,InterfaceGuid,DeviceName
InterfaceAlias             InterfaceIndex InterfaceGuid                          DeviceName
--------------             -------------- -------------                          ----------
Ethernet                               18 {E6BEA924-DC17-48E1-97CA-A55C600C0976} \Device\{E6BEA924-DC17-48E1-97CA-A55C600C0976}

PS C:\windows> .\SwitchValidationTool.exe -interfaceGUID "\Device\NPF_{89A8C9DB-D0A4-4E79-8D7D-2FB8A578A221}" -interfaceAlias "NIC1" 
```



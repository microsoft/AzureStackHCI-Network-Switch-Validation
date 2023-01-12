#Requires -RunAsAdministrator
$ErrorActionPreference = "Stop"

function Invoke-SwitchValidation {
  <#
    .SYNOPSIS
    Execute Invoke-SwitchValidation on Windows

    .DESCRIPTION
    Azure Stack Hub HCI Switch Validation

    .EXAMPLE
    Invoke-SwitchValidation -ifIndex 15

    .EXAMPLE
    Invoke-SwitchValidation -ifIndex 15 -nativeVlanID 710 -allVlanIDs "710,711,712" -mtu 9216

    .EXAMPLE
    Invoke-SwitchValidation -ifIndex 15 -nativeVlanID 710 -allVlanIDs "710,711,712" -mtu 9216 -etsMaxClass 8 -etsBWbyPG "0:48,1:0,2:0,3:50,4:0,5:2,6:0,7:0" -pfcMaxClass 8 -pfcPriorityEnabled "0:0,1:0,2:0,3:1,4:0,5:0,6:0,7:0"

    .INPUTS
    $ifIndex: Host interface index based on "Get-NetAdapter". Mandatory.
    $nativeVlanID: Native Vlan ID
    $allVlanIDs: Vlan list string separate with comma. Default: "710,711,712".
    $mtu: MTU value configured on the switch interface. Default: 9214.
    $etsMaxClass: Maximum number of traffic classes in ETS configuration. Default: 8.
    $etsBWbyPG: Bandwidth for PGID in ETS configuration. Default: "0:48,1:0,2:0,3:50,4:0,5:2,6:0,7:0".
    $pfcMaxClass: Maximum PFC enabled traffic classes in PFC configuration. Default: 8.
    $pfcPriorityEnabled: PFC for priority in PFC configuration. Default: "0:0,1:0,2:0,3:1,4:0,5:0,6:0,7:0".

    .OUTPUTS
    XML file
#>
  [CmdLetBinding()]
  param (

    [parameter(Mandatory = $true)]
    [ValidateNotNullOrEmpty()]
    [UInt32]
    $ifIndex,

    [parameter(Mandatory = $true)]
    [ValidateNotNullOrEmpty()]
    [UInt32]
    $nativeVlanID,

    [parameter(Mandatory = $true)]
    [ValidateNotNullOrEmpty()]
    [string]
    $allVlanIDs,

    [parameter(Mandatory = $false)]
    [ValidateNotNullOrEmpty()]
    [UInt32]
    $mtu,

    [parameter(Mandatory = $false)]
    [ValidateNotNullOrEmpty()]
    [UInt32]
    $etsMaxClass,

    [parameter(Mandatory = $false)]
    [ValidateNotNullOrEmpty()]
    [string]
    $etsBWbyPG,

    [parameter(Mandatory = $false)]
    [ValidateNotNullOrEmpty()]
    [UInt32]
    $pfcMaxClass,

    [parameter(Mandatory = $false)]
    [ValidateNotNullOrEmpty()]
    [string]
    $pfcPriorityEnabled
  )

  try {
    $exefile = ".\SwitchValidationTool.exe"
    $intfList = Get-NetAdapter | Select-Object InterfaceAlias, InterfaceGuid, ifIndex
  
    foreach ($intf in $intfList) {
      if ($intf.ifIndex -eq $ifIndex ) {
        $inftAlias = $intf.InterfaceAlias -replace " ",""
        $inftGUID = $intf.InterfaceGuid
        # Windows NIC name is different from Linux, which need to be converted before use in gopacket lib.
        # PS C:\Users\liunick\Downloads\Test> Get-NetAdapter | Select-Object InterfaceAlias,InterfaceIndex,InterfaceGuid,DeviceName
        # InterfaceAlias InterfaceIndex InterfaceGuid                          DeviceName
        # -------------- -------------- -------------                          ----------
        # NIC1                      18 {A91A8E1F-C8B3-4D96-A403-78B9E758EA38} \Device\{A91A8E1F-C8B3-4D96-A403-78B9E758EA38}
        #####
        # Powershell interfaceGUID: "{A91A8E1F-C8B3-4D96-A403-78B9E758EA38}"
        # Gopacket interface format: "\Device\NPF_{A91A8E1F-C8B3-4D96-A403-78B9E758EA38}"
        $interfaceGUID="Device\NPF_$($inftGUID)"
        write-host "## Interface Name $inftAlias is selected, [Debug InterfaceGUID: $interfaceGUID] ##"
      } 
    }
    
    if ($interfaceGUID -ne "" -and $inftAlias -ne "") {
      $arguments += "-interfaceAlias `"$inftAlias`" -interfaceGUID `"$($interfaceGUID)`""
      if ($nativeVlanID -ne 0) {
        $arguments += " -nativeVlanID `"$nativeVlanID`""
      }
      if ($allVlanIDs -ne "") {
        $arguments += " -allVlanIDs `"$allVlanIDs`""
      }
      if ($mtu -ne 0) {
        $arguments += " -mtu $mtu"
      }
      if ($etsMaxClass -ne 0) {
        $arguments += " -etsMaxClass $etsMaxClass"
      }
      if ($etsBWbyPG -ne "") {
        $arguments += " -etsBWbyPG `"$etsBWbyPG`""
      }
      if ($pfcMaxClass -ne 0) {
        $arguments += " -pfcMaxClass $pfcMaxClass"
      }
      if ($pfcPriorityEnabled -ne "") {
        $arguments += " -pfcPriorityEnabled `"$pfcPriorityEnabled`""
      }
      Start-Process -NoNewWindow -FilePath $exefile -ArgumentList $arguments
    }
    else {
      write-host "no interface founded"
    }
  }
  catch {
    write-host "Please use 'Get-Help' to check full instruction of the function."
  }
}
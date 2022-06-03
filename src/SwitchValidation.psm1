#Requires -RunAsAdministrator
$ErrorActionPreference = "Stop"

function Invoke-SwitchValidation {
      <#
    .SYNOPSIS
    Execute Invoke-SwitchValidation on Windows

    .DESCRIPTION
    Azure Stack Hub HCI Switch Validation

    .EXAMPLE
    Invoke-SwitchValidation -ifIndex 12

    .INPUTS
    $ifIndex: Host interface index based on "Get-NetAdapter"

    .OUTPUTS
    XML file
#>
  [CmdLetBinding()]
  param (
    # Local encrypted XML - Switch SSH Credential
    [parameter(Mandatory = $true)]
    [ValidateNotNullOrEmpty()]
    [UInt32]
    $ifIndex
  )
  $exefile=".\SwitchValidationTool.exe"
  $intfList= Get-NetAdapter | Select-Object InterfaceAlias,InterfaceGuid,ifIndex

  foreach ($intf in $intfList)
  {
    if ($intf.ifIndex -eq $ifIndex ) {
      $inftAlias=$intf.InterfaceAlias
      $inftGUID=$intf.InterfaceGuid
      write-host "interface $($intf.InterfaceAlias) is selected"
      write-host "$($intf.InterfaceGuid)"
    } 
  }
  
  if ($inftGUID -ne "" -and $inftAlias -ne "") {
      $arguments="-interfaceAlias '$inftAlias' -interfaceGUID '$inftGUID'"
      write-host $arguments
      & $exefile -interfaceAlias $inftAlias -interfaceGUID $inftGUID
  }else{
      write-host "no interface founded"
  }
}
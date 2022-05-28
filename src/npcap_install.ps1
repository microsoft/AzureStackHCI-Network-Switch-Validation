function Check_Program_Installed( $programName ) {
    $x86_check = ((Get-ChildItem "HKLM:Software\Microsoft\Windows\CurrentVersion\Uninstall") |
    Where-Object { $_."Name" -like "*$programName*" } ).Length -gt 0;
      
    if(Test-Path 'HKLM:Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall')  
    {
    $x64_check = ((Get-ChildItem "HKLM:Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall") |
    Where-Object { $_."Name" -like "*$programName*" } ).Length -gt 0;
    }
    return $x86_check -or $x64_check;
    }

function Install_Npcap {
    $urlPath = "https://nmap.org/npcap/dist/npcap-0.96.exe"
    $checksum = "83667e1306fdcf7f9967c10277b36b87e50ee8812e1ee2bb9443bdd065dc04a1"

    # Download the file
    echo "Downloading... ($urlPath)"
    Invoke-WebRequest $urlPath -UseBasicParsing -OutFile $PSScriptRoot"\npcap.exe"

    # Now let's check its checksum
    $_chksum = $(CertUtil -hashfile $PSScriptRoot"\npcap.exe" SHA256)[1] -replace " ", ""
    if ($_chksum -ne $checksum) {
        echo "Hashes NOT match."
        exit
    }
    else {
        echo "Hashes match."
        echo "Installing Npcap..."
    }
    # Run installer
    Start-Process $PSScriptRoot"\npcap.exe" -ArgumentList "/S /npf_startup=yes /loopback_support=yes /dlt_null=no /admin_only=no /dot11_support=no /vlan_support=no /winpcap_mode=yes" -wait
    echo "Npcap has been installed."
}

$npcap_installed=Check_Program_Installed("Npcap")

if ($npcap_installed) {
    echo "Npcap has been installed."
} else {
    Install_Npcap
}
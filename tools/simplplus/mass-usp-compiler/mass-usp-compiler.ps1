param (
    [Parameter(Mandatory=$true)]
    [string]$Path,
    [switch]$series2,
    [switch]$series3
)
$Source = Split-Path -Path $PSScriptRoot -Parent
$FolderPath = Split-Path -Path $Source -Parent
$FileExtension = ".usp"

$Dir = get-childitem $Path -recurse
$List = $Dir | where {$_.extension -eq $FileExtension}

$List | ForEach-Object {
    if($series2.IsPresent -and $series3.IsPresent)
    {
        & "C:\Program Files (x86)\Crestron\Simpl\SPlusCC.exe" \build $_.FullName \target series2 series3
    }
    elseif ($series2.IsPresent -and !($series3.IsPresent))
    {
        & "C:\Program Files (x86)\Crestron\Simpl\SPlusCC.exe" \build $_.FullName \target series2
    }
    elseif ($series3.IsPresent -and !($series2.IsPresent))
    {
        
        & "C:\Program Files (x86)\Crestron\Simpl\SPlusCC.exe" \build $_.FullName \target series3
    }
    else
    {
        
        & "C:\Program Files (x86)\Crestron\Simpl\SPlusCC.exe" \build $_.FullName \target series2
    }
}

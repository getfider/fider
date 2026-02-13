# PowerShell script to add test domains to Windows hosts file
# Run as Administrator

$hostsFile = "$env:SystemRoot\System32\drivers\etc\hosts"
$domains = @(
    "127.0.0.1 fider.local",
    "127.0.0.1 app.local",
    "127.0.0.1 multi.local",
    "127.0.0.1 tenant1.multi.local",
    "127.0.0.1 tenant2.multi.local"
)

Write-Host "Adding test domains to hosts file..." -ForegroundColor Green

# Check if running as admin
$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
if (-not $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Host "ERROR: This script must be run as Administrator" -ForegroundColor Red
    Write-Host "Right-click PowerShell and select 'Run as Administrator'" -ForegroundColor Yellow
    exit 1
}

# Read current hosts file
$hostsContent = Get-Content $hostsFile

# Add domains if not already present
$modified = $false
foreach ($domain in $domains) {
    if ($hostsContent -notcontains $domain) {
        Add-Content -Path $hostsFile -Value $domain
        Write-Host "Added: $domain" -ForegroundColor Cyan
        $modified = $true
    } else {
        Write-Host "Already exists: $domain" -ForegroundColor Gray
    }
}

if ($modified) {
    Write-Host "`nHosts file updated successfully!" -ForegroundColor Green
    Write-Host "Location: $hostsFile" -ForegroundColor Gray
} else {
    Write-Host "`nAll domains already configured." -ForegroundColor Green
}

Write-Host "`nYou can now run: docker-compose -f docker-compose-test.yml up -d" -ForegroundColor Yellow

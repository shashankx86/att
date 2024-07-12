# Function to download a file
function Download-File {
    param (
        [string]$url,
        [string]$outputPath
    )
    Write-Output "Downloading $url..."
    Invoke-WebRequest -Uri $url -OutFile $outputPath
    Write-Output "Downloaded to $outputPath"
}

# Function to determine the platform and architecture
function Get-PlatformArch {
    $platform = "windows"
    $arch = (Get-WmiObject Win32_Processor).Architecture
    switch ($arch) {
        9 { $arch = "amd64" }
        5 { $arch = "arm64" }
        default { Write-Output "Unsupported architecture: $arch"; exit 1 }
    }
    return "$platform-$arch"
}

# GitHub repository details
$repoOwner = "shashankx86"
$repoName = "att"

# Get the latest release from GitHub API
$latestReleaseUrl = "https://api.github.com/repos/$repoOwner/$repoName/releases/latest"
Write-Output "Fetching latest release info from $latestReleaseUrl..."
$latestReleaseJson = Invoke-RestMethod -Uri $latestReleaseUrl

# Extract the tag name and assets URL
$tagName = $latestReleaseJson.tag_name
$assetUrls = $latestReleaseJson.assets | ForEach-Object { $_.browser_download_url }

Write-Output "Latest release: $tagName"

# Get the platform and architecture
$platformArch = Get-PlatformArch

# Find the appropriate asset for the current platform and architecture
$assetUrl = $assetUrls | Where-Object { $_ -match "att-$platformArch" }

if (-not $assetUrl) {
    Write-Output "No matching asset found for $platformArch"
    exit 1
}

# Download the asset to temp path
$tmpDownloadPath = "$env:TEMP\att.exe"
Download-File -url $assetUrl -outputPath $tmpDownloadPath

# Move the file to a directory in the system PATH
$downloadDir = "C:\Program Files\att"
if (-not (Test-Path $downloadDir)) {
    New-Item -ItemType Directory -Path $downloadDir
}
Move-Item -Path $tmpDownloadPath -Destination "$downloadDir\att.exe"

# Add the directory to the system PATH
$env:PATH += ";$downloadDir"
[Environment]::SetEnvironmentVariable("PATH", $env:PATH, [EnvironmentVariableTarget]::Machine)

Write-Output "att installed successfully in $downloadDir"
Write-Output "You can run it using the command: att"

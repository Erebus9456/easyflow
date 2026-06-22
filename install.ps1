Write-Host "🚀 Starting EasyFlow Installer (Windows Powershell)..." -ForegroundColor Cyan
Write-Host "----------------------------------------"

# 1. Verify / Install Git & GitHub CLI via Winget package manager
function Enforce-Dependency ($Command, $PackageId) {
    if (Get-Command $Command -ErrorAction SilentlyContinue) {
        Write-Host "✅ $Command is already installed." -ForegroundColor Green
    } else {
        Write-Host "📦 $Command not found. Installing package via Winget ($PackageId)..." -ForegroundColor Yellow
        winget install --id $PackageId --silent --accept-source-agreements --accept-package-agreements
        # Refresh environment block inside active process path scope
        $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
    }
}

Enforce-Dependency "git" "Git.Git"
Enforce-Dependency "gh" "GitHub.cli"

# 2. Check and evaluate Go installation status
$UpdateGo = $false
if (Get-Command go -ErrorAction SilentlyContinue) {
    $GoVersionString = (go version) -split " " | Select-Object -Index 2
    Write-Host "🔍 Found installed Go version: $GoVersionString" -ForegroundColor Gray
    
    # Simple semantic major verification matching format checks
    if ($GoVersionString -match "go1\.(1[0-9]|2[0])\.") {
        Write-Host "⚠️ Go runtime version is outdated. Upgrading required..." -ForegroundColor Yellow
        $UpdateGo = $true
    }
} else {
    Write-Host "⚠️ Go compiler runtime not found on system." -ForegroundColor Yellow
    $UpdateGo = $true
}

# 3. Upgrade Go via Winget if missing or ancient
if ($UpdateGo) {
    Write-Host "📥 Installing stable Go Programming Language workspace..." -ForegroundColor Yellow
    winget install --id GoLang.Go --silent --accept-source-agreements --accept-package-agreements
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
} else {
    Write-Host "✅ Go toolchain satisfies running application dependencies." -ForegroundColor Green
}

# 4. Clear caching modules footprints and execute global tool pull string
Write-Host "🧼 Cleaning module cache footprints..." -ForegroundColor Gray
go clean -modcache

Write-Host "🏎️ Compiling and installing EasyFlow v1.0.2..." -ForegroundColor Cyan
$env:GOPROXY="direct"
go install github.com/Erebus9456/easyflow@v1.0.2

# 5. Inject Go Bin paths explicitly into User Environment Variable Blocks persistently
$UserPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
$GoBinPath = "$env:USERPROFILE\go\bin"
if ($UserPath -notlike "*$GoBinPath*") {
    Write-Host "⚙️ Attaching Go runtime target paths to User Environment Matrix variables..." -ForegroundColor Gray
    [System.Environment]::SetEnvironmentVariable("Path", $UserPath + ";" + $GoBinPath, "User")
    $env:Path += ";$GoBinPath"
}

# 6. Verify GitHub CLI (gh) Identity Authentication States
Write-Host "----------------------------------------"
Write-Host "🔐 Verifying GitHub CLI authentication context..." -ForegroundColor Cyan
gh auth status 2>$null
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ GitHub CLI is authenticated and connected properly!" -ForegroundColor Green
} else {
    Write-Host "⚠️ GitHub CLI credentials token mismatch or missing." -ForegroundColor Yellow
    Write-Host "🔄 Initializing secure login flow sequence tool prompt..." -ForegroundColor Cyan
    gh auth login
}

Write-Host "----------------------------------------"
Write-Host "🎉 EasyFlow installation completed successfully!" -ForegroundColor Green
Write-Host "💡 Please open a NEW terminal screen to refresh system environment execution matrices." -ForegroundColor Yellow
Write-Host "🚀 Type 'easyflow' to begin execution!" -ForegroundColor Green
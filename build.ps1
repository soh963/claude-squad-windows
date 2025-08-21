# Build script for Claude Squad Windows version

param(
    [string]$Version = "1.0.12-windows",
    [string]$Output = ".\bin\claude-squad-windows.exe",
    [switch]$Test = $false
)

Write-Host "Building Claude Squad Windows Edition..." -ForegroundColor Green

# Ensure bin directory exists
$binDir = Split-Path -Parent $Output
if (-not (Test-Path $binDir)) {
    New-Item -ItemType Directory -Path $binDir -Force | Out-Null
}

# Copy necessary Go files from original project
$sourceFiles = @(
    "app\app.go",
    "cmd\cmd.go", 
    "config\state.go",
    "daemon\daemon.go",
    "log\log.go",
    "session\git\git.go",
    "session\session.go"
)

Write-Host "Copying source files..." -ForegroundColor Yellow

foreach ($file in $sourceFiles) {
    $sourcePath = "..\claude-squad\$file"
    $destPath = ".\$file"
    $destDir = Split-Path -Parent $destPath
    
    if (-not (Test-Path $destDir)) {
        New-Item -ItemType Directory -Path $destDir -Force | Out-Null
    }
    
    if (Test-Path $sourcePath) {
        Copy-Item $sourcePath $destPath -Force
        Write-Host "  Copied: $file" -ForegroundColor Gray
    } else {
        Write-Host "  Warning: Source file not found: $sourcePath" -ForegroundColor Yellow
    }
}

# Initialize go module
Write-Host "Initializing Go module..." -ForegroundColor Yellow
& go mod init claude-squad-windows 2>$null
& go mod tidy

# Set build variables
$env:CGO_ENABLED = "0"
$env:GOOS = "windows"
$env:GOARCH = "amd64"

# Build flags
$buildFlags = @(
    "-ldflags", "-X main.version=$Version -w -s",
    "-o", $Output,
    "."
)

if ($Test) {
    Write-Host "Running tests..." -ForegroundColor Yellow
    & go test ./...
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Tests failed!" -ForegroundColor Red
        exit 1
    }
    Write-Host "Tests passed!" -ForegroundColor Green
}

Write-Host "Building executable..." -ForegroundColor Yellow
& go build @buildFlags

if ($LASTEXITCODE -eq 0) {
    $fileInfo = Get-Item $Output
    $sizeKB = [math]::Round($fileInfo.Length / 1KB, 2)
    
    Write-Host "Build successful!" -ForegroundColor Green
    Write-Host "Output: $Output" -ForegroundColor White
    Write-Host "Size: $sizeKB KB" -ForegroundColor White
    Write-Host "Version: $Version" -ForegroundColor White
    
    Write-Host "`nTo test the build:" -ForegroundColor Cyan
    Write-Host "  $Output version" -ForegroundColor Gray
    Write-Host "  $Output debug" -ForegroundColor Gray
} else {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}
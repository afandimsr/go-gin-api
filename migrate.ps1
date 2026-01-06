param (
    [Parameter(Mandatory = $true)]
    [ValidateSet("up", "down", "force")]
    [string]$Action,

    [Parameter(Mandatory = $false)]
    [string]$Version,

    [Parameter(Mandatory = $false)]
    [switch]$ConfirmInput
)

# Load .env file
if (Test-Path .env) {
    Get-Content .env | Where-Object { $_ -match '=' } | ForEach-Object {
        $key, $value = $_.Split('=', 2)
        [Environment]::SetEnvironmentVariable($key.Trim(), $value.Trim(), "Process")
    }
}
else {
    Write-Error ".env file not found"
    exit 1
}

# Construct Database URL
$DB_USER = $env:DB_USER
$DB_PASS = $env:DB_PASS
$DB_HOST = $env:DB_HOST
$DB_PORT = $env:DB_PORT
$DB_NAME = $env:DB_NAME
$APP_ENV = $env:APP_ENV

$MIGRATE_URL = "mysql://${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}"

# Safeguard against accidental production rollback
if ($Action -eq "down" -and $APP_ENV -eq "production" -and -not $ConfirmInput) {
    Write-Host "ERROR: DANGER! You are attempting to rollback in PRODUCTION." -ForegroundColor Red
    Write-Host "This can cause DATA LOSS. To proceed, use: .\migrate.ps1 -Action down -ConfirmInput" -ForegroundColor Yellow
    exit 1
}

# Execute migrate command
$Cmd = "migrate -database ""$MIGRATE_URL"" -path migrations"

if ($Action -eq "force") {
    if (-not $Version) {
        Write-Error "Version parameter is required for force action"
        exit 1
    }
    Invoke-Expression "$Cmd force $Version"
}
else {
    Invoke-Expression "$Cmd $Action"
}

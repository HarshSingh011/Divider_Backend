Get-ChildItem -Path "." -Filter "*.go" -Recurse | ForEach-Object {
    $file = $_
    $content = Get-Content $file.FullName -Raw
    $cleaned = $content -replace '^\s*//.*$', ''
    $cleaned = $cleaned -replace '`n`n`n+', "`n`n"
    Set-Content -Path $file.FullName -Value $cleaned -Encoding UTF8
    Write-Host "Cleaned: $($file.Name)"
}

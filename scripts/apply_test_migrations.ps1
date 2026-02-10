$files = Get-ChildItem "migrations\*.up.sql" | Sort-Object Name
foreach ($file in $files) {
    Write-Host "Applying migration: $($file.Name) to omnichat_test"
    Get-Content $file.FullName | docker exec -i omnichat-postgres psql -U omnichat -d omnichat_test
}

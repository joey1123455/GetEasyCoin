# Hardcoded list of directories
$directories = @(
    "./handlers/",
    "./services/",
    "./utils/",
    "./middleware/"
)

# Loop through the hardcoded list of directories
foreach ($dir in $directories) {
    Write-Host "Testing $dir"
    go test $dir
}

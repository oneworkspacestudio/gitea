param([string]$Go='go')
$ErrorActionPreference='Stop'
Push-Location $PSScriptRoot
try {
 New-Item -ItemType Directory -Force payload,dist | Out-Null
 $assets=@(
  @{name='gitea.exe';url='https://github.com/go-gitea/gitea/releases/download/v28.1.0/gitea-28.1.0-windows-amd64.exe';sha='68495c963bab1bfcbbe5b48916c7f362a49aca931583575e09763a02ddf63e2d'},
  @{name='mingit.zip';url='https://github.com/git-for-windows/git/releases/download/v2.56.0.windows.2/MinGit-2.56.0.2-64-bit.zip';sha='da35e72aa21c005a5a0d298cfbae110bc1609a815730ea0dde84b01a1b3cd3be'}
 )
 foreach($asset in $assets){$path=Join-Path 'payload' $asset.name;if(!(Test-Path -LiteralPath $path)){Invoke-WebRequest $asset.url -OutFile $path};if((Get-FileHash $path -Algorithm SHA256).Hash.ToLowerInvariant() -ne $asset.sha){throw "Payload checksum mismatch: $($asset.name)"}}
 Add-Type -AssemblyName System.IO.Compression
 $stream=[IO.File]::Create((Join-Path $PSScriptRoot 'payload.zip'));$zip=[IO.Compression.ZipArchive]::new($stream,[IO.Compression.ZipArchiveMode]::Create)
 try {
  $entry=$zip.CreateEntry('gitea.exe');$output=$entry.Open();$input=[IO.File]::OpenRead((Join-Path $PSScriptRoot 'payload/gitea.exe'));try{$input.CopyTo($output)}finally{$input.Dispose();$output.Dispose()}
  $git=[IO.Compression.ZipFile]::OpenRead((Join-Path $PSScriptRoot 'payload/mingit.zip'));try{foreach($file in $git.Entries){if(!$file.Name){continue};$entry=$zip.CreateEntry('git/'+$file.FullName);$input=$file.Open();$output=$entry.Open();try{$input.CopyTo($output)}finally{$input.Dispose();$output.Dispose()}}}finally{$git.Dispose()}
 }finally{$zip.Dispose();$stream.Dispose()}
 & $Go build -buildvcs=false -trimpath -ldflags '-s -w' -o dist/Gitea.exe .
 if($LASTEXITCODE){throw 'Gitea Workspace launcher build failed'}
 Get-FileHash dist/Gitea.exe -Algorithm SHA256
}finally{Pop-Location}

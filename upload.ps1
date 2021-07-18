$remotePath = "/www/history";
$composeFilePath = "~/dev/"
$fileList = "./dist/*"
foreach ($item in $fileList) {
    scp -r $item server2:$remotePath;
}
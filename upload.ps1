$remotePath = "/www/history";
$composeFilePath = "~/dev/"
$fileList = "./dist/*"
yarn build
foreach ($item in $fileList) {
    scp -r $item server2:$remotePath;
}
ssh server2 docker exec nginx nginx -s reload
$remotePath = "/www/history";
$composeFilePath = "~/dev/"
$fileList = "./dist/*"
yarn build
foreach ($item in $fileList) {
    scp -r $item server:$remotePath;
}
ssh server docker exec nginx nginx -s reload
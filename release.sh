rm -r release
mkdir release
go build -o=release/server
npm run build
tar -cvzf release.tar.gz release/
gsutil cp release.tar.gz gs://cdn.tcole.me
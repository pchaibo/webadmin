# 
git fetch --all
git reset --hard origin/main
cd ./Backend
rm -rf webadmin
go build .
echo "ok"
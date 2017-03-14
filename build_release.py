import os

if not os.path.exists("release"):
    os.mkdir("release")

os.system("rm ./error_*.log")
if os.system("GOOS=linux GOARCH=amd64 go build -o ./release/nginx_log_analy main.go") > 0:
    raise ValueError("build error")

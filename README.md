# myclone
a reddit clone test

## Without docker
use
```
go get github.com/thinkong/myclone
```

```
go install github.com/thinkong/myclone/main
```
copy `templates` folder to bin folder

then run with myclone.exe

## With Docker
I have added a Docker file.

So now you can clone the repo
```
git clone https://github.com/thinkong/myclone
```

If you have docker installed you can run

```
cd myclone
docker build -t myclone .
docker run --publish 8080:8080 --name test --rm myclone
```

after docker is running, you can go to http://localhost:8080


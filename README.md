```
# 도커 테스트 (로컬)
docker build -t test-server:local .
docker run --rm -p 8080:8080 test-server:local
```

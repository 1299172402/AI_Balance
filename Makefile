build-docker-in-china:
	docker build \
	--build-arg APK_MIRROR=mirrors.tuna.tsinghua.edu.cn \
	-t ai-balance .

build-docker:
	docker build -t ai-balance .

run-docker:
	docker run --rm -p 8080:8080 --env-file .env ai-balance

run-native:
	bash -c 'set -a && source .env && go run -tags dev ./cmd/ai-balance'


build-dev:
	docker build -t go-webrtc -f containers/images/Dockerfile . && docker build -t turn -f containers/images/Dockerfile .

clean-dev:
	docker-compose -f containers/compse/dc.dev.yml down

run-dev:
	docker-compose -f containers/compse/dc.dev.yml up

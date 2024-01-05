.PHONY: proto_go

GENERATE_DIR := generated
BUILD_TARGET := ./build
ANDROID_SDK_DIR=./android

proto_go:
	rm -rf ${GENERATE_DIR}/*
	mkdir -p $(GENERATE_DIR)
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
	protoc proto/*.proto -Iproto/ \
		--go_out=$(GENERATE_DIR) --go_opt=paths=import \
		--experimental_allow_proto3_optional


proto_oc:
	rm -rf $(BUILD_TARGET)/build_ios/proto
	mkdir -p $(BUILD_TARGET)/build_ios/proto
	protoc proto/*.proto -Iproto \
		--swift_opt=Visibility=Public \
		--swift_out=$(BUILD_TARGET)/build_ios/proto

proto_kt:
	cd  ${ANDROID_SDK_DIR}/ && ./gradlew -Dhttp.proxyHost :bean:clean && ./gradlew -Dhttp.proxyHost :bean:assembleRelease

ios:proto_go proto_oc
	rm -rf $(BUILD_TARGET)/build_ios/BtcApi.xcframework
	mkdir -p $(BUILD_TARGET)/build_ios
	go get golang.org/x/mobile
	go mod download golang.org/x/exp
	GOARCH=arm64 gomobile bind -v -trimpath -ldflags "-s -w" \
 	-o ${BUILD_TARGET}/build_ios/Wallet.xcframework -target=ios ./api

android:proto_go
	mkdir -p $(BUILD_TARGET)/build_android
	rm -rf $(BUILD_TARGET)/build_android/*
	go get golang.org/x/mobile
	go mod download golang.org/x/exp
	time GOARCH=arm64 gomobile bind -v -trimpath -ldflags "-s -w" \
	-o ${BUILD_TARGET}/build_android/Wallet.aar -target=android ./api
	unzip -d $(BUILD_TARGET)/build_android/sources $(BUILD_TARGET)/build_android/Wallet-sources.jar
GOPKG ?=	moul.io/pkgman
DOCKER_IMAGE ?=	moul/pkgman
GOBINS ?=	.

PRE_UNITTEST_STEPS = testdata/gomobile-ipfs-demo-1.1.0.ipa
PRE_TEST_STEPS = $(PRE_UNITTEST_STEPS)
testdata/gomobile-ipfs-demo-1.1.0.ipa:
	mkdir -p testdata
	wget -N -O $@ "https://bintray.com/berty/gomobile-ipfs-demo/download_file?file_path=ios%2F1.1.0%2Fios-demo-1.1.0.ipa"

include rules.mk

package ipa

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func ExampleOpen() {
	pkg, err := Open("path/to/ipa")
	if err != nil {
		panic(err)
	}
	defer pkg.Close()

	fmt.Printf("Loaded IPA: %v\n", pkg)
}

func TestOpen_gomobileipfsdemo(t *testing.T) {
	testdata := "../../testdata/gomobile-ipfs-demo-1.1.0.ipa"
	pkg, err := Open("../../testdata/gomobile-ipfs-demo-1.1.0.ipa")
	if err != nil {
		t.Skipf("cannot open testdata %q, skipping", testdata)
	}
	require.NoError(t, err)
	defer pkg.Close()
	require.Len(t, pkg.Apps(), 1)
	require.Len(t, pkg.Files(), 49)

	app := pkg.Apps()[0]
	require.Equal(t, app.Name, "GomobileIPFS Example.app")
	require.Len(t, app.Files(), 48)
	plist, err := app.Plist()
	require.NoError(t, err)
	expected := &Plist{
		BuildMachineOSBuild:              "19D76",
		CFBundleDevelopmentRegion:        "en",
		CFBundleExecutable:               "GomobileIPFS Example",
		CFBundleIdentifier:               "ipfs.gomobile.example",
		CFBundleInfoDictionaryVersion:    "6.0",
		CFBundleName:                     "GomobileIPFS Example",
		CFBundlePackageType:              "APPL",
		CFBundleShortVersionString:       "0.0.1",
		CFBundleSignature:                "????",
		CFBundleSupportedPlatforms:       []string{"iPhoneOS"},
		CFBundleVersion:                  "1",
		DTCompiler:                       "com.apple.compilers.llvm.clang.1_0",
		DTPlatformBuild:                  "17B102",
		DTPlatformName:                   "iphoneos",
		DTPlatformVersion:                "13.2",
		DTSDKBuild:                       "17B102",
		DTSDKName:                        "iphoneos13.2",
		DTXcode:                          "1130",
		DTXcodeBuild:                     "11C504",
		LSRequiresIPhoneOS:               true,
		MinimumOSVersion:                 "10.0",
		UIDeviceFamily:                   []int{1},
		UILaunchStoryboardName:           "LaunchScreen",
		UIMainStoryboardFile:             "Main",
		UIRequiredDeviceCapabilities:     []string{"armv7"},
		UISupportedInterfaceOrientations: []string{"UIInterfaceOrientationPortrait", "UIInterfaceOrientationLandscapeLeft"},
	}
	require.Equal(t, plist, expected)
}

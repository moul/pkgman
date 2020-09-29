package apk_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"moul.io/pkgman/pkg/apk"
)

func ExampleOpen() {
	pkg, err := apk.Open("path/to/apk")
	if err != nil {
		panic(err)
	}
	defer pkg.Close()

	fmt.Printf("Loaded APK: %v\n", pkg)
}

func TestOpen_gomobileipfsdemo(t *testing.T) {
	testdata := "../../testdata/gomobile-ipfs-demo-1.2.1.apk"
	pkg, err := apk.Open(testdata)
	if err != nil {
		t.Skipf("cannot open testdata %q, skipping", testdata)
	}
	require.NoError(t, err)
	defer pkg.Close()
	require.Len(t, pkg.Files(), 451)

	manifest, err := pkg.Manifest()
	require.NoError(t, err)
	require.Equal(t, manifest.Package, "null.example")
	require.Equal(t, manifest.MainActivity().Label, "Gomobile IPFS Example")
	require.Equal(t, manifest.MainActivity().Name, "ipfs.gomobile.example.MainActivity")

	noExists, err := pkg.FileBytes("blahblah")
	require.Nil(t, noExists)
	require.Error(t, err)

	exists, err := pkg.FileBytes("resources.arsc")
	require.NoError(t, err)
	require.Equal(t, len(exists), 257156)
}

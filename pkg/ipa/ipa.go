package ipa

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"howett.net/plist"
)

type Package struct {
	r *zip.ReadCloser
}

func Open(path string) (*Package, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("open path: %w", err)
	}
	return &Package{r: r}, nil
}

func (p *Package) Close() error {
	return p.r.Close()
}

func (p *Package) Apps() []*App {
	apps := []*App{}
	for _, f := range p.r.File {
		match, err := path.Match("Payload/*.app/Info.plist", f.Name)
		if match && err == nil {
			name := path.Base(path.Dir(f.Name))
			apps = append(apps, &App{
				p:    p,
				Name: name,
			})
		}
	}
	return apps
}

func (p Package) Files() []*zip.File {
	return p.r.File
}

func (p Package) File(name string) *zip.File {
	for _, f := range p.r.File {
		if f.Name == name {
			return f
		}
	}
	return nil
}

type App struct {
	p    *Package
	Name string
}

func (a App) dir() string {
	return fmt.Sprintf("Payload/%s", a.Name)
}

func (a App) Files() []*zip.File {
	files := []*zip.File{}
	prefix := a.dir()
	for _, f := range a.p.r.File {
		if strings.HasPrefix(f.Name, prefix) {
			files = append(files, f)
		}
	}
	return files
}

func (a App) File(name string) *zip.File {
	target := path.Join(a.dir(), name)
	for _, f := range a.Files() {
		if f.Name == target {
			return f
		}
	}
	return nil
}

func (a App) Plist() (*Plist, error) {
	f := a.File("Info.plist")
	if f == nil {
		return nil, fmt.Errorf("no such file: Info.plist")
	}

	r, err := f.Open()
	if err != nil {
		return nil, err
	}

	bin, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(bin)

	var ret Plist
	decoder := plist.NewDecoder(buf)
	err = decoder.Decode(&ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

type Plist struct {
	BuildMachineOSBuild       string `json:"BuildMachineOSBuild"`
	CFBundleDevelopmentRegion string `json:"CFBundleDevelopmentRegion"`
	CFBundleDisplayName       string `json:"CFBundleDisplayName"`
	CFBundleExecutable        string `json:"CFBundleExecutable"`
	CFBundleIcons             struct {
		CFBundlePrimaryIcon struct {
			CFBundleIconFiles []string `json:"CFBundleIconFiles"`
			CFBundleIconName  string   `json:"CFBundleIconName"`
		} `json:"CFBundlePrimaryIcon"`
	} `json:"CFBundleIcons"`
	CFBundleIdentifier            string   `json:"CFBundleIdentifier"`
	CFBundleInfoDictionaryVersion string   `json:"CFBundleInfoDictionaryVersion"`
	CFBundleName                  string   `json:"CFBundleName"`
	CFBundlePackageType           string   `json:"CFBundlePackageType"`
	CFBundleShortVersionString    string   `json:"CFBundleShortVersionString"`
	CFBundleSignature             string   `json:"CFBundleSignature"`
	CFBundleSupportedPlatforms    []string `json:"CFBundleSupportedPlatforms"`
	CFBundleURLTypes              []struct {
		CFBundleTypeRole   string   `json:"CFBundleTypeRole"`
		CFBundleURLName    string   `json:"CFBundleURLName"`
		CFBundleURLSchemes []string `json:"CFBundleURLSchemes"`
	} `json:"CFBundleURLTypes"`
	CFBundleVersion        string `json:"CFBundleVersion"`
	DTCompiler             string `json:"DTCompiler"`
	DTPlatformBuild        string `json:"DTPlatformBuild"`
	DTPlatformName         string `json:"DTPlatformName"`
	DTPlatformVersion      string `json:"DTPlatformVersion"`
	DTSDKBuild             string `json:"DTSDKBuild"`
	DTSDKName              string `json:"DTSDKName"`
	DTXcode                string `json:"DTXcode"`
	DTXcodeBuild           string `json:"DTXcodeBuild"`
	LSRequiresIPhoneOS     bool   `json:"LSRequiresIPhoneOS"`
	MinimumOSVersion       string `json:"MinimumOSVersion"`
	NSAppTransportSecurity struct {
		NSAllowsArbitraryLoads bool `json:"NSAllowsArbitraryLoads"`
		NSExceptionDomains     struct {
			Localhost struct {
				NSExceptionAllowsInsecureHTTPLoads bool `json:"NSExceptionAllowsInsecureHTTPLoads"`
			} `json:"localhost"`
		} `json:"NSExceptionDomains"`
	} `json:"NSAppTransportSecurity"`
	NSCameraUsageDescription                 string   `json:"NSCameraUsageDescription"`
	NSLocationWhenInUseUsageDescription      string   `json:"NSLocationWhenInUseUsageDescription"`
	UIAppFonts                               []string `json:"UIAppFonts"`
	UIDeviceFamily                           []int    `json:"UIDeviceFamily"`
	UILaunchStoryboardName                   string   `json:"UILaunchStoryboardName"`
	UIRequiredDeviceCapabilities             []string `json:"UIRequiredDeviceCapabilities"`
	UISupportedInterfaceOrientations         []string `json:"UISupportedInterfaceOrientations"`
	UIViewControllerBasedStatusBarAppearance bool     `json:"UIViewControllerBasedStatusBarAppearance"`
	UIMainStoryboardFile                     string   `json:"UIMainStoryboardFile"`
}

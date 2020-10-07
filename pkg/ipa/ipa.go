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

func (p Package) FileBytes(name string) ([]byte, error) {
	f := p.File(name)
	if f == nil {
		return nil, fmt.Errorf("no such file: %s", name)
	}

	r, err := f.Open()
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(r)
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

func (a App) FileBytes(name string) ([]byte, error) {
	target := path.Join(a.dir(), name)
	return a.p.FileBytes(target)
}

func (a App) Plist() (*Plist, error) {
	b, err := a.p.FileBytes(path.Join(a.dir(), "Info.plist"))
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(b)
	var ret Plist
	decoder := plist.NewDecoder(buf)
	err = decoder.Decode(&ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

type Plist struct {
	BuildMachineOSBuild                 string   `json:"BuildMachineOSBuild,omitempty"`
	BGTaskSchedulerPermittedIdentifiers []string `json:"BGTaskSchedulerPermittedIdentifiers,omitempty"`
	CFBundleDevelopmentRegion           string   `json:"CFBundleDevelopmentRegion,omitempty"`
	CFBundleDisplayName                 string   `json:"CFBundleDisplayName,omitempty"`
	CFBundleExecutable                  string   `json:"CFBundleExecutable,omitempty"`
	CFBundleIconName                    string   `json:"CFBundleIconName,omitempty"`
	CFBundleIcons                       struct {
		CFBundlePrimaryIcon struct {
			CFBundleIconFiles []string `json:"CFBundleIconFiles,omitempty"`
			CFBundleIconName  string   `json:"CFBundleIconName,omitempty"`
		} `json:"CFBundlePrimaryIcon,omitempty"`
	} `json:"CFBundleIcons,omitempty"`
	CFBundleIdentifier            string   `json:"CFBundleIdentifier,omitempty"`
	CFBundleInfoDictionaryVersion string   `json:"CFBundleInfoDictionaryVersion,omitempty"`
	CFBundleName                  string   `json:"CFBundleName,omitempty"`
	CFBundlePackageType           string   `json:"CFBundlePackageType,omitempty"`
	CFBundleShortVersionString    string   `json:"CFBundleShortVersionString,omitempty"`
	CFBundleSignature             string   `json:"CFBundleSignature,omitempty"`
	CFBundleSupportedPlatforms    []string `json:"CFBundleSupportedPlatforms,omitempty"`
	CFBundleURLTypes              []struct {
		CFBundleTypeRole   string   `json:"CFBundleTypeRole,omitempty"`
		CFBundleURLName    string   `json:"CFBundleURLName,omitempty"`
		CFBundleURLSchemes []string `json:"CFBundleURLSchemes,omitempty"`
	} `json:"CFBundleURLTypes,omitempty"`
	CFBundleVersion         string `json:"CFBundleVersion,omitempty"`
	DTCompiler              string `json:"DTCompiler,omitempty"`
	DTPlatformBuild         string `json:"DTPlatformBuild,omitempty"`
	DTPlatformName          string `json:"DTPlatformName,omitempty"`
	DTPlatformVersion       string `json:"DTPlatformVersion,omitempty"`
	DTSDKBuild              string `json:"DTSDKBuild,omitempty"`
	DTSDKName               string `json:"DTSDKName,omitempty"`
	DTXcode                 string `json:"DTXcode,omitempty"`
	DTXcodeBuild            string `json:"DTXcodeBuild,omitempty"`
	LSRequiresIPhoneOS      bool   `json:"LSRequiresIPhoneOS"`
	MinimumOSVersion        string `json:"MinimumOSVersion,omitempty"`
	NSHighResolutionCapable bool   `json:"NSHighResolutionCapable"`
	NSAppTransportSecurity  struct {
		NSAllowsArbitraryLoads bool `json:"NSAllowsArbitraryLoads,omitempty"`
		NSExceptionDomains     map[string]struct {
			NSExceptionAllowsInsecureHTTPLoads bool   `json:"NSExceptionAllowsInsecureHTTPLoads,omitempty"`
			NSIncludesSubdomains               bool   `json:"NSIncludesSubdomains,omitempty"`
			NSExceptionMinimumTLSVersion       string `json:"NSExceptionMinimumTLSVersion,omitempty"`
			NSExceptionRequiresForwardSecrecy  bool   `json:"NSExceptionRequiresForwardSecrecy,omitempty"`
			NSRequiresCertificateTransparency  bool   `json:"NSRequiresCertificateTransparency,omitempty"`
		} `json:"NSExceptionDomains,omitempty"`
	} `json:"NSAppTransportSecurity,omitempty"`
	NSCameraUsageDescription                 string   `json:"NSCameraUsageDescription,omitempty"`
	NSLocationWhenInUseUsageDescription      string   `json:"NSLocationWhenInUseUsageDescription,omitempty"`
	UIAppFonts                               []string `json:"UIAppFonts,omitempty"`
	UIBackgroundModes                        []string `json:"UIBackgroundModes,omitempty"`
	UIDeviceFamily                           []int    `json:"UIDeviceFamily,omitempty"`
	UILaunchStoryboardName                   string   `json:"UILaunchStoryboardName,omitempty"`
	UIRequiredDeviceCapabilities             []string `json:"UIRequiredDeviceCapabilities,omitempty"`
	UISupportedInterfaceOrientations         []string `json:"UISupportedInterfaceOrientations,omitempty"`
	UIViewControllerBasedStatusBarAppearance bool     `json:"UIViewControllerBasedStatusBarAppearance"`
	UIMainStoryboardFile                     string   `json:"UIMainStoryboardFile,omitempty"`

	// not standard
	Shake map[string]interface{}
}

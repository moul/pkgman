package apk

import (
	"encoding/xml"
	"strings"
)

// Manifest implements https://developer.android.com/guide/topics/manifest/manifest-element
type Manifest struct {
	XMLName                   xml.Name `xml:"manifest"`
	Android                   string   `xml:"android,attr,omitempty"`
	VersionCode               string   `xml:"versionCode,attr,omitempty"`
	VersionName               string   `xml:"versionName,attr,omitempty"`
	CompileSdkVersion         string   `xml:"compileSdkVersion,attr,omitempty"`
	CompileSdkVersionCodename string   `xml:"compileSdkVersionCodename,attr,omitempty"`
	Package                   string   `xml:"package,attr,omitempty"`
	PlatformBuildVersionCode  string   `xml:"platformBuildVersionCode,attr,omitempty"`
	PlatformBuildVersionName  string   `xml:"platformBuildVersionName,attr,omitempty"`
	UsesSdk                   struct {
		MinSdkVersion    string `xml:"minSdkVersion,attr,omitempty"`
		TargetSdkVersion string `xml:"targetSdkVersion,attr,omitempty"`
	} `xml:"uses-sdk,omitempty"`
	UsesPermission []struct {
		Name string `xml:"name,attr,omitempty"`
	} `xml:"uses-permission,omitempty"`
	UsesFeature []struct {
		Name     string `xml:"name,attr,omitempty"`
		Required bool   `xml:"required,attr,omitempty"`
	} `xml:"uses-feature,omitempty"`
	Application ManifestApplication `xml:"application"`
}

// ManifestApplication implements https://developer.android.com/guide/topics/manifest/application-element
type ManifestApplication struct {
	Theme                 string             `xml:"theme,attr,omitempty"`
	Label                 string             `xml:"label,attr,omitempty"`
	Icon                  string             `xml:"icon,attr,omitempty"`
	Name                  string             `xml:"name,attr,omitempty"`
	AllowBackup           bool               `xml:"allowBackup,attr,omitempty"`
	SupportsRtl           bool               `xml:"supportsRtl,attr,omitempty"`
	NetworkSecurityConfig string             `xml:"networkSecurityConfig,attr,omitempty"`
	RoundIcon             string             `xml:"roundIcon,attr,omitempty"`
	AppComponentFactory   string             `xml:"appComponentFactory,attr,omitempty"`
	Activity              []ManifestActivity `xml:"activity,omitempty"`
	MetaData              []ManifestMetadata `xml:"meta-data,omitempty"`
	Service               []ManifestService  `xml:"service,omitempty"`
	Provider              []ManifestProvider `xml:"provider,omitempty"`
	Receiver              struct {
		Name     string `xml:"name,attr,omitempty"`
		Exported bool   `xml:"exported,attr,omitempty"`
	} `xml:"receiver,omitempty"`
}

// ManifestProvider implements https://developer.android.com/guide/topics/manifest/provider-element
type ManifestProvider struct {
	Name                string               `xml:"name,attr,omitempty"`
	Exported            bool                 `xml:"exported,attr,omitempty"`
	Authorities         string               `xml:"authorities,attr,omitempty"`
	GrantURIPermissions bool                 `xml:"grantUriPermissions,attr,omitempty"`
	InitOrder           string               `xml:"initOrder,attr,omitempty"`
	DirectBootAware     bool                 `xml:"directBootAware,attr,omitempty"`
	MetaData            ManifestMetadata     `xml:"meta-data,omitempty"`
	IntentFilter        ManifestIntentFilter `xml:"intent-filter"`
}

// ManifestService implements https://developer.android.com/guide/topics/manifest/service-element
type ManifestService struct {
	Name                  string               `xml:"name,attr,omitempty"`
	StopWithTask          bool                 `xml:"stopWithTask,attr,omitempty"`
	Process               string               `xml:"process,attr,omitempty"`
	IsolatedProcess       bool                 `xml:"isolatedProcess,attr,omitempty"`
	Exported              bool                 `xml:"exported,attr,omitempty"`
	Description           string               `xml:"description,attr,omitempty"`
	DirectBootAware       bool                 `xml:"directBootAware,attr,omitempty"`
	ForegroundServiceType string               `xml:"foregroundServiceType,attr,omitempty"`
	Enabled               bool                 `xml:"enabled,attr,omitempty"`
	Permission            string               `xml:"permission,attr,omitempty"`
	MetaData              []ManifestMetadata   `xml:"meta-data,omitempty"`
	IntentFilter          ManifestIntentFilter `xml:"intent-filter"`
}

// ManifestMetadata implements https://developer.android.com/guide/topics/manifest/meta-data-element0
type ManifestMetadata struct {
	Name     string `xml:"name,attr,omitempty"`
	Value    string `xml:"value,attr,omitempty"`
	Resource string `xml:"resource,attr,omitempty"`
}

// ManifestActivity implements https://developer.android.com/guide/topics/manifest/activity-element
type ManifestActivity struct {
	Label               string               `xml:"label,attr,omitempty"`
	Name                string               `xml:"name,attr,omitempty"`
	LaunchMode          string               `xml:"launchMode,attr,omitempty"`
	ScreenOrientation   string               `xml:"screenOrientation,attr,omitempty"`
	ConfigChanges       string               `xml:"configChanges,attr,omitempty"`
	WindowSoftInputMode string               `xml:"windowSoftInputMode,attr,omitempty"`
	Theme               string               `xml:"theme,attr,omitempty"`
	Exported            bool                 `xml:"exported,attr,omitempty"`
	IntentFilter        ManifestIntentFilter `xml:"intent-filter"`
}

// ManifestIntentFilter implements https://developer.android.com/guide/topics/manifest/intent-filter-element
type ManifestIntentFilter struct {
	Action struct {
		Name string `xml:"name,attr,omitempty"`
	} `xml:"action"`
	Category struct {
		Name string `xml:"name,attr,omitempty"`
	} `xml:"category"`
}

func (m Manifest) MainActivity() *ManifestActivity {
	for _, activity := range m.Application.Activity {
		if strings.HasSuffix(activity.Name, ".MainActivity") {
			return &activity
		}
	}
	return nil
}

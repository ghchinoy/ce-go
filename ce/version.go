package ce

const (
	version     = "0.1.0"
	versionName = "ce-go"
)

// Version returns the version string
func Version() string {
	return version
}

// VersionName returns the name of the version
func VersionName() string {
	return versionName
}

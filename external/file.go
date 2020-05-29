package external

// MboxFile is a filesystem file that contains mail folder.
type MboxFile struct {
	// File is absolute (or relative) file location.
	File string
	// Name is folder name.
	Name string
}

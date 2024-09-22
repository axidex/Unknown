package archive

type Info struct {
	Archive *archive
}

type archive struct {
	Size int
	Hash string
}

package app

import "zelin/internal/files"

func (s *FileService) Image(path string) (string, error) {
	image, err := files.ReadImage(path)
	return image.DataURL, err
}

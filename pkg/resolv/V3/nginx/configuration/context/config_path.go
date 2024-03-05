package context

import (
	"github.com/marmotedu/errors"
	"path/filepath"
)

type ConfigPath interface {
	FullPath() string
	BaseDir() string
	RelativePath() string
}

type RelConfigPath struct {
	configHomeDir      string
	configRelativePath string
}

func (c RelConfigPath) FullPath() string {
	return filepath.Clean(filepath.Join(c.configHomeDir, c.configRelativePath))
}

func (c RelConfigPath) BaseDir() string {
	return c.configHomeDir
}

func (c RelConfigPath) RelativePath() string {
	return c.configRelativePath
}

func NewRelConfigPath(dir, file string) (ConfigPath, error) {
	cleanDir := filepath.Clean(dir)

	if filepath.IsAbs(file) {
		rel, err := filepath.Rel(cleanDir, file)
		if err != nil {
			return nil, err
		}
		return &RelConfigPath{
			configHomeDir:      cleanDir,
			configRelativePath: rel,
		}, nil
	}

	return &RelConfigPath{
		configHomeDir:      cleanDir,
		configRelativePath: file,
	}, nil
}

type AbsConfigPath struct {
	fullPath string
}

func (a AbsConfigPath) FullPath() string {
	return a.fullPath
}

func (a AbsConfigPath) BaseDir() string {
	return filepath.Dir(a.fullPath)
}

func (a AbsConfigPath) RelativePath() string {
	_, filename := filepath.Split(a.fullPath)
	return filename
}

func NewAbsConfigPath(absPath string) (ConfigPath, error) {
	if !filepath.IsAbs(absPath) {
		return nil, errors.Errorf("%s is not a absolute path", absPath)
	}
	return &AbsConfigPath{filepath.Clean(absPath)}, nil
}

type ConfigGraph interface {
	AddEdge(src, dst ConfigPath) error
	CropOffEde(src, dst ConfigPath)
	MainConfigPath() ConfigPath
}

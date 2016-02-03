/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package cursor

import (
	"os"
	"path"

	"pkg.deepin.io/dde/api/thumbnails/loader"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	presentCursorLeftPtr   = "left_ptr"
	presentCursorLeftWatch = "left_ptr_watch"
	presentCursorQuestion  = "question_arrow"
)

const (
	defaultWidth    = 128
	defaultHeight   = 72
	defaultIconSize = 24
)

func doGenThumbnail(src, bg, dest string, width, height int, force, theme bool) (string, error) {
	if !force && dutils.IsFileExist(dest) {
		return dest, nil
	}

	src = dutils.DecodeURI(src)
	bg = dutils.DecodeURI(bg)
	dir := path.Dir(src)
	tmp := loader.GetTmpImage()
	err := loader.CompositeIcons(getCursorIcons(dir), bg, tmp,
		defaultIconSize, defaultWidth, defaultHeight)
	os.RemoveAll(xcur2pngCache)
	if err != nil {
		return "", err
	}

	defer os.Remove(tmp)
	if !theme {
		err = loader.ThumbnailImage(tmp, dest, width, height)
	} else {
		err = loader.ScaleImage(tmp, dest, width, height)
	}
	if err != nil {
		return "", err
	}

	return dest, nil
}

func getCursorIcons(dir string) []string {
	presents := []string{
		presentCursorLeftPtr,
		presentCursorLeftWatch,
		presentCursorQuestion,
	}

	var files []string
	for _, name := range presents {
		tmp, err := XCursorToPng(path.Join(dir, "cursors", name))
		if err != nil {
			return nil
		}
		files = append(files, tmp)
	}
	return files
}

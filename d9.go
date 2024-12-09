package main

import (
	"bytes"
	"fmt"
	"io"
)

type File struct {
	Position int
	Size     int
	Id       int
}

func checksum(drive []int) int {
	checksum := 0
	for i, c := range drive {
		if c == -1 {
			continue
		}
		val := c
		checksum += val * i
	}
	return checksum
}

func (self *Day) D9() error {

	files := []File{}

	for {
		line, err := self.Input.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			break
		}

		position := 0
		for i := 0; i < len(line); i += 2 {
			id := i / 2
			size := int(line[i] - '0')
			padding := 0
			if i+1 < len(line) {
				padding = int(line[i+1] - '0')
			}
			files = append(files, File{
				Position: position,
				Size:     size,
				Id:       id,
			})
			position += size + padding
		}
	}

	last := files[len(files)-1]

	drive := make([]int, last.Position+last.Size)
	for i := range drive {
		drive[i] = -1
	}

	for _, file := range files {
		for i := 0; i < file.Size; i++ {
			drive[file.Position+i] = file.Id
		}
	}

	Debug(drive)
	lowest := 0
	found := 0
	for i := len(drive) - 1; i >= 0; i-- {
		if drive[i] == -1 {
			continue
		}
		val := drive[i]
		found += 1
		for d := lowest; d < len(drive); d++ {
			if d >= i {
				break
			}
			if drive[d] == -1 {
				drive[i] = -1
				drive[d] = val
				lowest = d
				// Debug(string(drive))
				break
			}
		}
	}

	Debug(drive)
	fmt.Println(checksum(drive))

	drive2 := make([]int, last.Position+last.Size)
	for i := range drive {
		drive2[i] = -1
	}
	for _, file := range files {
		for i := 0; i < file.Size; i++ {
			drive2[file.Position+i] = file.Id
		}
	}

	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		count := 0
		position := 0
		for i, id := range drive2 {
			if i > file.Position {
				break
			}
			if id != -1 {
				if count >= file.Size {
					Debug("Moving", file, "to", position)
					for i := position; i < position+file.Size; i++ {
						drive2[i] = file.Id
					}
					for i := file.Position; i < file.Position+file.Size; i++ {
						drive2[i] = -1
					}
					break
				}
				count = 0
			} else {
				if count == 0 {
					position = i
				}
				count += 1
			}
		}

	}
	Debug(drive2)
	fmt.Println(checksum(drive2))

	return nil
}

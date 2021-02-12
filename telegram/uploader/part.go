package uploader

import "golang.org/x/xerrors"

// https://core.telegram.org/api/files#uploading-files
const (
	// Use upload.saveBigFilePart in case the full size of the file is more than 10 MB
	// and upload.saveFilePart for smaller files
	bigFileLimit = 10 * 1024 * 1024 // 10 MB

	// Each part should have a sequence number, file_part, with a value ranging from 0 to 2,999.
	partsLimit = 2999

	defaultPartSize = 1024 // 1 KB
	// The file’s binary content is then split into parts. All parts must have the same size (part_size)
	// and the following conditions must be met:

	// `part_size % 1024 = 0` (divisible by 1KB)
	paddingPartSize = 1024
	// `524288 % part_size = 0` (512KB must be evenly divisible by part_size)
	MaximumPartSize = 524288
)

func (u *Uploader) checkPartSize() error {
	switch {
	case u.partSize == 0:
		return xerrors.New("is equal to zero")
	case u.partSize%paddingPartSize != 0:
		return xerrors.Errorf("%d is not divisible by 1024", u.partSize)
	case MaximumPartSize%u.partSize != 0:
		return xerrors.Errorf("524288 is not divisible by %d", u.partSize)
	}

	return nil
}

func (u *Uploader) computeParts(upload *Upload) int {
	if upload.totalBytes <= 0 {
		return 0
	}

	totalBytes := int(upload.totalBytes)
	parts := totalBytes / u.partSize
	if totalBytes%u.partSize != 0 {
		parts++
	}
	return parts
}

func (u *Uploader) initUpload(upload *Upload) error {
	big := upload.totalBytes > bigFileLimit
	totalParts := u.computeParts(upload)
	if !big && totalParts > partsLimit {
		return xerrors.Errorf(
			"part size is too small: total size = %d, part size = %d, %d / %d > %d",
			upload.totalBytes, u.partSize, upload.totalBytes, u.partSize, partsLimit,
		)
	}

	if upload.id == 0 {
		id, err := u.id()
		if err != nil {
			return xerrors.Errorf("id generation: %w", err)
		}

		upload.id = id
		upload.partSize = u.partSize
	} else if upload.partSize != u.partSize {
		return xerrors.Errorf(
			"previous upload has part size %d, but uploader size is %d",
			upload.partSize, u.partSize,
		)
	}

	upload.big = big
	upload.totalParts = totalParts
	return nil
}

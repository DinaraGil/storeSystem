package helpers

type FileDataType struct {
	FileName string
	Data     []byte
}

type OperationError struct {
	ObjectID string
	Error    error
}

type UploadedFile struct {
	ObjectID string
	Link     string
}

package err

import "errors"

var ErrDocumentNotFound = errors.New("document not found")
var ErrCollectionAlreadyExists = errors.New("collection already exists")
var ErrCollectionNotFound = errors.New("collection not found")
var ErrUnsupportedDocumentField = errors.New("unsupported document field")
var ErrCreatedUser = errors.New("error creating user")
var ErrListEmpty = errors.New("the list is empty")
var ErrNotFound = errors.New("not found")
var ErrAddUser = errors.New("error adding user")

package fileupload

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/alejogs4/blog/src/post/domain/post"
	"github.com/alejogs4/blog/src/shared/infraestructure/httputils"
	"github.com/alejogs4/blog/src/shared/infraestructure/middleware"
)

var errCopyingFile = errors.New("File: file must be present")

func uploadFile(request *http.Request, formField, folder string) (string, error) {
	file, _, err := request.FormFile(formField)
	if err != nil {
		return "", post.ErrMissingPostPicture
	}
	defer file.Close()

	newFile, err := ioutil.TempFile(folder, "upload-*.jpeg")
	if err != nil {
		return "", errCopyingFile
	}
	picturePath := "/" + newFile.Name()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", errCopyingFile
	}

	_, err = newFile.Write(fileBytes)
	if err != nil {
		return "", errCopyingFile
	}
	return picturePath, nil
}

// UploadFile middleware that upload an incoming file in a given existing folder
func UploadFile(formField, folder string) middleware.Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(response http.ResponseWriter, request *http.Request) {
			response.Header().Set("Content-Type", "application/json")

			filePath, err := uploadFile(request, formField, folder)
			if err != nil {
				httpError := mapFileErrorToHttpError(err)
				httputils.DispatchNewHttpError(response, httpError.Message, httpError.Status)
				return
			}

			// Look for what else I can use here
			newContext := context.WithValue(request.Context(), "file", filePath)
			f(response, request.WithContext(newContext))
		}
	}
}

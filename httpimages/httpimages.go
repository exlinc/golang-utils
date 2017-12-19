package httpimages

import (
	"errors"
	"fmt"
	"git.exlhub.io/exlinc/golang-utils/uuid"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"image"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
)

type ImgType int

const (
	PNG ImgType = iota
	JPG
)

type UploaderConfig struct {
	UploadsDir   string
	S3Session    *session.Session
	S3BucketName string
	S3ACL        string
	ImageType    ImgType
}

func FormatAndContentTypeForImgType(t ImgType) (format imgio.Format, contentType, extension string) {
	switch t {
	case JPG:
		return imgio.JPEG, "image/jpeg", ".jpg"
	default:
		return imgio.PNG, "image/png", ".png"
	}
}

func (cfg *UploaderConfig) HandleAvatarUploadToS3(r *http.Request, fileKey string, exportWidth, exportHeight int, exportType ImgType) (url string, status int, err error) {
	resourceName, resourcePath, _, status, err := cfg.ImageUploadToFile(r, fileKey)
	if err != nil {
		return "", status, err
	}
	defer os.Remove(resourcePath)
	resourcePathAfter := fmt.Sprintf("%s.after", resourcePath)
	defer os.Remove(resourcePathAfter)

	img, err := ResizeImage(resourcePath, exportHeight, exportWidth)
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	format, contentType, extension := FormatAndContentTypeForImgType(exportType)
	err = imgio.Save(resourcePathAfter, img, format)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	var fileName string
	fileName = fmt.Sprintf("%s_%s_%s%s",
		resourceName,
		fmt.Sprint(exportWidth),
		fmt.Sprint(exportHeight),
		extension)

	avatarUrl, err := cfg.SendImageToS3(resourcePath, fileName, contentType)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	return avatarUrl, http.StatusInternalServerError, err
}

func (cfg *UploaderConfig) SendImageToS3(localFilePath, remoteFilePath, contentType string) (string, error) {
	if cfg.S3BucketName == "" {
		return "", errors.New("missing s3 bucket name")
	}
	if cfg.S3Session == nil {
		return "", errors.New("missing s3 session")
	}
	file, err := os.Open(localFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	acl := cfg.S3ACL
	if acl == "" {
		acl = "private"
	}

	res, err := s3manager.NewUploader(cfg.S3Session).Upload(&s3manager.UploadInput{
		Bucket:      aws.String(cfg.S3BucketName),
		Key:         aws.String(remoteFilePath),
		ACL:         aws.String(acl),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	return res.Location, nil
}

func (cfg *UploaderConfig) ImageUploadToFile(r *http.Request, fileKey string) (resourceName string, resourcePath string, contentType string, status int, err error) {
	if cfg.UploadsDir == "" {
		return "", "", "", http.StatusInternalServerError, errors.New("missing upload directory")
	}
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile(fileKey)
	if err != nil {
		return "", "", "", http.StatusBadRequest, err
	}
	defer file.Close()

	contentType = handler.Header.Get("Content-Type")
	if contentType != "image/gif" && contentType != "image/jpeg" && contentType != "image/png" {
		return "", "", "", http.StatusUnsupportedMediaType, errors.New("invalid content-type. Requires gif, jpeg, or png format")
	}

	resourceExtension, err := mime.ExtensionsByType(contentType)
	if err != nil || resourceExtension == nil || len(resourceExtension) < 1 {
		return "", "", "", http.StatusBadRequest, errors.New("error determining file extension")
	}

	resourceName = uuid.NewV4().String()
	resourcePath = path.Join(cfg.UploadsDir, resourceName, resourceExtension[0])
	localImage, err := os.OpenFile(resourcePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", "", "", http.StatusInternalServerError, err
	}
	defer localImage.Close()
	_, err = io.Copy(localImage, file)
	if err != nil {
		return "", "", "", http.StatusInternalServerError, err
	}
	return resourceName, resourcePath, contentType, http.StatusOK, nil
}

func ResizeImage(path string, width, height int) (*image.RGBA, error) {
	imgBefore, err := imgio.Open(path)
	if err != nil {
		return nil, err
	}
	imgAfter := transform.Resize(imgBefore, width, height, transform.Lanczos)
	return imgAfter, nil
}

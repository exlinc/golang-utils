package httpimages

import (
	"errors"
	"fmt"
	"exlgit.com/golang/utils/uuid"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"image"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"time"
	"path/filepath"
)

// ImgType is a type used to track image formats
type ImgType int

const (
	// PNG is a PNG image format
	PNG ImgType = iota
	// JPG is a JPG/JPEG image format
	JPG
)

// NewS3UploaderConfig creates an S3Session to bypass vendored lib conflicts across project -- it also enforces using the AWS standard config
func NewS3UploaderConfig() (*UploaderConfig, error) {
	s, err := session.NewSession(&aws.Config{})
	if err != nil {
		return nil, err
	}
	uc := UploaderConfig{
		S3Session: s,
	}
	return &uc, nil
}

// UploaderConfig outlines various configuration options for upload-related requests, they are required/optional on a per-method basis
type UploaderConfig struct {
	// UploadsDir is the local file path where images will be stored by certain methods
	UploadsDir string
	// S3Session is the S3 session used for S3-related methods
	S3Session *session.Session
	// S3BucketName is the S3 bucket used for S3-related methods
	S3BucketName string
	// S3ACL is the ACL used by S3-related methods -- it defaults to 'private' in most methods
	S3ACL string
}

// FormatAndContentTypeForImgType is a utility to get the imgio.Format, content type, and file extension from the ImgType const
func FormatAndContentTypeForImgType(t ImgType) (format imgio.Format, contentType, extension string) {
	switch t {
	case JPG:
		return imgio.JPEG, "image/jpeg", ".jpg"
	default:
		return imgio.PNG, "image/png", ".png"
	}
}

// HandleAvatarUploadToFile takes an HTTP request, handles the file upload, resizes the image, and then puts the result onto the FS in the upload directory
func (cfg *UploaderConfig) HandleAvatarUploadToFile(r *http.Request, fileKey string, exportWidth, exportHeight int, exportType ImgType) (destName string, status int, err error) {
	resourceName, resourcePath, _, status, err := cfg.HandleImageUploadToFile(r, fileKey)
	if err != nil {
		return "", status, err
	}
	defer os.Remove(resourcePath)
	resourcePathAfter := fmt.Sprintf("%s.after", resourcePath)

	img, err := ResizeImage(resourcePath, exportHeight, exportWidth)
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	format, _, extension := FormatAndContentTypeForImgType(exportType)
	f, err := os.Create(resourcePathAfter)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	defer f.Close()
	err = imgio.Encode(f, img, format)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	destName = fmt.Sprintf("%s_%s_%s%s",
		resourceName,
		fmt.Sprint(exportWidth),
		fmt.Sprint(exportHeight),
		extension)

	err = os.Rename(resourcePathAfter, filepath.Join(cfg.UploadsDir, destName))
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	return destName, http.StatusOK, err
}

// HandleAvatarUploadToS3 takes an HTTP request, handles the file upload, resizes the image, and then puts the result onto S3 -- returning the URL and HTTP status code to the caller
func (cfg *UploaderConfig) HandleAvatarUploadToS3(r *http.Request, fileKey string, exportWidth, exportHeight int, exportType ImgType) (url string, destKey string, status int, err error) {
	resourceName, resourcePath, _, status, err := cfg.HandleImageUploadToFile(r, fileKey)
	if err != nil {
		return "", "", status, err
	}
	defer os.Remove(resourcePath)
	resourcePathAfter := fmt.Sprintf("%s.after", resourcePath)
	defer os.Remove(resourcePathAfter)

	img, err := ResizeImage(resourcePath, exportHeight, exportWidth)
	if err != nil {
		return "", "", http.StatusBadRequest, err
	}
	format, contentType, extension := FormatAndContentTypeForImgType(exportType)
	err = imgio.Save(resourcePathAfter, img, format)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	var fileName string
	fileName = fmt.Sprintf("%s_%s_%s%s",
		resourceName,
		fmt.Sprint(exportWidth),
		fmt.Sprint(exportHeight),
		extension)

	avatarUrl, err := cfg.SendImageToS3(resourcePath, fileName, contentType)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	return avatarUrl, fileName, http.StatusOK, err
}

// SendImageToS3 takes an image from the filesystem and puts it onto S3
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

// HandleImageUploadToFile takes an HTTP request, verifies the upload is a supported image, and manages putting the file onto the filesystem -- returning data about the file and the HTTP status code to the caller
func (cfg *UploaderConfig) HandleImageUploadToFile(r *http.Request, fileKey string) (resourceName string, resourcePath string, contentType string, status int, err error) {
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
	resourcePath = path.Join(cfg.UploadsDir, fmt.Sprintf("%s%s", resourceName, resourceExtension[0]))
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

// ResizeImage is a utility to get the resized version of an image using the Lanczos transform
func ResizeImage(path string, width, height int) (*image.RGBA, error) {
	imgBefore, err := imgio.Open(path)
	if err != nil {
		return nil, err
	}
	imgAfter := transform.Resize(imgBefore, width, height, transform.Lanczos)
	return imgAfter, nil
}

// PresignS3GetObjectRequest is a utility to presign an S3 get object request for a period of time -- returning the presigned URL to the caller. It uses the AWS standard environment/CLI configuration
func PresignS3GetObjectRequest(bucketName, itemName string, duration time.Duration) (url string, err error) {
	s, err := session.NewSession(&aws.Config{})
	if err != nil {
		return "", err
	}
	svc := s3.New(s)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(itemName),
	})
	return req.Presign(duration)
}

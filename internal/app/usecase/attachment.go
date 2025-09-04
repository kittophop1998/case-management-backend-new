package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"case-management/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttachmentUseCase struct {
	repo repository.AttachmentsRepository
}

func NewAttachmentUseCase(repo repository.AttachmentsRepository) *AttachmentUseCase {
	return &AttachmentUseCase{repo: repo}
}

func (u *AttachmentUseCase) UploadAttachment(c *gin.Context, files []*multipart.FileHeader, caseID, userID uuid.UUID, isilonBaseURL, isilonUser, isilonPass string) error {
	for _, file := range files {
		extension := filepath.Ext(file.Filename)

		randomName, err := utils.RandStringRunes(10)
		if err != nil {
			return err
		}

		newFileName := fmt.Sprintf("/ifs/case/%s/%s%s", caseID.String(), randomName, extension)

		fileContent, err := file.Open()
		if err != nil {
			return err
		}
		defer fileContent.Close()

		// upload ไปยัง Isilon
		url := fmt.Sprintf("%s%s", isilonBaseURL, newFileName)
		req, err := http.NewRequest("PUT", url, fileContent)
		if err != nil {
			return err
		}
		req.SetBasicAuth(isilonUser, isilonPass)
		req.Header.Set("Content-Type", file.Header.Get("Content-Type"))

		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			return fmt.Errorf("isilon upload failed: %s", resp.Status)
		}

		attachment := model.Attachment{
			CaseId:           caseID,
			FileName:         file.Filename,
			FilePath:         newFileName,
			FileType:         file.Header.Get("Content-Type"),
			FileSizeBytes:    uint64(file.Size),
			UploadedByUserId: userID,
			UploadedAt:       time.Now(),
		}

		err = u.SaveAttachmentRecord(c, caseID, attachment)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *AttachmentUseCase) SaveAttachmentRecord(ctx *gin.Context, caseID uuid.UUID, attachment model.Attachment) error {
	// ใช้ repo ของคุณ
	_, err := u.repo.UploadAttachment(ctx, caseID, attachment)
	return err
}

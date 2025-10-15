package upload

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"auto-forge/internal/dto/response"
	"auto-forge/internal/models"
	uploadRepo "auto-forge/internal/repositories/upload"
	"auto-forge/pkg/common"
	"auto-forge/pkg/constants"
	"auto-forge/pkg/database"
	"auto-forge/pkg/logger"
	"auto-forge/pkg/upload"
)

var uploadService *UploadService
var config *upload.Config


type UploadService struct {
	uploadRepo *uploadRepo.UploadRepository
	storage    upload.Storage
	config     *upload.Config
}


func InitUploadService() {
	config = upload.NewDefaultConfig()
	db := database.GetDB()
	uploadService = &UploadService{
		uploadRepo: uploadRepo.NewUploadRepository(db),
		storage:    upload.NewLocalStorage(config.UploadDir),
		config:     config,
	}
}


func GetUploadService() *UploadService {
	return uploadService
}


func SimpleUpload(file *multipart.FileHeader, userID string) (*response.SimpleUploadResponse, error) {
	if uploadService == nil {
		InitUploadService()
	}
	return uploadService.SimpleUpload(file, userID)
}


func InitChunkUpload(filename string, fileSize int64, md5Hash string, chunkSize int64, userID string) (*response.ChunkUploadInitResponse, error) {
	if uploadService == nil {
		InitUploadService()
	}
	return uploadService.InitChunkUpload(filename, fileSize, md5Hash, chunkSize, userID)
}


func UploadChunk(fileID string, chunkIndex int, md5Hash string, chunk *multipart.FileHeader) (*response.ChunkUploadResponse, error) {
	if uploadService == nil {
		InitUploadService()
	}
	return uploadService.UploadChunk(fileID, chunkIndex, md5Hash, chunk)
}


func MergeChunks(fileID string) (*response.ChunkMergeResponse, error) {
	if uploadService == nil {
		InitUploadService()
	}
	return uploadService.MergeChunks(fileID)
}


func GetUploadProgress(fileID string) (*response.UploadProgressResponse, error) {
	if uploadService == nil {
		InitUploadService()
	}
	return uploadService.GetUploadProgress(fileID)
}


func (s *UploadService) SimpleUpload(file *multipart.FileHeader, userID string) (*response.SimpleUploadResponse, error) {

	if err := s.validateFile(file); err != nil {
		return nil, err
	}


	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %v", err)
	}
	defer src.Close()


	md5Hash, err := upload.CalculateFileMD5(src)
	if err != nil {
		return nil, fmt.Errorf("计算文件MD5失败: %v", err)
	}


	if existingFile, err := s.uploadRepo.GetUploadFileByMD5(md5Hash); err == nil {
		logger.Info("文件已存在，执行秒传", "md5", md5Hash, "fileID", existingFile.ID.String())
		return &response.SimpleUploadResponse{
			FileID:     existingFile.ID.String(),
			Filename:   existingFile.Filename,
			StoredName: existingFile.StoredName,
			FileSize:   existingFile.FileSize,
			MimeType:   existingFile.MimeType,
			Extension:  existingFile.Extension,
			MD5Hash:    existingFile.MD5Hash,
			FilePath:   fmt.Sprintf("/files/preview/%s/%s", existingFile.ID.String(), existingFile.Filename),
			UploadedAt: *existingFile.UploadedAt,
		}, nil
	}


	storedName := upload.GenerateStoredFilename(file.Filename)


	filePath, err := s.storage.SaveFile(storedName, src)
	if err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}


	now := time.Now()
	uploadFile := &models.UploadFile{
		BaseModel:     models.BaseModel{ID: common.NewUUID()},
		Filename:      file.Filename,
		StoredName:    storedName,
		FilePath:      filePath,
		FileSize:      file.Size,
		MimeType:      upload.GetMimeTypeFromExtension(file.Filename),
		Extension:     upload.GetFileExtension(file.Filename),
		MD5Hash:       md5Hash,
		UploadStatus:  constants.UploadStatusCompleted,
		ChunkTotal:    1,
		ChunkUploaded: 1,
		UserID:        userID,
		UploadedAt:    &now,
		IsPublic:      true,
	}

	if err := s.uploadRepo.CreateUploadFile(uploadFile); err != nil {

		s.storage.DeleteFile(filePath)
		return nil, fmt.Errorf("保存文件记录失败: %v", err)
	}

	return &response.SimpleUploadResponse{
		FileID:     uploadFile.ID.String(),
		Filename:   uploadFile.Filename,
		StoredName: uploadFile.StoredName,
		FileSize:   uploadFile.FileSize,
		MimeType:   uploadFile.MimeType,
		Extension:  uploadFile.Extension,
		MD5Hash:    uploadFile.MD5Hash,
		FilePath:   fmt.Sprintf("/files/preview/%s/%s", uploadFile.ID.String(), uploadFile.Filename),
		UploadedAt: *uploadFile.UploadedAt,
	}, nil
}


func (s *UploadService) InitChunkUpload(filename string, fileSize int64, md5Hash string, chunkSize int64, userID string) (*response.ChunkUploadInitResponse, error) {

	if err := upload.ValidateFilename(filename); err != nil {
		return nil, err
	}

	if !s.config.ValidateFileSize(fileSize) {
		return nil, fmt.Errorf("文件大小超出限制，最大允许 %d 字节", s.config.MaxFileSize)
	}

	mimeType := upload.GetMimeTypeFromExtension(filename)
	if !s.config.ValidateMimeType(mimeType) {
		return nil, fmt.Errorf("不支持的文件类型: %s", mimeType)
	}


	if existingFile, err := s.uploadRepo.GetUploadFileByMD5(md5Hash); err == nil {
		logger.Info("文件已存在，执行秒传", "md5", md5Hash, "fileID", existingFile.ID.String())
		return &response.ChunkUploadInitResponse{
			FileID:     existingFile.ID.String(),
			ChunkSize:  chunkSize,
			ChunkTotal: 1,
		}, nil
	}


	chunkTotal := upload.CalculateChunkTotal(fileSize, chunkSize)


	fileID := common.NewUUID()
	uploadFile := &models.UploadFile{
		BaseModel:     models.BaseModel{ID: fileID},
		Filename:      filename,
		StoredName:    upload.GenerateStoredFilename(filename),
		FileSize:      fileSize,
		MimeType:      mimeType,
		Extension:     upload.GetFileExtension(filename),
		MD5Hash:       md5Hash,
		UploadStatus:  constants.UploadStatusUploading,
		ChunkTotal:    chunkTotal,
		ChunkUploaded: 0,
		UserID:        userID,
		IsPublic:      true,
	}

	if err := s.uploadRepo.CreateUploadFile(uploadFile); err != nil {
		return nil, fmt.Errorf("创建文件记录失败: %v", err)
	}


	for i := 0; i < chunkTotal; i++ {
		chunkPath := upload.GenerateChunkPath(s.config.TempDir, fileID.String(), i)
		chunk := &models.ChunkInfo{
			BaseModel:  models.BaseModel{ID: common.NewUUID()},
			FileID:     fileID.String(),
			ChunkIndex: i,
			ChunkPath:  chunkPath,
			IsUploaded: false,
		}

		if err := s.uploadRepo.CreateChunkInfo(chunk); err != nil {
			return nil, fmt.Errorf("创建分片记录失败: %v", err)
		}
	}

	logger.Info("分片上传初始化成功", "fileID", fileID.String(), "chunkTotal", chunkTotal)

	return &response.ChunkUploadInitResponse{
		FileID:     fileID.String(),
		ChunkSize:  chunkSize,
		ChunkTotal: chunkTotal,
	}, nil
}


func (s *UploadService) UploadChunk(fileID string, chunkIndex int, md5Hash string, chunk *multipart.FileHeader) (*response.ChunkUploadResponse, error) {

	uploadFile, err := s.uploadRepo.GetUploadFileByID(fileID)
	if err != nil {
		return nil, fmt.Errorf("文件记录不存在: %v", err)
	}


	if s.uploadRepo.CheckChunkExists(fileID, chunkIndex) {

		count, _ := s.uploadRepo.GetUploadedChunksCount(fileID)
		return &response.ChunkUploadResponse{
			FileID:        fileID,
			ChunkIndex:    chunkIndex,
			ChunkUploaded: int(count),
			ChunkTotal:    uploadFile.ChunkTotal,
			IsCompleted:   int(count) == uploadFile.ChunkTotal,
		}, nil
	}


	chunkInfo, err := s.uploadRepo.GetChunkInfo(fileID, chunkIndex)
	if err != nil {
		return nil, fmt.Errorf("分片信息不存在: %v", err)
	}


	src, err := chunk.Open()
	if err != nil {
		return nil, fmt.Errorf("无法打开分片文件: %v", err)
	}
	defer src.Close()


	if md5Hash != "" {
		calculatedMD5, err := upload.CalculateChunkMD5(src)
		if err != nil {
			return nil, fmt.Errorf("计算分片MD5失败: %v", err)
		}
		if calculatedMD5 != md5Hash {
			return nil, fmt.Errorf("分片MD5校验失败")
		}
		src.Seek(0, 0)
	}


	if err := s.storage.SaveChunk(chunkInfo.ChunkPath, src); err != nil {
		return nil, fmt.Errorf("保存分片失败: %v", err)
	}


	chunkInfo.ChunkSize = chunk.Size
	chunkInfo.MD5Hash = md5Hash
	chunkInfo.IsUploaded = true

	if err := s.uploadRepo.UpdateChunkStatus(fileID, chunkIndex, true); err != nil {
		return nil, fmt.Errorf("更新分片状态失败: %v", err)
	}


	count, err := s.uploadRepo.GetUploadedChunksCount(fileID)
	if err != nil {
		return nil, fmt.Errorf("获取上传进度失败: %v", err)
	}

	if err := s.uploadRepo.UpdateUploadProgress(fileID, int(count)); err != nil {
		return nil, fmt.Errorf("更新上传进度失败: %v", err)
	}

	logger.Info("分片上传成功", "fileID", fileID, "chunkIndex", chunkIndex, "progress", fmt.Sprintf("%d/%d", count, uploadFile.ChunkTotal))

	return &response.ChunkUploadResponse{
		FileID:        fileID,
		ChunkIndex:    chunkIndex,
		ChunkUploaded: int(count),
		ChunkTotal:    uploadFile.ChunkTotal,
		IsCompleted:   int(count) == uploadFile.ChunkTotal,
	}, nil
}


func (s *UploadService) MergeChunks(fileID string) (*response.ChunkMergeResponse, error) {

	uploadFile, err := s.uploadRepo.GetUploadFileByID(fileID)
	if err != nil {
		return nil, fmt.Errorf("文件记录不存在: %v", err)
	}


	count, err := s.uploadRepo.GetUploadedChunksCount(fileID)
	if err != nil {
		return nil, fmt.Errorf("获取上传进度失败: %v", err)
	}

	if int(count) != uploadFile.ChunkTotal {
		return nil, fmt.Errorf("分片尚未完全上传，进度: %d/%d", count, uploadFile.ChunkTotal)
	}


	chunks, err := s.uploadRepo.GetFileChunks(fileID)
	if err != nil {
		return nil, fmt.Errorf("获取分片信息失败: %v", err)
	}


	var chunkPaths []string
	for _, chunk := range chunks {
		chunkPaths = append(chunkPaths, chunk.ChunkPath)
	}


	targetPath := filepath.Join(s.config.UploadDir, uploadFile.StoredName)
	if err := s.storage.MergeChunks(chunkPaths, targetPath); err != nil {
		return nil, fmt.Errorf("合并分片失败: %v", err)
	}


	now := time.Now()
	uploadFile.FilePath = targetPath
	uploadFile.UploadStatus = constants.UploadStatusCompleted
	uploadFile.UploadedAt = &now

	if err := s.uploadRepo.UpdateUploadFile(uploadFile); err != nil {
		return nil, fmt.Errorf("更新文件记录失败: %v", err)
	}


	if err := s.uploadRepo.DeleteFileChunks(fileID); err != nil {
		logger.Error("清理分片记录失败", "fileID", fileID, "error", err)
	}

	logger.Info("分片合并成功", "fileID", fileID, "filename", uploadFile.Filename)

	return &response.ChunkMergeResponse{
		FileID:     uploadFile.ID.String(),
		Filename:   uploadFile.Filename,
		StoredName: uploadFile.StoredName,
		FileSize:   uploadFile.FileSize,
		MimeType:   uploadFile.MimeType,
		Extension:  uploadFile.Extension,
		MD5Hash:    uploadFile.MD5Hash,
		FilePath:   fmt.Sprintf("/files/preview/%s/%s", uploadFile.ID.String(), uploadFile.Filename),
		UploadedAt: *uploadFile.UploadedAt,
	}, nil
}


func (s *UploadService) GetUploadProgress(fileID string) (*response.UploadProgressResponse, error) {
	uploadFile, err := s.uploadRepo.GetUploadFileByID(fileID)
	if err != nil {
		return nil, fmt.Errorf("文件记录不存在: %v", err)
	}

	progress := float64(uploadFile.ChunkUploaded) / float64(uploadFile.ChunkTotal) * 100

	return &response.UploadProgressResponse{
		FileID:        uploadFile.ID.String(),
		Filename:      uploadFile.Filename,
		FileSize:      uploadFile.FileSize,
		ChunkTotal:    uploadFile.ChunkTotal,
		ChunkUploaded: uploadFile.ChunkUploaded,
		Progress:      progress,
		Status:        uploadFile.UploadStatus,
	}, nil
}


func (s *UploadService) validateFile(file *multipart.FileHeader) error {

	if err := upload.ValidateFilename(file.Filename); err != nil {
		return err
	}


	if !s.config.ValidateFileSize(file.Size) {
		return fmt.Errorf("文件大小超出限制，最大允许 %d 字节", s.config.MaxFileSize)
	}


	mimeType := upload.GetMimeTypeFromExtension(file.Filename)
	if !s.config.ValidateMimeType(mimeType) {
		return fmt.Errorf("不支持的文件类型: %s", mimeType)
	}


	if !upload.ValidateFileExtension(file.Filename) {
		return fmt.Errorf("不支持的文件扩展名: %s", upload.GetFileExtension(file.Filename))
	}

	return nil
}

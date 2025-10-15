package upload

import (
	"strconv"

	"auto-forge/internal/dto/request"
	"auto-forge/internal/middleware"
	uploadService "auto-forge/internal/services/upload"
	"auto-forge/pkg/common"
	"auto-forge/pkg/upload"

	"github.com/gin-gonic/gin"
)












func SimpleUpload(c *gin.Context) {

	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		common.Unauthorized(c, err.Error())
		return
	}


	file, err := c.FormFile("file")
	if err != nil {
		common.BadRequest(c, "未找到上传文件")
		return
	}


	var req request.SimpleUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		common.BadRequest(c, "参数错误")
		return
	}


	result, err := uploadService.SimpleUpload(file, user.UserID)
	if err != nil {
		common.ServerError(c, err.Error())
		return
	}


	common.Success(c, result, "文件上传成功")
}











func InitChunkUpload(c *gin.Context) {

	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		common.Unauthorized(c, err.Error())
		return
	}


	req, err := common.ValidateRequest[request.ChunkUploadInitRequest](c)
	if err != nil {
		common.BadRequest(c, err.Error())
		return
	}


	result, err := uploadService.InitChunkUpload(
		req.Filename,
		req.FileSize,
		req.MD5Hash,
		req.ChunkSize,
		user.UserID,
	)
	if err != nil {
		common.ServerError(c, err.Error())
		return
	}


	common.Success(c, result, "分片上传初始化成功")
}














func UploadChunk(c *gin.Context) {

	_, err := middleware.GetUserFromContext(c)
	if err != nil {
		common.Unauthorized(c, err.Error())
		return
	}


	fileID := c.PostForm("fileID")
	chunkIndexStr := c.PostForm("chunkIndex")
	md5Hash := c.PostForm("md5Hash")

	if fileID == "" || chunkIndexStr == "" {
		common.BadRequest(c, "缺少必要参数")
		return
	}

	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil {
		common.BadRequest(c, "分片索引格式错误")
		return
	}


	chunk, err := c.FormFile("chunk")
	if err != nil {
		common.BadRequest(c, "未找到分片文件")
		return
	}


	result, err := uploadService.UploadChunk(fileID, chunkIndex, md5Hash, chunk)
	if err != nil {
		common.ServerError(c, err.Error())
		return
	}


	common.Success(c, result, "分片上传成功")
}











func MergeChunks(c *gin.Context) {

	_, err := middleware.GetUserFromContext(c)
	if err != nil {
		common.Unauthorized(c, err.Error())
		return
	}


	req, err := common.ValidateRequest[request.ChunkMergeRequest](c)
	if err != nil {
		common.BadRequest(c, err.Error())
		return
	}


	result, err := uploadService.MergeChunks(req.FileID)
	if err != nil {
		common.ServerError(c, err.Error())
		return
	}


	common.Success(c, result, "分片合并成功")
}











func GetUploadProgress(c *gin.Context) {

	_, err := middleware.GetUserFromContext(c)
	if err != nil {
		common.Unauthorized(c, err.Error())
		return
	}


	fileID := c.Param("fileID")
	if fileID == "" {
		common.BadRequest(c, "文件ID不能为空")
		return
	}


	result, err := uploadService.GetUploadProgress(fileID)
	if err != nil {
		common.ServerError(c, err.Error())
		return
	}


	common.Success(c, result, "获取上传进度成功")
}









func GetUploadConfig(c *gin.Context) {

	config := upload.NewDefaultConfig()


	result := map[string]interface{}{
		"maxFileSize":      config.MaxFileSize,
		"allowedMimeTypes": config.AllowedMimeTypes,
		"chunkSize":        config.ChunkSize,
	}

	common.Success(c, result, "获取上传配置成功")
}

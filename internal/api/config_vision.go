package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/photoprism/photoprism/internal/ai/vision"
	"github.com/photoprism/photoprism/internal/auth/acl"
	"github.com/photoprism/photoprism/internal/photoprism/get"
	"github.com/photoprism/photoprism/pkg/i18n"
)

// GetVisionConfig returns backend vision config.
//
//	@Summary	returns backend vision config
//	@Id			GetVisionConfig
//	@Tags		Config, Settings
//	@Produce	json
//	@Success	200			{object}	vision.ConfigValues
//	@Failure	401,403,429	{object}	i18n.Response
//	@Router		/api/v1/config/vision [get]
func GetVisionConfig(router *gin.RouterGroup) {
	router.GET("/config/vision", func(c *gin.Context) {
		s := Auth(c, acl.ResourceConfig, acl.AccessAll)
		conf := get.Config()

		if s.Invalid() || conf.Public() || conf.DisableSettings() {
			AbortForbidden(c)
			return
		}

		c.JSON(http.StatusOK, vision.Config)
	})
}

// SaveVisionConfig updates backend vision config.
//
//	@Summary	updates backend vision config
//	@Id			SaveVisionConfig
//	@Tags		Config, Settings
//	@Accept		json
//	@Produce	json
//	@Success	200					{object}	vision.ConfigValues
//	@Failure	400,401,403,429,500	{object}	i18n.Response
//	@Router		/api/v1/config/vision [post]
func SaveVisionConfig(router *gin.RouterGroup) {
	router.POST("/config/vision", func(c *gin.Context) {
		s := Auth(c, acl.ResourceConfig, acl.ActionManage)
		conf := get.Config()

		if s.Invalid() || conf.Public() || conf.DisableSettings() {
			AbortForbidden(c)
			return
		}

		var newConfig vision.ConfigValues

		LimitRequestBodyBytes(c, MaxSettingsRequestBytes)

		if err := c.BindJSON(&newConfig); err != nil {
			if IsRequestBodyTooLarge(err) {
				AbortRequestTooLarge(c, i18n.ErrBadRequest)
				return
			}
			AbortBadRequest(c, err)
			return
		}

		vision.Config.Models = newConfig.Models
		if err := vision.Config.Save(conf.VisionYaml()); err != nil {
			log.Errorf("config: failed saving vision config (%s)", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		// Reload to apply defaults
		vision.Config.Load(conf.VisionYaml())

		c.JSON(http.StatusOK, vision.Config)
	})
}

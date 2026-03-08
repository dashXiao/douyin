package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
	"tiktok/pkg/utils"
	"tiktok/services/api/biz/pack"
)

func AuthToken() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token, ok := c.GetQuery("token")

		if !ok {
			token, ok = c.GetPostForm("token")
		}

		if !ok {
			pack.SendFailResponse(c, errno.AuthorizationFailedError)
			c.Abort()
			return
		}

		claims, err := utils.CheckToken(token)

		if err != nil {
			pack.SendFailResponse(c, errno.AuthorizationFailedError)
			c.Abort()
			return
		}

		if claims.UserId < constants.StartID {
			pack.SendFailResponse(c, errno.AuthorizationFailedError)
			c.Abort()
		}
		// c.Set("current_user_id", claims.UserId)

		c.Next(ctx)
	}
}

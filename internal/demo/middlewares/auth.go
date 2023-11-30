package middlewares

//func Auth(ctx *gin.Context) {
//	tokenStr := ctx.Request.Header.Get("token")
//	osType := ctx.Request.Header.Get("os")
//	uInfo, err := token.CheckToken(tokenStr, osType)
//	if err != nil {
//		ctx.Abort()
//		response.ResPubErr(ctx, err)
//		return
//	}
//	ctx.Set("uid", uInfo.ID)
//}

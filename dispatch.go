package redisChan 

// 3.-------------dispatch------------------------
func waitDispatch(cr chan chanReq){
    for {
        tunnel := <- cr
        tunnel.response <- dispatch(tunnel.request)
    }
    
}

//分发请求到响应的action操作函数上去
func dispatch(req *reqBulk) (res *responseBulk) {
    if h, ok := actions[req.cmd]; ok {
        res = &responseBulk{}
        h(req, res) // 调用相应的函数
    } else {
        return notFound(req)
    }
    return
}

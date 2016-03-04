package redisChan 
type CommandCode string

//request command
const (
    GET     CommandCode = "get"
    SET     CommandCode = "set"
    DELETE  CommandCode = "delete"
    STATS   CommandCode = "stats"
)

type Status int
//response 状态
const (
    SUCCESS Status = iota
    ERROR
    STORED
    NOT_STORED
    END
    DELETED
    NOT_FOUND
    UNKNOWN_COMMAND
)

// 4.-------------action------------------------
type action func(req *reqBulk, res *responseBulk)

var actions = map[CommandCode]action{
    STATS: StatsAction,
}

//给request绑定上处理程序
func BindAction(opcode CommandCode, h action) {
    actions[opcode] = h
}

//未支持命令
func notFound(req *reqBulk) *responseBulk {
    var response responseBulk
    response.status = UNKNOWN_COMMAND
    return &response
}

var data map[string]string = make(map[string]string)

func init(){
    BindAction(GET, GetAction)
    BindAction(SET, SetAction)
}

func GetAction(req *reqBulk, res *responseBulk) {
    res.cmd = req.cmd
    res.key = req.key
    if tmp, ok := data[req.key]; ok {
         res.value = tmp
         res.status = SUCCESS
    }else{
         res.status = NOT_STORED
    }
}

func SetAction(req *reqBulk, res *responseBulk) {
    res.cmd = req.cmd
    res.key = req.key
    data[req.key] = req.value
    res.status = STORED
}

func StatsAction(req *reqBulk, res *responseBulk) {
    res.cmd = req.cmd
    res.key = req.key
    res.value = "END\r\n"
    res.status = END
}

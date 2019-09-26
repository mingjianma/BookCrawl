package model

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import (
    "strings"
    "errors"
    "BookCrawl/config"
    "log"
)

//数据库操作类型
type MysqlDb struct {
    Dbname    string
    Tblname   string
    Fields    []string
    Gfield   string
    Where     string
    Limit     string
    OrderBy   string
    GroupBy   string
    Dbcon     *sql.DB
}

func (this *MysqlDb) NewMysqlDb() error {
    //用户名 密码 IP 端口 在config.go中配置
    //构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
    path := strings.Join([]string{config.MQ_USERNAME, ":", config.MQ_PASSWD, "@tcp(", config.MQ_HOST, ":", config.MQ_PORT, ")/", this.Dbname, "?charset=utf8"}, "")

    //打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
    DB, _ := sql.Open("mysql", path)
    //设置数据库最大连接数
    DB.SetConnMaxLifetime(100)
    //设置上数据库最大闲置连接数
    DB.SetMaxIdleConns(10)

    //验证连接
    if err := DB.Ping(); err != nil{
        return err
    }

    this.Dbcon = DB
    //设置select字段默认值
    this.Gfield = "*"
    this.getFields()

    return nil
}

func (this *MysqlDb) Close(){
    this.Dbcon.Close()
}


/**
 * 获取当前表的所有字段
 */
func (this *MysqlDb) getFields(){
 
    //查看表结构
    sql := "DESC " + this.Tblname
    //执行并发送SQL
    result, err := this.Dbcon.Query(sql)
    defer result.Close() 
    if err != nil{
        log.Printf("sql fail ! [%s]",err)
    }
 
    this.Fields = make([]string,0)
 
    for result.Next() {
        var field string
        var Type interface{}
        var Null string
        var Key string
        var Default interface{}
        var Extra string
        err :=result.Scan(&field,&Type,&Null,&Key,&Default,&Extra)
        if err != nil{
            log.Printf("scan fail ! [%s]",err)
        }
        this.Fields = append(this.Fields, field)
    }
 }

/**
 * 查询单条数据
 * @return res    查询结果
 * @return err    错误信息
 */
func (this *MysqlDb) Find() (res map[string]string, err error) {
    //log.Printf("SELECT " + this.Gfield + " FROM `" + this.Tblname + "` "+ this.Where + this.OrderBy + this.GroupBy + " limit 1")
    row, err := this.Dbcon.Query("SELECT " + this.Gfield + " FROM `" + this.Tblname + "` "+ this.Where + this.OrderBy + this.GroupBy + " limit 1")
    defer row.Close()
    if err != nil {
        return 
    }

    //获取列名
    columns, err := row.Columns()
    if err != nil {
        return 
    }

    if !row.Next() {
        return 
    }
    
    //定义一个切片,长度是字段的个数,切片里面的元素类型是sql.RawBytes
    values := make([]sql.RawBytes, len(columns))
    //定义一个切片,元素类型是interface{} 接口
    scanArgs := make([]interface{}, len(values))
    for i := range values {
        //把sql.RawBytes类型的地址存进去了
        scanArgs[i] = &values[i]
    }

    tmp := make(map[string]string)
    row.Scan(scanArgs...)
    for i, col := range values {
        tmp[string(columns[i])] = string(col)
    }
    res = tmp
    row.Close()

    return 
}

/**
 * 查询多条数据
 * @return result    查询结果
 * @return err       错误信息
 */
func (this *MysqlDb) Select() (result []map[string]string, err error) {
    
    rows, err := this.Dbcon.Query("SELECT " + this.Gfield + " FROM `" + this.Tblname + "` "+ this.Where + this.OrderBy + this.GroupBy + this.Limit)
    defer rows.Close()
    if err != nil {
        return 
    }

     //获取列名
    columns, err := rows.Columns()
    if err != nil {
        return 
    }

    //定义一个切片,长度是字段的个数,切片里面的元素类型是sql.RawBytes
    values := make([]sql.RawBytes, len(columns))
    //定义一个切片,元素类型是interface{} 接口
    scanArgs := make([]interface{}, len(values))
    for i := range values {
        //把sql.RawBytes类型的地址存进去了
        scanArgs[i] = &values[i]
    }
    //获取字段值
    for rows.Next() {
        res := make(map[string]string)
        rows.Scan(scanArgs...)
        for i, col := range values {
            res[columns[i]] = string(col)
        }
        result = append(result, res)
    }

    rows.Close()

    return 
}

/**
 * 单条数据插入
 * @param  data      插入数据
 * @return insert_id 插入成功的主键id
 * @return err       错误信息
 */
func (this *MysqlDb) Insert(data map[string]string) (insert_id int, err error){
    //开启事务
    tx, err := this.Dbcon.Begin()

    if err != nil{
        return 
    }

    if data == nil {
        err = errors.New("params data empty")
        return
    }
    keys := " ("
    values := " ("
    for k, v := range data {
        keys += "`" + k + "`,"
        values += "'" + v + "',"
    }
    new_keys := strings.TrimRight(keys, ",") + ") "
    new_values := strings.TrimRight(values, ",") + ") "

    //log.Printf("INSERT INTO `" + this.Tblname + "` " + new_keys + " VALUES " + new_values)
    //准备sql语句
    stmt, err := tx.Prepare("INSERT INTO `" + this.Tblname + "` " + new_keys + " VALUES " + new_values)
    if err != nil{
        return 
    }
    //将参数传递到sql语句中并且执行
    res, err := stmt.Exec()
    if err != nil{
        return 
    }

    //将事务提交
    tx.Commit()

    //获得上一个插入自增的id
    last_id, err := res.LastInsertId()
    insert_id = int(last_id)
    return 
}


/**
 * 设置要查询的字段信息
 * @param string field    要查询的字段
 * @return this           返回自己，保证连贯操作
 */
func (this *MysqlDb) SetField(field string) *MysqlDb{
    this.Gfield = field
    return this
}

/**
 * where条件
 * @param string where   输入的where条件
 * @return this          返回自己，保证连贯操作
 */
func (this *MysqlDb) SetWhere(data map[string]string) *MysqlDb{
    var where_str = " WHERE 1 "
    if data != nil {
        for k, v := range data {
            where_str += " AND `" + k + "` = '" + v +"'"
        }
    }

   this.Where = where_str

   return this
}

/**
 * limit限制
 * @param  beigin        起始位置
 * @param  offset        偏移位置
 * @return this          返回自己，保证连贯操作
 */
func (this *MysqlDb) SetLimit(beigin int, offset int) *MysqlDb{
   this.Limit = " LIMIT " + string(beigin) + ", " + string(offset) + " "
   return this
}

/**
 * order by条件
 * @param  string field  输入的order by条件
 * @return this          返回自己，保证连贯操作
 */
func (this *MysqlDb) SetOrderBy(field string) *MysqlDb{
   this.OrderBy = field
   return this
}

/**
 * group by条件
 * @param string field  输入的group by条件
 * @return this         返回自己，保证连贯操作
 */
func (this *MysqlDb) SetGroupBy(field string) *MysqlDb{
   this.GroupBy = field
   return this
}

/**
 * 修改操作
 * @param  array $data  要修改的数组
 * @return bool 修改成功返回true，失败返回false
 */
func (this *MysqlDb) Update(data map[string]string ) (affect_num int64, err error){
    str := ""
    //过滤非法字段
    for k,v:=range data{
        if res:=in_array(k, this.Fields); res != true {
            delete(data, k)
        }else{
            str += k +` = '` + v + `',`
        }
    }
 
    //去掉最右侧的逗号
    str =strings.TrimRight(str, ",")
    
    //log.Printf(`update ` + this.Tblname + ` set ` + str  +` `+ this.Where)
    sql := `update ` + this.Tblname + ` set ` + str  +` `+ this.Where
    res, err := this.Dbcon.Exec(sql)
    if err != nil {
        return 
    }

    affect_num, err = res.RowsAffected()

    return 
}


//是否存在数组内
func in_array(need interface{}, needArr []string) bool {
    for _,v := range needArr{
        if need == v{
            return true
        }
    }
    return false
}


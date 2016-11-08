package main

import (
    	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main(){
	if db,error := sql.Open("mysql","root:12345@tcp(127.0.0.1:3306)/usageinfo?charset=utf8");error!=nil{
		fmt.Println(error.Error())
	}else {
		if tx,error:=db.Begin();error!=nil{
			fmt.Println(error.Error())
		}else{
			if stmt,error:=tx.Prepare("update User set userid=? where id=?");error!=nil{
				fmt.Println(error.Error())
			}else{
				if obj,error:=stmt.Exec("testuid",1);error!=nil{
					fmt.Println(error.Error())
					defer tx.Rollback()
				}else{
					fmt.Println(obj.RowsAffected())
				}
			}
			//i:=1/0
			//fmt.Print(i)
			if rows,error:=tx.Query("select * from User");error!=nil{
				fmt.Println(error.Error())
			}else{
				for rows.Next(){
					var row User
					if error:=rows.Scan(&row.id, &row.userid,&row.name, &row.password, &row.company, &row.department, &row.email, &row.create_time, &row.update_time);error!=nil{
						fmt.Println(error.Error())
					}else{
						fmt.Println(row)
					}
				}
			}
			tx.Commit()
		}
		db.Close()
	}

}

type User struct  {
	id  int
	userid string
	name   string
	password string
	company  string
	department string
	email string
	create_time []uint8
	update_time []uint8
}
package main

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

type Alpha struct {
	Id               int
	StringIdentifier string
}

type Dossier struct {
	Id           int
	DepartmentId int
	EmployeeId   int
	Summary      string
}

func (d *Dossier) TableName() string {
	return "dossier"
}

func init() {
	orm.RegisterModel(new(Alpha), new(Dossier))
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(192.168.137.3:3306)/yiitest")
}

func main() {
	// 创建表
	//orm.RunSyncdb("default", false, true)
	o := orm.NewOrm()

	// crud
	alpha := &Alpha{Id: 2}
	err := o.Read(alpha)
	if err != nil {
		panic(err)
	}
	fmt.Printf("id: %d, string_identifier: %s\n", alpha.Id, alpha.StringIdentifier)

	// query builder
	qb, _ := orm.NewQueryBuilder("mysql")
	// 需要严格按照sql顺序拼接
	qb.Select("*").
		From(new(Dossier).TableName()).
		Where("id > ?").
		Limit(10)
	sql := qb.String()
	var dossiers []*Dossier
	_, err = o.Raw(sql, 1).QueryRows(&dossiers)
	if err != nil {
		panic(err)
	}
	for _, dossier := range dossiers {
		fmt.Printf("id: %d, department_id: %d, employee_id: %d, summary: %s\n",
			dossier.Id, dossier.DepartmentId, dossier.EmployeeId, dossier.Summary)
	}

	// query seter
	dossiers = []*Dossier{}
	_, err = o.QueryTable(new(Dossier)).Filter("id__gt", 1).All(&dossiers)
	if err != nil {
		panic(err)
	}
	for _, dossier := range dossiers {
		fmt.Printf("id: %d, department_id: %d, employee_id: %d, summary: %s\n",
			dossier.Id, dossier.DepartmentId, dossier.EmployeeId, dossier.Summary)
	}
}

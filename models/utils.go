package models

//添加一个企业信息.
func DefaultAdd(o interface{})error{
	return g.Create(o).Error

}

func DefaultRemoveById(id int,o interface{})error{
	var err error
	tx:=g.Begin()
	defer func() {
		if err!=nil{
			tx.Rollback()
		}else {
			tx.Commit()
		}
	}()
	if err=tx.Where("id=?",id).Delete(o).Error;err!=nil{
		return err
	}

	return nil

}
func DefaultUpdateById(id int,o interface{})error{
	var err error
	tx:=g.Begin()
	defer func() {
		if err!=nil{
			tx.Rollback()
		}else {
			tx.Commit()
		}
	}()
	if err=tx.Where("id=?",id).Save(o).Error;err!=nil{
		return err
	}

	return nil

}
func DefaultGetByOrgId(cid int,o interface{})(err error)  {

	if err:=g.Where("org_id=?",cid).Find(o).Error;err!=nil{
		return err
	}
	return nil
}
func DefaultGetById(cid int,o interface{})(err error)  {

	if err:=g.Where("id=?",cid).First(o).Error;err!=nil{
		return err
	}
	return nil
}
func DefaultGetAll(o interface{})(err error)  {

	if err:=g.Find(o).Error;err!=nil{
		return err
	}
	return nil
}
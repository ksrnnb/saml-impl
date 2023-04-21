package model

const defaultCompanyID = "38azqp4z"
const defaultCompanyName = "サンプル会社"

type Company struct {
	ID   string `gorm:"primaryKey"`
	Name string
}

func FindCompany(cid string) (*Company, error) {
	c := &Company{}
	res := db.Limit(1).Find(c, "id=?", cid)
	if err := res.Error; err != nil {
		return nil, err
	}
	return c, nil
}

func ListAllCompanies() ([]*Company, error) {
	var companies []*Company
	if err := db.Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

func (c *Company) IsZero() bool {
	return c.ID == "" && c.Name == ""
}

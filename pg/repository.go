package pg

import (
	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{Db: db}
}

type Repository struct {
	Db *gorm.DB
}

func (p *Repository) Load(id uint) (*Blog, error) {
	blog := &Blog{}
	err := p.Db.Where(`id = ?`, id).First(blog).Error
	return blog, err
}

// ListAll list all blogs
func (p *Repository) ListAll() ([]*Blog, error) {
	var l []*Blog
	err := p.Db.Find(&l).Error
	return l, err
}

func (p *Repository) List(offset, limit int) ([]*Blog, error) {
	var l []*Blog
	err := p.Db.Offset(offset).Limit(limit).Find(&l).Error
	return l, err
}

func (p *Repository) Save(blog *Blog) error {
	return p.Db.Save(blog).Error
}

func (p *Repository) Delete(id uint) error {
	return p.Db.Delete(&Blog{ID: id}).Error
}

func (p *Repository) SearchByTitle(q string, offset, limit int) ([]*Blog, error) {
	var l []*Blog
	err := p.Db.Where(`title like ?`, "%"+q+"%").
		Offset(offset).Limit(limit).Find(&l).Error
	return l, err
}

func (p *Repository) Migrate() error {
	return p.Db.AutoMigrate(&Blog{})
}

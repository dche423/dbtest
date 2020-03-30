package pg

import (
	"github.com/jinzhu/gorm"
)

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type Repository struct {
	db *gorm.DB
}

func (p *Repository) Load(id uint) (*Blog, error) {
	blog := &Blog{}
	err := p.db.Where(`id = ?`, id).First(blog).Error
	return blog, err
}

// ListAll list all blogs
func (p *Repository) ListAll() ([]*Blog, error) {
	var l []*Blog
	err := p.db.Find(&l).Error
	return l, err
}

func (p *Repository) List(offset, limit int) ([]*Blog, error) {
	var l []*Blog
	err := p.db.Offset(offset).Limit(limit).Find(&l).Error
	return l, err
}

func (p *Repository) Save(blog *Blog) error {
	return p.db.Save(blog).Error
}

func (p *Repository) Delete(id uint) error {
	return p.db.Delete(&Blog{ID: id}).Error
}

func (p *Repository) SearchByTitle(q string, offset, limit int) ([]*Blog, error) {
	var l []*Blog
	err := p.db.Where(`title like ?`, "%"+q+"%").
		Offset(offset).Limit(limit).Find(&l).Error
	return l, err
}

func (p *Repository) Migrate() error {
	return p.db.AutoMigrate(&Blog{}).Error
}

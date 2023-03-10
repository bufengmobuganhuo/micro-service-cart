package repository

import (
	"errors"
	"github.com/bufengmobuganhuo/micro-service-cart/domain/model"
	"github.com/jinzhu/gorm"
)

type ICartRepository interface {
	InitTable() error
	FindCartByID(int64) (*model.Cart, error)
	CreateCart(*model.Cart) (int64, error)
	DeleteCartByID(int64) error
	UpdateCart(*model.Cart) error
	FindAllByUserID(int642 int64) ([]model.Cart, error)

	CleanCart(int64) error
	IncrNum(int64, int64) error
	DecrNum(int64, int64) error
}

//创建cartRepository
func NewCartRepository(db *gorm.DB) ICartRepository {
	return &CartRepository{mysqlDb: db}
}

type CartRepository struct {
	mysqlDb *gorm.DB
}

// InitTable 初始化表
func (u *CartRepository) InitTable() error {
	return u.mysqlDb.CreateTable(&model.Cart{}).Error
}

// FindCartByID 根据ID查找Cart信息
func (u *CartRepository) FindCartByID(cartID int64) (cart *model.Cart, err error) {
	cart = &model.Cart{}
	return cart, u.mysqlDb.First(cart, cartID).Error
}

// CreateCart 创建Cart信息
func (u *CartRepository) CreateCart(cart *model.Cart) (int64, error) {
	// 如果记录已经存在，就不应该再插入
	db := u.mysqlDb.FirstOrCreate(cart, model.Cart{
		ProductID: cart.ProductID,
		Num:       cart.Num,
		SizeID:    cart.SizeID,
		UserID:    cart.UserID,
	})
	if db.Error != nil {
		return 0, db.Error
	}
	// 已经存在了
	if db.RowsAffected == 0 {
		return 0, errors.New("购物车插入失败")
	}
	return cart.ID, nil
}

// DeleteCartByID 根据ID删除Cart信息
func (u *CartRepository) DeleteCartByID(cartID int64) error {
	return u.mysqlDb.Where("id = ?", cartID).Delete(&model.Cart{}).Error
}

// UpdateCart 更新Cart信息
func (u *CartRepository) UpdateCart(cart *model.Cart) error {
	return u.mysqlDb.Model(cart).Update(cart).Error
}

// FindAll 获取结果集
func (u *CartRepository) FindAllByUserID(userId int64) (cartAll []model.Cart, err error) {
	return cartAll, u.mysqlDb.Where("user_id = ?", userId).Find(&cartAll).Error
}

func (u *CartRepository) CleanCart(userId int64) error {
	return u.mysqlDb.Where("user_id = ?", userId).Delete(&model.Cart{}).Error
}

func (u *CartRepository) IncrNum(cartID int64, num int64) error {
	cart := &model.Cart{ID: cartID}
	return u.mysqlDb.Model(cart).UpdateColumn("num", gorm.Expr("num + ?", num)).Error
}

func (u *CartRepository) DecrNum(cartId int64, num int64) error {
	cart := &model.Cart{ID: cartId}
	// 已有数量必须>要减少的数量
	db := u.mysqlDb.Model(cart).Where("num >= ?", num).UpdateColumn("num", gorm.Expr("num - ?", num))
	if db.Error != nil {
		return db.Error
	}
	// 已有数量必须>要减少的数量
	if db.RowsAffected == 0 {
		return errors.New("减少失败")
	}
	return nil
}

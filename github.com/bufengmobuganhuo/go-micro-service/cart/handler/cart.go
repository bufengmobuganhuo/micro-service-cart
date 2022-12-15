package handler

import (
	"context"
	"github.com/bufengmobuganhuo/go-micro-service/cart/domain/model"
	"github.com/bufengmobuganhuo/go-micro-service/cart/domain/service"
	cart "github.com/bufengmobuganhuo/go-micro-service/cart/proto/cart"
	common "github.com/bufengmobuganhuo/micro-service-common"
)

type Cart struct {
	CartDataService service.ICartDataService
}

func (c Cart) AddCart(ctx context.Context, req *cart.CartInfo, resp *cart.ResponseAdd) error {
	cart := &model.Cart{}
	err := common.SwapTo(req, cart)
	if err != nil {
		return err
	}
	_, err = c.CartDataService.AddCart(cart)
	if err != nil {
		return err
	}
	resp.CartId = cart.ID
	resp.Msg = "插入成功"
	return err
}

func (c Cart) CleanCart(ctx context.Context, req *cart.Clean, resp *cart.Response) error {
	err := c.CartDataService.CleanCart(req.UserId)
	if err != nil {
		return err
	}
	resp.Msg = "清空购物车成功"
	return nil
}

func (c Cart) Incr(ctx context.Context, req *cart.Item, resp *cart.Response) error {
	err := c.CartDataService.IncrNum(req.Id, req.ChangeNum)
	if err != nil {
		return err
	}
	resp.Msg = "增加数量成功"
	return nil
}

func (c Cart) Decr(ctx context.Context, req *cart.Item, resp *cart.Response) error {
	err := c.CartDataService.DecrNum(req.Id, req.ChangeNum)
	if err != nil {
		return err
	}
	resp.Msg = "减少数量成功"
	return nil
}

func (c Cart) DeleteItemID(ctx context.Context, req *cart.CartID, resp *cart.Response) error {
	err := c.CartDataService.DeleteCart(req.Id)
	if err != nil {
		return err
	}
	resp.Msg = "删除商品成功"
	return nil
}

func (c Cart) GetAll(ctx context.Context, req *cart.CartFindAll, resp *cart.CartAll) error {
	carts, err := c.CartDataService.FindAllCart(req.UserId)
	if err != nil {
		return err
	}
	infos := []*cart.CartInfo{}
	for _, crt := range carts {
		info := &cart.CartInfo{}
		e := common.SwapTo(crt, info)
		if e != nil {
			continue
		}
		infos = append(infos, info)
	}
	resp.CartInfo = infos
	return nil
}

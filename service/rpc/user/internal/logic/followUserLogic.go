package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"h68u-tiktok-app-microservice/common/error/rpcErr"
	"h68u-tiktok-app-microservice/common/model"
	"h68u-tiktok-app-microservice/common/mq"
	"h68u-tiktok-app-microservice/common/utils"
	"h68u-tiktok-app-microservice/service/rpc/user/internal/svc"
	"h68u-tiktok-app-microservice/service/rpc/user/types/user"
)

type FollowUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowUserLogic {
	return &FollowUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowUserLogic) FollowUser(in *user.FollowUserRequest) (*user.Empty, error) {
	err := l.svcCtx.DBList.Mysql.Transaction(func(tx *gorm.DB) error {
		var users *model.User
		var followUser *model.User
		tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", in.UserId).First(&users)
		tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", in.FollowUserId).First(&followUser)

		//处理关注用户
		result, err := l.svcCtx.DBList.Redis.Exists(l.ctx, utils.GenUserInfoCacheKey(in.UserId)).Result()
		if result == 1 {
			// 如果是大V（在redis中有缓存），就只更新缓存，交给定时任务更新数据库
			task, err := mq.NewAddCacheValueTask(utils.GenUserInfoCacheKey(in.UserId), "FollowCount", 1)
			if err != nil {
				logx.WithContext(l.ctx).Errorf("创建任务失败: %v", err)
				return err
			}
			if _, err := l.svcCtx.AsynqClient.Enqueue(task); err != nil {
				logx.WithContext(l.ctx).Errorf("发送任务失败: %v", err)
				return err
			}
		} else {
			if err != nil {
				l.Logger.Error(rpcErr.CacheError.Code, err.Error())
			}
			// 如果是普通用户，就直接更新数据库
			users.FollowCount++
			err := tx.Save(&users).Error
			if err != nil {
				return status.Error(rpcErr.DataBaseError.Code, err.Error())
			}
		}

		//处理被关注用户
		result, err = l.svcCtx.DBList.Redis.Exists(l.ctx, utils.GenUserInfoCacheKey(in.FollowUserId)).Result()
		if result == 1 {
			// 如果是大V（在redis中有缓存），就只更新缓存，交给定时任务更新数据库
			task, err := mq.NewAddCacheValueTask(utils.GenUserInfoCacheKey(in.FollowUserId), "FanCount", 1)
			if err != nil {
				logx.WithContext(l.ctx).Errorf("创建任务失败: %v", err)
				return err
			}
			if _, err := l.svcCtx.AsynqClient.Enqueue(task); err != nil {
				logx.WithContext(l.ctx).Errorf("发送任务失败: %v", err)
				return err
			}
		} else {
			if err != nil {
				l.Logger.Error(rpcErr.CacheError.Code, err.Error())
			}
			// 如果是普通用户，就直接更新数据库

			followUser.FanCount++
			err = tx.Save(&followUser).Error
			if err != nil {
				return status.Error(rpcErr.DataBaseError.Code, err.Error())
			}
		}
		//建立关注关系
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(users).Association("Follows").Append(followUser)
		if err != nil {
			return status.Error(rpcErr.DataBaseError.Code, err.Error())
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}
	return &user.Empty{}, nil
}

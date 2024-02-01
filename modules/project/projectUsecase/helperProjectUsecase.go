package projectusecase

import (
	datacenterPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterPb"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users"
	usersPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersPb"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"time"
)

func (u *projectUsecase) assignProjectRes(data *datacenterPb.ProjectRes, fav bool, parsedTime time.Time) *project.ProjectRes {
	return &project.ProjectRes{
		Id:             data.Id,
		Name:           data.Name,
		LogoUrl:        data.LogoUrl,
		BannerUrl:      data.BannerUrl,
		WebsiteUrl:     data.WebsiteUrl,
		CryptoCategory: data.CryptoCategory,
		Description:    data.Description,
		Reason:         data.Reason,
		Category:       data.Category,
		Contact:        data.Contact,
		FavCount:       data.FavCount,
		CommentCount:   data.CommentCount,
		CreatedBy:      data.CreatedBy,
		CreatedAt:      parsedTime,
		FavOrNot:       fav,
	}
}

func (u *projectUsecase) convertPModelToPWithUser(m *project.ProjectModel, us *usersPb.GetUserInfoRes) (*project.ProjectResWithUser, error) {

	loc, err := utils.LocationTime()
	if err != nil {
		return nil, err
	}
	return &project.ProjectResWithUser{
		Id:             m.Id.Hex(),
		Name:           m.Name,
		LogoUrl:        m.LogoUrl,
		BannerUrl:      m.BannerUrl,
		WebsiteUrl:     m.WebsiteUrl,
		CryptoCategory: m.CryptoCategory,
		Description:    m.Description,
		Reason:         m.Reason,
		Category:       m.Category,
		Contact:        m.Contact,
		FavCount:       m.FavCount,
		CommentCount:   m.CommentCount,
		User: users.UserProfileRes{
			Id:        us.UserId,
			Email:     us.Email,
			Profile:   us.Profile,
			Username:  us.UserName,
			Firstname: us.FirstName,
			Lastname:  us.LastName,
		},
		CreateAt:  m.CreateAt.In(loc),
		UpdatedAt: m.UpdatedAt.In(loc),
	}, nil
}

func (u *projectUsecase) convertPDatacenterToPWithUser(m *datacenterPb.GetSingleProjectDataCenterRes, us *usersPb.GetUserInfoRes, fav bool) (*project.ProjectResWithUser, error) {

	loc, err := utils.LocationTime()
	if err != nil {
		return nil, err
	}
	return &project.ProjectResWithUser{
		Id:             m.Projects.Id,
		Name:           m.Projects.Name,
		LogoUrl:        m.Projects.LogoUrl,
		BannerUrl:      m.Projects.BannerUrl,
		WebsiteUrl:     m.Projects.WebsiteUrl,
		CryptoCategory: m.Projects.CryptoCategory,
		Description:    m.Projects.Description,
		Reason:         m.Projects.Reason,
		Category:       m.Projects.Category,
		Contact:        m.Projects.Contact,
		FavCount:       m.Projects.FavCount,
		CommentCount:   m.Projects.CommentCount,
		FavOrNot:       fav,
		User: users.UserProfileRes{
			Id:        us.UserId,
			Email:     us.Email,
			Profile:   us.Profile,
			Username:  us.UserName,
			Firstname: us.FirstName,
			Lastname:  us.LastName,
		},
		CreateAt:  utils.ConvertStringTimeToTime(m.Projects.CreatedAt).In(loc),
		UpdatedAt: utils.ConvertStringTimeToTime(m.Projects.UpdatedAt).In(loc),
	}, nil
}

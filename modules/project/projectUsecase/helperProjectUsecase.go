package projectusecase

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment"
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
		Owner: users.UserProfileRes{
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

func (u *projectUsecase) convertPDatacenterToPWithUser(m *datacenterPb.GetSingleProjectDataCenterRes, fav bool, rawC *comment.CommentModel, rawU []*usersPb.UserProfile) (*project.FullProjectRes, error) {

	owner := new(users.UserProfileRes)
	comments := make([]comment.CommentAResWithUser, 0)

	for _, v := range rawU {
		if v.UserId == m.Projects.CreatedBy {
			owner.Id = v.UserId
			owner.Email = v.Email
			owner.Firstname = v.FirstName
			owner.Lastname = v.LastName
			owner.Username = v.UserName
			owner.Profile = v.Profile
		}
	}

	for _, c := range rawC.Comments {
		for _, u := range rawU {
			if c.CreatedBy == u.UserId {
				comments = append(comments, comment.CommentAResWithUser{
					CommentId: c.CommentId.Hex(),
					Title:     c.Title,
					Content:   c.Content,
					CreateAt:  c.CreateAt,
					UpdatedAt: c.UpdatedAt,
					CreatedBy: users.UserProfileRes{
						Id:        u.UserId,
						Email:     u.Email,
						Profile:   u.Profile,
						Username:  u.UserName,
						Firstname: u.FirstName,
						Lastname:  u.LastName,
					},
				})
			}
		}
	}

	loc, err := utils.LocationTime()
	if err != nil {
		return nil, err
	}
	return &project.FullProjectRes{
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
		Comment:        comments,
		FavOrNot:       fav,
		Owner:          *owner,
		CreateAt:       utils.ConvertStringTimeToTime(m.Projects.CreatedAt).In(loc),
		UpdatedAt:      utils.ConvertStringTimeToTime(m.Projects.UpdatedAt).In(loc),
	}, nil

}

func (u *projectUsecase) accumateUserId(rawC *comment.CommentModel) []string {

	createdBy := make([]string, 0)

	for _, v := range rawC.Comments {
		createdBy = append(createdBy, v.CreatedBy)
	}
	return createdBy
}

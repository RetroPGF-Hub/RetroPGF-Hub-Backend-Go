package projectusecase

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users"
	usersPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersPb"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"time"
)

func (u *projectUsecase) assignProjectRes(data *project.ProjectModel, fav bool, parsedTime time.Time, user *users.UserProfileRes) *project.ProjectResWithUser {
	return &project.ProjectResWithUser{
		Id:           data.Id.Hex(),
		Name:         data.Name,
		LogoUrl:      data.LogoUrl,
		GithubUrl:    data.GithubUrl,
		WebsiteUrl:   data.WebsiteUrl,
		Description:  data.Description,
		Feedback:     data.Feedback,
		Category:     data.Category,
		FavCount:     data.FavCount,
		CommentCount: data.CommentCount,
		Owner: users.UserProfileRes{
			Id:        user.Id,
			Email:     user.Email,
			Profile:   user.Profile,
			Username:  user.Username,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
		},
		CreatedAt: parsedTime,
		FavOrNot:  fav,
	}
}

func (u *projectUsecase) convertPModelToPWithUser(m *project.ProjectModel, us *usersPb.GetUserInfoRes) (*project.ProjectResWithUser, error) {

	loc, err := utils.LocationTime()
	if err != nil {
		return nil, err
	}
	return &project.ProjectResWithUser{
		Id:           m.Id.Hex(),
		Name:         m.Name,
		LogoUrl:      m.LogoUrl,
		GithubUrl:    m.GithubUrl,
		WebsiteUrl:   m.WebsiteUrl,
		Description:  m.Description,
		Feedback:     m.Feedback,
		Category:     m.Category,
		FavCount:     m.FavCount,
		CommentCount: m.CommentCount,
		Owner: users.UserProfileRes{
			Id:        us.UserId,
			Email:     us.Email,
			Profile:   us.Profile,
			Username:  us.UserName,
			Firstname: us.FirstName,
			Lastname:  us.LastName,
		},
		CreatedAt: m.CreateAt.In(loc),
		UpdatedAt: m.UpdatedAt.In(loc),
	}, nil
}

func (u *projectUsecase) convertPDatacenterToPWithUser(m *project.ProjectModel, fav bool, rawC *comment.CommentProjectModel, rawU []*usersPb.UserProfile) (*project.FullProjectRes, error) {

	owner := new(users.UserProfileRes)
	comments := make([]comment.CommentAResWithUser, 0)

	for _, v := range rawU {
		if v.UserId == m.CreatedBy {
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
		Id:           m.Id.Hex(),
		Name:         m.Name,
		LogoUrl:      m.LogoUrl,
		GithubUrl:    m.GithubUrl,
		WebsiteUrl:   m.WebsiteUrl,
		Description:  m.Description,
		Feedback:     m.Feedback,
		Category:     m.Category,
		FavCount:     m.FavCount,
		CommentCount: m.CommentCount,
		Comment:      comments,
		FavOrNot:     fav,
		Owner:        *owner,
		CreateAt:     m.CreateAt.In(loc),
		UpdatedAt:    m.UpdatedAt.In(loc),
	}, nil

}

func (u *projectUsecase) accumateUserId(rawC *comment.CommentProjectModel) []string {

	createdBy := make([]string, 0)

	for _, v := range rawC.Comments {
		createdBy = append(createdBy, v.CreatedBy)
	}
	return createdBy
}

func (u *projectUsecase) accumateUserIdByProjects(rawP []*project.ProjectModel) []string {
	createdBy := make([]string, 0)
	for _, v := range rawP {
		createdBy = append(createdBy, v.CreatedBy)
	}
	return createdBy
}

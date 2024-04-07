package projectusecase

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/comment"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/favorite"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"
	"fmt"
	"sync"

	usersPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/users/usersPb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ProjectUsecaseService interface {
		CreateNewProjectUsecase(pctx context.Context, grpcCfg *config.Grpc, req *project.InsertProjectReq) (*project.ProjectRes, error)
		FindOneProjectUsecase(pctx context.Context, grpcCfg *config.Grpc, projectId, userId string) (*project.FullProjectRes, error)
		UpdateOneProjectUsecase(pctx context.Context, grpcCfg *config.Grpc, projectId, userId string, req *project.InsertProjectReq) (*project.ProjectResWithUser, error)
		DeleteOneProjectUsecase(pctx context.Context, projectId, userId string) error
		FindAllProjectUsecase(pctx context.Context, grpcCfg *config.Grpc, limit, skip, pageCount int64, userId, sort, category, search, projectType string) ([]*project.ProjectResWithUser, int64, error)
		CreateNewQuestionUsecase(pctx context.Context, grpcCfg *config.Grpc, req *project.InsertQuestionReq) (*project.ProjectRes, error)
		GetNewestProject(pctx context.Context, limit int64, exceptProjectId string) ([]*project.RandomProjectDisplay, error)
	}

	projectUsecase struct {
		pActor modules.ProjectSvcInteractor
	}
)

func NewProjectUsecase(pActor modules.ProjectSvcInteractor) ProjectUsecaseService {
	return &projectUsecase{
		pActor: pActor,
	}
}

func (u *projectUsecase) CreateNewProjectUsecase(pctx context.Context, grpcCfg *config.Grpc, req *project.InsertProjectReq) (*project.ProjectRes, error) {
	projectId, err := u.pActor.ProjectRepo.InsertOneProject(pctx, &project.ProjectModel{
		Id:           primitive.NewObjectID(),
		Name:         req.Name,
		LogoUrl:      req.LogoUrl,
		Type:         req.Type,
		FavCount:     0,
		CommentCount: 0,
		GithubUrl:    req.GithubUrl,
		WebsiteUrl:   req.WebsiteUrl,
		Description:  req.Description,
		Feedback:     req.Feedback,
		Category:     req.Category,
		CreatedBy:    req.CreatedBy,
		CreateAt:     utils.LocalTime(),
		UpdatedAt:    utils.LocalTime(),
	})
	if err != nil {
		return nil, err
	}
	countU, err := u.pActor.FavoriteRepo.CountUserFav(pctx, utils.ConvertToObjectId(req.CreatedBy))
	if err != nil {
		if err := u.pActor.ProjectRepo.DeleteProject(pctx, projectId, req.CreatedBy); err != nil {
			return nil, err
		}
	}
	fmt.Println("this is countU", countU)
	if countU == 0 {
		// create empty docs to fav
		// in case something wrong with this the project going to ge remove
		if err := u.pActor.FavoriteRepo.InsertOneFav(pctx, &favorite.FavProjectModel{
			User:      utils.ConvertToObjectId(req.CreatedBy),
			ProjectId: []string{},
			CreateAt:  utils.LocalTime(),
			UpdatedAt: utils.LocalTime(),
		}); err != nil {
			if err := u.pActor.ProjectRepo.DeleteProject(pctx, projectId, req.CreatedBy); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	countP, err := u.pActor.CommentRepo.CountCommentProject(pctx, projectId)
	if err != nil {
		if err := u.pActor.ProjectRepo.DeleteProject(pctx, projectId, req.CreatedBy); err != nil {
			return nil, err
		}
		if err := u.pActor.FavoriteRepo.DeleteFav(pctx, projectId); err != nil {
			return nil, err
		}
	}
	if countP == 0 {
		// create empty docs to comment
		// in case someething wrong with this the project going to ge remove
		if err := u.pActor.CommentRepo.InsertEmptyComment(pctx, &comment.CommentProjectModel{
			ProjectId: projectId,
			Comments:  []comment.CommentA{},
			CreateAt:  utils.LocalTime(),
			UpdatedAt: utils.LocalTime(),
		}); err != nil {
			if err := u.pActor.ProjectRepo.DeleteProject(pctx, projectId, req.CreatedBy); err != nil {
				return nil, err
			}
			if err := u.pActor.FavoriteRepo.DeleteFav(pctx, projectId); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	rawP, err := u.pActor.ProjectRepo.FindOneProject(pctx, projectId.Hex(), grpcCfg.DatacenterUrl)
	if err != nil {
		return nil, err
	}

	loc, err := utils.LocationTime()
	if err != nil {
		return nil, err
	}

	return &project.ProjectRes{
		Id:           rawP.Id.Hex(),
		Name:         rawP.Name,
		Type:         rawP.Type,
		LogoUrl:      rawP.LogoUrl,
		GithubUrl:    rawP.GithubUrl,
		WebsiteUrl:   rawP.WebsiteUrl,
		Description:  rawP.Description,
		Feedback:     rawP.Feedback,
		Category:     rawP.Category,
		FavCount:     rawP.FavCount,
		CommentCount: rawP.CommentCount,
		CreatedBy:    rawP.CreatedBy,
		CreatedAt:    rawP.CreateAt.In(loc),
		UpdatedAt:    rawP.UpdatedAt.In(loc),
	}, nil
}

func (u *projectUsecase) CreateNewQuestionUsecase(pctx context.Context, grpcCfg *config.Grpc, req *project.InsertQuestionReq) (*project.ProjectRes, error) {
	projectId, err := u.pActor.ProjectRepo.InsertOneProject(pctx, &project.ProjectModel{
		Id:           primitive.NewObjectID(),
		Name:         req.Name,
		LogoUrl:      "",
		Type:         req.Type,
		FavCount:     0,
		CommentCount: 0,
		GithubUrl:    "",
		WebsiteUrl:   "",
		Description:  req.Description,
		Feedback:     "",
		Category:     req.Category,
		CreatedBy:    req.CreatedBy,
		CreateAt:     utils.LocalTime(),
		UpdatedAt:    utils.LocalTime(),
	})
	if err != nil {
		return nil, err
	}

	countU, err := u.pActor.FavoriteRepo.CountUserFav(pctx, utils.ConvertToObjectId(req.CreatedBy))
	if err != nil {
		if err := u.pActor.ProjectRepo.DeleteProject(pctx, projectId, req.CreatedBy); err != nil {
			return nil, err
		}
	}
	fmt.Println("this is countU", countU)
	if countU == 0 {
		// create empty docs to fav
		// in case something wrong with this the project going to ge remove
		if err := u.pActor.FavoriteRepo.InsertOneFav(pctx, &favorite.FavProjectModel{
			User:      utils.ConvertToObjectId(req.CreatedBy),
			ProjectId: []string{},
			CreateAt:  utils.LocalTime(),
			UpdatedAt: utils.LocalTime(),
		}); err != nil {
			if err := u.pActor.ProjectRepo.DeleteProject(pctx, projectId, req.CreatedBy); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	countP, err := u.pActor.CommentRepo.CountCommentProject(pctx, projectId)
	if err != nil {
		if err := u.pActor.ProjectRepo.DeleteProject(pctx, projectId, req.CreatedBy); err != nil {
			return nil, err
		}
		if err := u.pActor.FavoriteRepo.DeleteFav(pctx, projectId); err != nil {
			return nil, err
		}
	}
	if countP == 0 {
		// create empty docs to comment
		// in case someething wrong with this the project going to ge remove
		if err := u.pActor.CommentRepo.InsertEmptyComment(pctx, &comment.CommentProjectModel{
			ProjectId: projectId,
			Comments:  []comment.CommentA{},
			CreateAt:  utils.LocalTime(),
			UpdatedAt: utils.LocalTime(),
		}); err != nil {
			if err := u.pActor.ProjectRepo.DeleteProject(pctx, projectId, req.CreatedBy); err != nil {
				return nil, err
			}
			if err := u.pActor.FavoriteRepo.DeleteFav(pctx, projectId); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	rawP, err := u.pActor.ProjectRepo.FindOneProject(pctx, projectId.Hex(), grpcCfg.DatacenterUrl)
	if err != nil {
		return nil, err
	}

	loc, err := utils.LocationTime()
	if err != nil {
		return nil, err
	}

	return &project.ProjectRes{
		Id:           rawP.Id.Hex(),
		Name:         rawP.Name,
		Type:         rawP.Type,
		LogoUrl:      rawP.LogoUrl,
		GithubUrl:    rawP.GithubUrl,
		WebsiteUrl:   rawP.WebsiteUrl,
		Description:  rawP.Description,
		Feedback:     rawP.Feedback,
		Category:     rawP.Category,
		FavCount:     rawP.FavCount,
		CommentCount: rawP.CommentCount,
		CreatedBy:    rawP.CreatedBy,
		CreatedAt:    rawP.CreateAt.In(loc),
		UpdatedAt:    rawP.UpdatedAt.In(loc),
	}, nil
}

func (u *projectUsecase) UpdateOneProjectUsecase(pctx context.Context, grpcCfg *config.Grpc, projectId, userId string, req *project.InsertProjectReq) (*project.ProjectResWithUser, error) {
	projectD, err := u.pActor.ProjectRepo.UpdateProject(pctx, &project.ProjectModel{
		Id:          utils.ConvertToObjectId(projectId),
		Name:        req.Name,
		LogoUrl:     req.LogoUrl,
		GithubUrl:   req.GithubUrl,
		WebsiteUrl:  req.WebsiteUrl,
		Description: req.Description,
		Feedback:    req.Feedback,
		Category:    req.Category,
		UpdatedAt:   utils.LocalTime(),
	}, userId)
	if err != nil {
		return nil, err
	}

	user, err := u.pActor.ProjectRepo.FindOneUserWithId(pctx, grpcCfg.UserUrl, &usersPb.GetUserInfoReq{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	return u.convertPModelToPWithUser(projectD, user)
}

func (u *projectUsecase) DeleteOneProjectUsecase(pctx context.Context, projectId, userId string) error {
	projectIdPri := utils.ConvertToObjectId(projectId)
	if err := u.pActor.ProjectRepo.DeleteProject(pctx, projectIdPri, userId); err != nil {
		return err
	}

	if err := u.pActor.FavoriteRepo.DeleteFav(pctx, projectIdPri); err != nil {
		return err
	}

	if err := u.pActor.CommentRepo.DeleteCommentDoc(pctx, projectIdPri); err != nil {
		return err
	}

	return nil
}

func (u *projectUsecase) FindAllProjectUsecase(pctx context.Context, grpcCfg *config.Grpc, limit, skip, pageCount int64, userId, sort, category, search, projectType string) ([]*project.ProjectResWithUser, int64, error) {

	var pjType string
	switch projectType {
	case "project":
		pjType = "p"
	case "question":
		pjType = "q"
	default:
		pjType = projectType
	}

	var count int64 = pageCount
	if count == 0 {
		temp, err := u.pActor.ProjectRepo.CountProject(pctx, pjType)
		if err != nil {
			return nil, temp, err
		}
		count = temp
	}

	rawProjects, err := u.pActor.ProjectRepo.FindAllProjectDatacenter(pctx, grpcCfg.DatacenterUrl, limit, skip, userId, sort, category, search, pjType)
	if err != nil {
		return nil, count, err
	}

	usersId := u.accumateUserIdByProjects(rawProjects)
	usersInfo, err := u.pActor.ProjectRepo.FindManyUserInfo(pctx, grpcCfg.UserUrl, &usersPb.GetManyUserInfoForProjectReq{UsersId: usersId})
	if err != nil {
		return nil, count, err
	}

	var projectRes []*project.ProjectResWithUser
	loc, err := utils.LocationTime()
	if err != nil {
		return nil, count, err
	}

	// authen user
	if len(userId) > 5 {

		rawFav, err := u.pActor.FavoriteRepo.GetAllProjectInUser(pctx, utils.ConvertToObjectId(userId))
		if err != nil {
			return nil, count, err
		}

		for _, p := range rawProjects {
			parsedTime := p.CreateAt.In(loc)
			var wg sync.WaitGroup
			var match bool

			owner := new(users.UserProfileRes)
			wg.Add(len(usersInfo.UsersProfile))
			for _, v := range usersInfo.UsersProfile {
				go func(v *usersPb.UserProfile) {
					defer wg.Done()
					if v.UserId == p.CreatedBy {
						owner.Id = v.UserId
						owner.Email = v.Email
						owner.Firstname = v.FirstName
						owner.Lastname = v.LastName
						owner.Username = v.UserName
						owner.Profile = v.Profile
					}
				}(v)
			}
			wg.Wait()

			for _, fp := range rawFav.ProjectId {
				if p.Id.Hex() == fp {
					projectRes = append(projectRes, u.assignProjectRes(p, true, parsedTime, owner))
					match = true
					break
				}
			}
			if !match {
				projectRes = append(projectRes, u.assignProjectRes(p, false, parsedTime, owner))
			}
		}

		// unauthen user no check fav
	} else {
		var wg sync.WaitGroup
		for _, p := range rawProjects {

			owner := new(users.UserProfileRes)
			wg.Add(len(usersInfo.UsersProfile))
			for _, v := range usersInfo.UsersProfile {
				go func(v *usersPb.UserProfile) {
					defer wg.Done()
					if v.UserId == p.CreatedBy {
						owner.Id = v.UserId
						owner.Email = v.Email
						owner.Firstname = v.FirstName
						owner.Lastname = v.LastName
						owner.Username = v.UserName
						owner.Profile = v.Profile
					}
				}(v)
			}
			wg.Wait()

			parsedTime := p.CreateAt.In(loc)
			projectRes = append(projectRes, u.assignProjectRes(p, false, parsedTime, owner))
		}
	}

	return projectRes, count, nil
}

func (u *projectUsecase) FindOneProjectUsecase(pctx context.Context, grpcCfg *config.Grpc, projectId, userId string) (*project.FullProjectRes, error) {

	projectD, err := u.pActor.ProjectRepo.FindOneProject(pctx, projectId, grpcCfg.DatacenterUrl)
	if err != nil {
		return nil, err
	}

	rawComment, err := u.pActor.CommentRepo.FindCommentByProjectId(pctx, utils.ConvertToObjectId(projectId))
	if err != nil {
		return nil, err
	}

	usersId := u.accumateUserId(rawComment)
	// add the owner of the project to get the info of owner
	usersId = append(usersId, projectD.CreatedBy)

	usersInfo, err := u.pActor.ProjectRepo.FindManyUserInfo(pctx, grpcCfg.UserUrl, &usersPb.GetManyUserInfoForProjectReq{UsersId: usersId})
	if err != nil {
		return nil, err
	}

	if len(userId) > 5 {
		countFav, err := u.pActor.FavoriteRepo.CountUserFavByProjectId(pctx, utils.ConvertToObjectId(userId), projectD.Id)
		if err != nil {
			return nil, err
		}
		var fav bool = false
		if countFav != 0 {
			fav = true
		}
		return u.convertPDatacenterToPWithUser(projectD, fav, rawComment, usersInfo.UsersProfile)
	}

	return u.convertPDatacenterToPWithUser(projectD, false, rawComment, usersInfo.UsersProfile)
}

func (u *projectUsecase) GetNewestProject(pctx context.Context, limit int64, exceptProjectId string) ([]*project.RandomProjectDisplay, error) {
	rawP, err := u.pActor.ProjectRepo.FindLatestProjects(pctx, 3, utils.ConvertToObjectId(exceptProjectId))
	if err != nil {
		return nil, err
	}

	var projectRes []*project.RandomProjectDisplay
	for _, v := range rawP {
		projectRes = append(projectRes, &project.RandomProjectDisplay{
			Id:           v.Id.Hex(),
			Name:         v.Name,
			Type:         v.Type,
			LogoUrl:      v.LogoUrl,
			Category:     v.Category,
			Description:  v.Description,
			FavCount:     v.FavCount,
			CommentCount: v.CommentCount,
		})
	}

	return projectRes, nil

}

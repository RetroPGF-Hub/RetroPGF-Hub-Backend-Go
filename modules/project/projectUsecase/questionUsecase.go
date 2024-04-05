package projectusecase

// func (u *projectUsecase) CreateNewQuestionUsecase(pctx context.Context, grpcCfg *config.Grpc, req *project.InsertQuestionReq) (*project.QuestionRes, error) {
// 	projectId, err := u.pActor.ProjectRepo.InsertOneQuestion(pctx, &project.QuestionModel{
// 		Id:        primitive.NewObjectID(),
// 		Title:     req.Title,
// 		Detail:    req.Detail,
// 		CreatedBy: req.CreatedBy,
// 		CreateAt:  utils.LocalTime(),
// 		UpdatedAt: utils.LocalTime(),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	countU, err := u.pActor.FavoriteRepo.CountUserFav(pctx, utils.ConvertToObjectId(req.CreatedBy))
// 	if err != nil {
// 		if err := u.pActor.ProjectRepo.DeleteProject(pctx, projectId, req.CreatedBy); err != nil {
// 			return nil, err
// 		}
// 	}
// 	fmt.Println("this is countU", countU)
// 	if countU == 0 {
// 		// create empty docs to fav
// 		// in case something wrong with this the project going to ge remove
// 		if err := u.pActor.FavoriteRepo.InsertOneFav(pctx, &favorite.FavProjectModel{
// 			User:      utils.ConvertToObjectId(req.CreatedBy),
// 			ProjectId: []string{},
// 			CreateAt:  utils.LocalTime(),
// 			UpdatedAt: utils.LocalTime(),
// 		}); err != nil {
// 			if err := u.pActor.ProjectRepo.DeleteProject(pctx, projectId, req.CreatedBy); err != nil {
// 				return nil, err
// 			}
// 			return nil, err
// 		}
// 	}

// 	countP, err := u.pActor.CommentRepo.CountCommentProject(pctx, projectId)
// 	if err != nil {
// 		if err := u.pActor.ProjectRepo.DeleteProject(pctx, projectId, req.CreatedBy); err != nil {
// 			return nil, err
// 		}
// 		if err := u.pActor.FavoriteRepo.DeleteFav(pctx, projectId); err != nil {
// 			return nil, err
// 		}
// 	}
// 	if countP == 0 {
// 		// create empty docs to comment
// 		// in case someething wrong with this the project going to ge remove
// 		if err := u.pActor.CommentRepo.InsertEmptyComment(pctx, &comment.CommentProjectModel{
// 			ProjectId: projectId,
// 			Comments:  []comment.CommentA{},
// 			CreateAt:  utils.LocalTime(),
// 			UpdatedAt: utils.LocalTime(),
// 		}); err != nil {
// 			if err := u.pActor.ProjectRepo.DeleteProject(pctx, projectId, req.CreatedBy); err != nil {
// 				return nil, err
// 			}
// 			if err := u.pActor.FavoriteRepo.DeleteFav(pctx, projectId); err != nil {
// 				return nil, err
// 			}
// 			return nil, err
// 		}
// 	}

// 	rawP, err := u.pActor.ProjectRepo.FindOneProject(pctx, projectId.Hex(), grpcCfg.DatacenterUrl)
// 	if err != nil {
// 		return nil, err
// 	}

// 	loc, err := utils.LocationTime()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &project.ProjectRes{
// 		Id:             rawP.Id.Hex(),
// 		Name:           rawP.Name,
// 		LogoUrl:        rawP.LogoUrl,
// 		GithubUrl:      rawP.GithubUrl,
// 		WebsiteUrl:     rawP.WebsiteUrl,
// 		CryptoCategory: rawP.CryptoCategory,
// 		Description:    rawP.Description,
// 		Feedback:       rawP.Feedback,
// 		Category:       rawP.Category,
// 		FavCount:       rawP.FavCount,
// 		CommentCount:   rawP.CommentCount,
// 		Contact:        rawP.Contact,
// 		CreatedBy:      rawP.CreatedBy,
// 		CreatedAt:      rawP.CreateAt.In(loc),
// 		UpdatedAt:      rawP.UpdatedAt.In(loc),
// 	}, nil
// }

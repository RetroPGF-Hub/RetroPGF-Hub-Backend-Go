package datacenterusecase

import (
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/config"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter"
	datacenterPb "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterPb"
	datacenterrepository "RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/datacenter/datacenterRepository"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/modules/project"
	"RetroPGF-Hub/RetroPGF-Hub-Backend-Go/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// ambigouse variable name
// urlId == cacheId

type (
	DatacenterUsecaseService interface {
		GetAllProjectUsecase(pctx context.Context, limit, skip int64) (*datacenterPb.GetProjectDataCenterRes, error)
		GetSingleProjectUsecase(pctx context.Context, projectId string) (*datacenterPb.GetSingleProjectDataCenterRes, error)
		InsertUrlCache(pctx context.Context, url string) (string, error)
		DeleteUrlCache(pctx context.Context, url string) error
		FindManyUrlsCache(pctx context.Context) ([]*datacenter.CacheModel, error)
		FindCacheData(pctx context.Context, cacheId string) (*any, error)
		CronJobUpdateCache(pctx context.Context) error
		TriggerUpdateCache(pctx context.Context, cacheId string) (*any, error)
	}

	datacenterUsecase struct {
		cfg            *config.Grpc
		datacenterRepo datacenterrepository.DatacenterRepositoryService
	}
)

func NewDatacenterUsecase(datacenterRepo datacenterrepository.DatacenterRepositoryService, cfg *config.Grpc) DatacenterUsecaseService {
	return &datacenterUsecase{
		datacenterRepo: datacenterRepo,
		cfg:            cfg,
	}
}

func (u *datacenterUsecase) GetAllProjectUsecase(pctx context.Context, limit, skip int64) (*datacenterPb.GetProjectDataCenterRes, error) {
	projects, err := u.datacenterRepo.GetAllProjectRepo(pctx, limit, skip)
	if err != nil {
		return nil, err
	}

	var result []*datacenterPb.ProjectRes

	for _, v := range projects {
		result = append(result, u.convertToPbProject(v))
	}

	return &datacenterPb.GetProjectDataCenterRes{
		Projects: result,
	}, nil

}

func (u *datacenterUsecase) GetSingleProjectUsecase(pctx context.Context, projectId string) (*datacenterPb.GetSingleProjectDataCenterRes, error) {
	projects, err := u.datacenterRepo.GetSingleProjectRepo(pctx, utils.ConvertToObjectId(projectId))
	if err != nil {
		return nil, err
	}

	return &datacenterPb.GetSingleProjectDataCenterRes{
		Projects: u.convertToPbProject(projects),
	}, nil

}

func (u *datacenterUsecase) convertToPbProject(p *project.ProjectModel) *datacenterPb.ProjectRes {
	return &datacenterPb.ProjectRes{
		Id:           p.Id.Hex(),
		Name:         p.Name,
		LogoUrl:      p.LogoUrl,
		GithubUrl:    p.GithubUrl,
		WebsiteUrl:   p.WebsiteUrl,
		Description:  p.Description,
		Feedback:     p.Feedback,
		Category:     p.Category,
		FavCount:     p.FavCount,
		CommentCount: p.CommentCount,
		CreatedBy:    p.CreatedBy,
		CreatedAt:    p.CreateAt.String(),
		UpdatedAt:    p.UpdatedAt.String(),
	}
}

func (u *datacenterUsecase) InsertUrlCache(pctx context.Context, url string) (string, error) {
	id, err := u.datacenterRepo.InsertUrlCache(pctx, &datacenter.CacheModel{Url: url})
	if err != nil {
		return "", err
	}

	body, err := getReq(url)
	if err != nil {
		return "", err
	}

	if err := u.datacenterRepo.InsertCacheToRedis(pctx, id.Hex(), string(body)); err != nil {
		if err := u.datacenterRepo.DeleteUrlCache(pctx, id); err != nil {
			return "", errors.New("failed to insert cache to redis and also failed to delete the url in db")
		}
		return "", errors.New("failed to insert cache to redis")
	}

	return id.Hex(), err
}

func (u *datacenterUsecase) DeleteUrlCache(pctx context.Context, cacheId string) error {
	err := u.datacenterRepo.DeleteUrlCache(pctx, utils.ConvertToObjectId(cacheId))
	if err != nil {
		return err
	}
	if err := u.datacenterRepo.DeleteCacheFromRedis(pctx, cacheId); err != nil {
		return err
	}

	return nil
}

func (u *datacenterUsecase) FindManyUrlsCache(pctx context.Context) ([]*datacenter.CacheModel, error) {
	data, err := u.datacenterRepo.GetAllUrlCache(pctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *datacenterUsecase) FindCacheData(pctx context.Context, cacheId string) (*any, error) {
	rawD, err := u.datacenterRepo.GetCacheFromRedis(pctx, cacheId)
	if err != nil {
		return nil, err
	}

	temp := new(any)
	err = json.Unmarshal([]byte(rawD), &temp)
	if err != nil {
		return temp, err
	}
	return temp, err

}

func (u *datacenterUsecase) TriggerUpdateCache(pctx context.Context, cacheId string) (*any, error) {
	rawD, err := u.datacenterRepo.FindOneCache(pctx, utils.ConvertToObjectId(cacheId))
	if err != nil {
		return nil, err
	}

	body, err := getReq(rawD.Url)
	if err != nil {
		return nil, err
	}

	if err := u.datacenterRepo.InsertCacheToRedis(pctx, rawD.UrlId.Hex(), string(body)); err != nil {
		return nil, errors.New("failed to insert cache to redis")
	}
	result := new(any)
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (u *datacenterUsecase) CronJobUpdateCache(pctx context.Context) error {
	data, err := u.datacenterRepo.GetAllUrlCache(pctx)
	if err != nil {
		return err
	}

	pipeLineData := make([]*datacenter.PipeLineCache, 0)

	for _, v := range data {
		body, err := getReq(v.Url)
		if err != nil {
			return err
		}
		pipeLineData = append(pipeLineData, &datacenter.PipeLineCache{
			CacheId:   v.UrlId.Hex(),
			CacheData: string(body),
		})

		// if err := u.datacenterRepo.InsertCacheToRedis(pctx, v.UrlId.Hex(), string(body)); err != nil {
		// 	return errors.New("failed to insert cache to redis")
		// }

		// pipeline.Set(pctx, v.UrlId.Hex(), string(body), 0)

	}

	if err := u.datacenterRepo.InsertManyCacheToRedis(pctx, pipeLineData); err != nil {
		return err
	}

	return nil
}

func getReq(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

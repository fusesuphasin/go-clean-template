package service

import (
	"strconv"

	"github.com/fusesuphasin/go-fiber/app/domain"
	"github.com/fusesuphasin/go-fiber/app/repository"
	"github.com/gofiber/fiber/v2"
)

var c *fiber.Ctx

type RoleService struct {
	RoleRepository repository.RoleRepository
}

func (rs RoleService) GetAll(page string, limit string) (roles *[]domain.Role, err int) {
	pageInt, errPage := strconv.Atoi(page)

	limitInt, errLimit := strconv.Atoi(limit)

	if page != "" && limit != "" {
		if errPage != nil || errLimit != nil {
			return nil, 1
		}
	}

	roleData, _ := rs.RoleRepository.GetAll(pageInt, limitInt)

	return roleData, 0
}

func (rs RoleService) CreateRole(Role *domain.Role) (user *domain.Role, err error) {
	data, err := rs.RoleRepository.Insert(Role)

	return data, err
}

func (rs RoleService) UpdateRole(role_id string, Role *domain.Role) (user *domain.Role, err error) {
	data, err := rs.RoleRepository.Update(role_id, Role)

	return data, err
}

func (rs RoleService) CheckDuplicateNameRole(name string) int64 {
	data := rs.RoleRepository.CountByName(name)

	return data
}

func (rs RoleService) GetById(role_id string) (role *domain.Role) {
	roleData, _ := rs.RoleRepository.FindById(role_id)

	return roleData
}

func (rs RoleService) DeleteRole(role_id string) error {
	deleteRole := rs.RoleRepository.Delete(role_id)

	return deleteRole
}
package contracts

import (
	"encoding/json"
	"strconv"

	"github.com/meshplus/bitxhub-core/boltvm"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub/internal/repo"
)

const (
	adminRolesKey = "admin-roles"
)

type Role struct {
	boltvm.Stub
}

func (r *Role) GetRole() *boltvm.Response {
	var admins []*repo.Admin
	r.GetObject(adminRolesKey, &admins)

	for _, admin := range admins {
		if admin.Address == r.Caller() {
			return boltvm.Success([]byte("admin"))
		}
	}

	res := r.CrossInvoke(constant.AppchainMgrContractAddr.String(), "Appchain")
	if !res.Ok {
		return boltvm.Success([]byte("none"))
	}

	return boltvm.Success([]byte("appchain_admin"))
}

func (r *Role) IsAdmin(address string) *boltvm.Response {
	var admins []*repo.Admin
	r.GetObject(adminRolesKey, &admins)

	for _, admin := range admins {
		if admin.Address == address {
			return boltvm.Success([]byte(strconv.FormatBool(true)))
		}
	}

	return boltvm.Success([]byte(strconv.FormatBool(false)))
}

func (r *Role) GetAdminRoles() *boltvm.Response {
	var admins []*repo.Admin
	r.GetObject(adminRolesKey, &admins)

	ret, err := json.Marshal(admins)
	if err != nil {
		return boltvm.Error(err.Error())
	}

	return boltvm.Success(ret)
}

func (r *Role) SetAdminRoles(addrs string) *boltvm.Response {
	as := make([]string, 0)
	if err := json.Unmarshal([]byte(addrs), &as); err != nil {
		return boltvm.Error(err.Error())
	}

	admins := make([]*repo.Admin, 0)
	for _, addr := range as {
		admins = append(admins, &repo.Admin{
			Address: addr,
			Weight:  1,
		})
	}

	r.SetObject(adminRolesKey, admins)
	return boltvm.Success(nil)
}

func (r *Role) GetRoleWeight(address string) *boltvm.Response {
	var admins []*repo.Admin
	r.GetObject(adminRolesKey, &admins)

	for _, admin := range admins {
		if admin.Address == address {
			return boltvm.Success([]byte(strconv.Itoa(int(admin.Weight))))
		}
	}

	return boltvm.Error("account at the address does not exist:" + address)
}

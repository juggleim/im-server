package services

import (
	"im-server/commons/tools"
	"im-server/services/admingateway/dbs"
	"time"
)

func CheckLogin(account, password string) (AdminErrorCode, *Account) {
	dao := dbs.AccountDao{}
	defaultAccount, err := dao.FindByAccount("admin")
	if err != nil || defaultAccount == nil {
		//init default account
		dao.Create(dbs.AccountDao{
			Account:       "admin",
			Password:      tools.SHA1("123456"),
			CreatedTime:   time.Now(),
			UpdatedTime:   time.Now(),
			State:         0,
			ParentAccount: "",
		})
	}

	password = tools.SHA1(password)
	admin, err := dao.FindByAccountPassword(account, password)
	if err == nil && admin != nil {
		if admin.State != 0 {
			return AdminErrorCode_AccountForbidden, nil
		}
		return AdminErrorCode_Success, &Account{
			Account:       admin.Account,
			State:         admin.State,
			ParentAccount: admin.ParentAccount,
			RoleId:        admin.RoleId,
			CreatedTime:   admin.CreatedTime.UnixMilli(),
			UpdatedTime:   admin.UpdatedTime.UnixMilli(),
		}
	}
	return AdminErrorCode_LoginFail, nil
}

func CheckAccountState(account string) AdminErrorCode {
	if account == "jugglechat" {
		return AdminErrorCode_Success
	}
	dao := dbs.AccountDao{}
	admin, err := dao.FindByAccount(account)
	if err != nil || admin == nil {
		return AdminErrorCode_AccountNotExist
	} else {
		if admin.State == 0 {
			return AdminErrorCode_Success
		} else {
			return AdminErrorCode_AccountForbidden
		}
	}
}

func UpdPassword(account, password, newPassword string) AdminErrorCode {
	dao := dbs.AccountDao{}
	password = tools.SHA1(password)
	admin, err := dao.FindByAccountPassword(account, password)
	if err != nil || admin == nil {
		return AdminErrorCode_UpdPwdFail
	}
	newPassword = tools.SHA1(newPassword)
	dao.UpdatePassword(account, newPassword)
	return AdminErrorCode_Success
}

func AddAccount(parentAccount, account, password string, roleId int) AdminErrorCode {
	dao := dbs.AccountDao{}
	password = tools.SHA1(password)
	err := dao.Create(dbs.AccountDao{
		Account:       account,
		Password:      password,
		ParentAccount: parentAccount,
		UpdatedTime:   time.Now(),
		CreatedTime:   time.Now(),
		RoleId:        roleId,
	})
	if err != nil {
		return AdminErrorCode_AccountExisted
	}
	return AdminErrorCode_Success
}

func DisableAccounts(accounts []string, isDisable int) AdminErrorCode {
	dao := dbs.AccountDao{}
	dao.UpdateState(accounts, isDisable)
	return AdminErrorCode_Success
}
func DeleteAccounts(accounts []string) AdminErrorCode {
	dao := dbs.AccountDao{}
	dao.DeleteAccounts(accounts)
	return AdminErrorCode_Success
}

func QryAccounts(limit int64, offset string) *Accounts {
	accounts := &Accounts{
		Items:   []*Account{},
		HasMore: false,
		Offset:  "",
	}
	dao := dbs.AccountDao{}
	offsetInt, err := tools.DecodeInt(offset)
	if err != nil {
		offsetInt = 0
	}
	dbAccounts, err := dao.QryAccounts(limit+1, offsetInt)
	if err == nil {
		if len(dbAccounts) > int(limit) {
			dbAccounts = dbAccounts[:len(dbAccounts)-1]
			accounts.HasMore = true
		}
		var id int64 = 0
		for _, dbAccount := range dbAccounts {
			accounts.Items = append(accounts.Items, &Account{
				Account:       dbAccount.Account,
				State:         dbAccount.State,
				CreatedTime:   dbAccount.CreatedTime.UnixMilli(),
				UpdatedTime:   dbAccount.UpdatedTime.UnixMilli(),
				ParentAccount: dbAccount.ParentAccount,
				RoleId:        dbAccount.RoleId,
			})
			if dbAccount.ID > id {
				id = dbAccount.ID
			}
		}
		if id > 0 {
			offset, _ := tools.EncodeInt(id)
			accounts.Offset = offset
		}
	}
	return accounts
}

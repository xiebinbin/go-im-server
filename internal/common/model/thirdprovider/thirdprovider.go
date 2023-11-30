package thirdprovider

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"imsdk/internal/common/dao/thirdprovider"
	"imsdk/internal/common/model/errors"
	"imsdk/pkg/eccsign"
	"imsdk/pkg/funcs"
	"imsdk/pkg/sdk"
)

func CreateProvider() {
	//ak := "68oni7jrg31qcsaijtg76qln"
	ak := "OFFICIAL"
	id := funcs.Md516(ak)
	fmt.Println("id:", id, ak)
	seed, _ := eccsign.GenerateSeed(id, ak)
	pubKey, _ := eccsign.GetPublicKeyStr(seed)
	priKey, _ := eccsign.GetPrivateKey(seed)
	t := funcs.GetMillis()
	addData := thirdprovider.ThirdProvider{
		ID:        id,
		AK:        ak,
		SK:        priKey,
		PK:        pubKey,
		Status:    sdk.StatusNormal,
		CreatedAt: t,
		UpdatedAt: t,
	}
	err := thirdprovider.New().UpsertOne(addData)
	if err != nil {
		return
	}
	return
}

func GetProviderByAKey(ak string) (thirdprovider.ThirdProvider, error) {
	data, err := thirdprovider.New().GetByAKey(ak)
	if err == mongo.ErrNoDocuments {
		return thirdprovider.ThirdProvider{}, errors.ErrSdkAKNotMatch
	}
	if err != nil {
		return thirdprovider.ThirdProvider{}, errors.ErrSdkDefErr
	}
	return data, err
}
